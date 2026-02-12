package memory

import (
	"Mengbot/config"
	"Mengbot/internal/core/llm"
	"Mengbot/internal/core/model"
	"Mengbot/pkg/logger"
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/pgvector/pgvector-go"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

func StartDiary() {
	c := cron.New(cron.WithSeconds())

	_, err := c.AddFunc("0 0 4 * * *", func() {
		botName := config.Conf.Bot.Name
		logger.Infof("🌙 可爱的%s开始尝试写日记了喵....", botName)

		// 引入重试机制：最多重试 3 次
		maxRetries := 3
		for i := 0; i < maxRetries; i++ {
			err := generateDiary()
			if err == nil {
				logger.Infof("✅ 可爱的%s成功写完了日记！", botName)
				return // 成功则退出
			}

			logger.Errorf("⚠️ 写日记失败 (第 %d/%d 次尝试): %v", i+1, maxRetries, err)

			// 失败等待一段时间再重试 (指数退避: 10s, 20s, 40s...)
			time.Sleep(time.Duration(10*(i+1)) * time.Second)
		}

		logger.Errorf("❌ 可爱的%s彻底放弃了今天的日记，呜呜呜...", botName)
	})
	if err != nil {
		logger.Errorf("定时任务建立失败, err: %v", err)
	}
	c.Start()
}

func generateDiary() error {
	masterId, err := strconv.ParseInt(config.Conf.Master.MasterID, 10, 64)
	if err != nil {
		return err
	}

	value, ok := Store.Load(masterId)
	if !ok {
		return nil
	}

	// 拿到 UserMemory 结构体指针
	userMem := value.(*model.UserMemory)

	userMem.Mu.Lock()

	masterUser := userMem.Messages

	if len(masterUser) == 0 {
		userMem.Mu.Unlock()
		return nil
	}

	var (
		splitIndex = len(masterUser)
		splitTime  = time.Now()
		//.Add(-20 * time.Minute)
		foundNew = false
	)

	for k, msg := range masterUser {
		if msg.Time.After(splitTime) {
			splitIndex = k
			foundNew = true
			break
		}
	}

	if foundNew && splitIndex == 0 {
		userMem.Mu.Unlock()
		return nil
	}

	historyMsg := masterUser[:splitIndex]

	userMem.Messages = masterUser[splitIndex:]

	userMem.Mu.Unlock()

	if len(historyMsg) == 0 {
		return nil
	}

	//拼接文本
	var r strings.Builder
	for _, message := range historyMsg {
		r.WriteString("\n" + message.TimeString + " " + message.NickName + "：" + message.Content + "\n" + message.ApplyName + "：" + message.ApplyContent + "\n")
	}
	today := r.String()

	// 调用 LLM 生成日记
	diaryContent, err := llm.CallDiary(context.Background(), today)
	if err != nil {
		return err
	}

	// 调用 LLM 生成日记元数据
	metadata, err := llm.CallMetadata(context.Background(), diaryContent)
	if err != nil {
		return err
	}
	embeddings, err := llm.CallDiaryEmbedding(context.Background(), metadata)
	if err != nil {
		return err
	}

	diary := &model.DiaryMessage{
		UserId:    config.Conf.Master.MasterID,
		Content:   diaryContent,
		Metadata:  *metadata,
		Embedding: pgvector.NewVector(embeddings),
	}
	// 保存到数据库
	err = db.Create(diary).Error
	if err != nil {
		return err
	}
	return nil
}

func SearchDiary(chat string) ([]model.DiaryMessage, error) {
	var results []model.SearchResult
	embeddings, err := llm.CallChatEmbedding(context.Background(), chat)
	if err != nil {
		return nil, err
	}

	targetVec := pgvector.NewVector(embeddings)

	err = db.Model(&model.DiaryMessage{}).
		Select("*, (1 - (embedding <=> ?)) as score", targetVec).
		Order(gorm.Expr("embedding <=> ?", targetVec)). // 分数越大越相似
		Limit(5).
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, nil
	}

	var diaries []model.DiaryMessage
	for _, result := range results {
		if result.Score < 0.7 {
			continue
		}
		diaries = append(diaries, result.DiaryMessage)
	}
	return diaries, nil
}
