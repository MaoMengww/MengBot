package core

import (
	"Mengbot/config"
	"Mengbot/internal/core/graph"
	"Mengbot/internal/core/llm"
	"Mengbot/internal/core/memory"
	"Mengbot/internal/core/prompt"
	"Mengbot/pkg/logger"
	"context"
	"math/rand"
	"strconv"
	"time"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func GetRandomReapOnWrong() string {
	reapOnWrong := []string{
		"呜... 脑子里突然乱糟糟的，刚才的话没听清喵！可以再说一遍吗？",
		"唔... 刚才想得太深奥，大脑过载了喵！",
		"脑子糊掉了喵，再试一次？",
		"看什么看？猫梦刚才只是去抓了一下路过的电子老鼠而已！再给你一次机会，重新说一遍喵！",
		"(哈欠) 刚才那句话被次元风暴吹走了喵。趁猫梦还没睡着，赶紧再说一遍。",
	}
	return reapOnWrong[rand.Intn(len(reapOnWrong))]
}

func Init() {
	llm.InitModel()
	run := graph.GraphBotInit(context.Background())
	if run == nil {
		logger.Errorf("初始化图灵机器人失败")
		return
	}
	zero.OnMessage(zero.OnlyToMe).Handle(func(ctx *zero.Ctx) {
		var userRole, userName string

		masterId, err := strconv.Atoi(config.Conf.Master.MasterID)
		if err != nil {
			logger.Errorf("转换主ID失败: %v", err)
			return
		}
		if ctx.Event.UserID == int64(masterId) {
			userRole = "master"
			userName = config.Conf.Master.MasterName
		} else {
			userRole = "user"
			userName = ctx.Event.Sender.NickName
		}
		time := time.Now().Format("2006-01-02 15:04:05")
		prompt := &prompt.Message{
			Query:     ctx.Event.Message.String(),
			UserRole:  userRole,
			UserName:  userName,
			History:   memory.GetShortMemory(ctx.Event.UserID),
			Time:      time,
			Documents: nil,
		}
		resp, err := (*run).Invoke(context.Background(), prompt)
		if err != nil {
			logger.Errorf("聊天失败: %v", err)
			ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("呜... 脑子里突然乱糟糟的，刚才的话没听清喵！可以再说一遍吗？")))
			return
		}
		memory.AppendShortMemory(ctx.Event.UserID, memory.MemoryMessage{
			Time:         time,
			NickName:     userName,
			Content:      ctx.Event.Message.String(),
			ApplyName:    config.Conf.Bot.Name,
			ApplyContent: resp,
		})
		ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text(resp)))
	})
}
