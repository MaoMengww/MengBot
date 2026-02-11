package prompt

import "Mengbot/config"

func GetRouterPrompt() string {
	return config.Conf.Prompt.RouterPrompt
}

func GetChatPrompt() string {
	return config.Conf.Prompt.ChatPrompt
}
