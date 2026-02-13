package memory

import (
	"Mengbot/config"
	"Mengbot/internal/core/model"
	"sync"
)

var (
	Store sync.Map
	RowStore sync.Map
)

// GetShortMemory 获取短期记忆（返回深拷贝，绝对安全）
func GetShortMemory(userId int64) []model.MemoryMessage {
	value, ok := Store.Load(userId)
	if !ok {
		return nil
	}
	userMem := value.(*model.UserMemory)

	userMem.Mu.Lock()
	defer userMem.Mu.Unlock()

	// 如果直接返回 userMem.Messages，外层读取时会发生并发冲突
	if len(userMem.Messages) == 0 {
		return nil
	}

	result := make([]model.MemoryMessage, len(userMem.Messages))
	copy(result, userMem.Messages)

	return result
}

// AppendShortMemory 追加记忆
func AppendShortMemory(userId int64, message model.MemoryMessage) {
	// 1. 原子性加载或创建 (LoadOrStore)
	actual, _ := Store.LoadOrStore(userId, &model.UserMemory{
		Messages: make([]model.MemoryMessage, 0, config.Conf.Bot.Memory.WindowLength),
	})
	row, _ := RowStore.LoadOrStore(userId, &model.UserMemory{
		Messages: make([]model.MemoryMessage, 0),
	})
	userMem := actual.(*model.UserMemory)
	rowMem := row.(*model.UserMemory)

	// 2. 加锁写入
	userMem.Mu.Lock()
	defer userMem.Mu.Unlock()

	windowLen := config.Conf.Bot.Memory.WindowLength

	if len(userMem.Messages) >= windowLen {
		rowMem.Messages = append(rowMem.Messages, userMem.Messages[0])
		userMem.Messages = userMem.Messages[1:]
	}
	userMem.Messages = append(userMem.Messages, message)
}
