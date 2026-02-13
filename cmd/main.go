package main

import (
	"Mengbot/config"
	"Mengbot/pkg/logger"
	"Mengbot/plugins"

	"Mengbot/internal/core"
	"Mengbot/internal/core/memory"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
)

func main() {
	logger.InitLogger()
	config.Init()
	memory.InitPgsql()
	plugins.Init()
	core.Init()
	zero.RunAndBlock(&zero.Config{
		NickName:      []string{"bot"},
		CommandPrefix: "/",
		SuperUsers:    []int64{123456},
		Driver: []zero.Driver{
			driver.NewWebSocketServer(16, "ws://127.0.0.1:3001/", ""),
		},
	}, nil)
}
