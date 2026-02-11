package memory

import "Mengbot/config"

var (
	ShortMemory = map[int64][]MemoryMessage{}
)

func GetShortMemory(key int64) []MemoryMessage {
	return ShortMemory[key]
}

type MemoryMessage struct {
	Time string
	NickName string
	Content  string

	ApplyName string
	ApplyContent string
}

func AppendShortMemory(key int64, message MemoryMessage) {
	windowLenth := len(ShortMemory[key])
	if windowLenth >= config.Conf.Bot.Memory.WindowLength {
		ShortMemory[key] = ShortMemory[key][1:]
	}
	ShortMemory[key] = append(ShortMemory[key], message)
}
