package memory

import (
	"Mengbot/config"
	"Mengbot/pkg/logger"
	"strings"

	"github.com/robfig/cron/v3"
)

func StartDiary(){
	c := cron.New(cron.WithSeconds())

	_, err := c.AddFunc("0 0 4 * * *", func() {
		logger.Infof("可爱的%s开始写日记了喵....", config.Conf.Bot.Name)
		if err := generateDiary(); err != nil{
			logger.Errorf("可爱的%s没有头绪选择放弃了今天的日记:%v", config.Conf.Bot.Name, err)
		}
	})
	if err != nil{
		logger.Errorf("定时任务建立失败, err: %v", err)
	}
}

func generateDiary() error {
	var r strings.Builder
	if len(ShortMemory) == 0{
		return  nil
	}
	for _, msg := range ShortMemory{
		
	}
}