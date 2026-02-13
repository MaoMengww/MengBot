package model

import (
	"sync"
	"time"

	"github.com/pgvector/pgvector-go"
)

type UserProfile struct {
	Id        int64           `gorm:"primaryKey"`
	UserId    string          `gorm:"index"`
	Profile   string          `gorm:"type:text"`
}

type UserMemory struct {
	Mu       sync.Mutex
	Messages []MemoryMessage
}

type MemoryMessage struct {
	Time         time.Time
	TimeString   string
	NickName     string
	Content      string
	ApplyName    string
	ApplyContent string
}

type Message struct {
	// 1. 核心交互数据
	Query    string          `json:"query"`    // 用户当前说的话
	History  []MemoryMessage `json:"history"`  // 之前的聊天记录
	Strategy string          `json:"strategy"` // 策略 (TOOL / CHAT)

	// 2. 状态与感知 (Context)
	UserName string `json:"user_name"` // 用户名称
	UserRole string `json:"user_role"` // 用户角色 (master / user)
	Time     string `json:"time"`      // 当前时间

	// 3. RAG / 知识库 (可选)
	Documents []string `json:"documents"` // 检索到的相关文档片段
}

func NewMessage(query string, history []MemoryMessage, mood string, userName string, userRole string, time string, documents []string) *Message {
	return &Message{
		Query:     query,
		History:   history,
		UserName:  userName,
		UserRole:  userRole,
		Time:      time,
		Documents: documents,
	}
}

type DiaryMessage struct {
	Id        int64           `gorm:"primaryKey"`
	UserId    string          `gorm:"index"`
	Content   string          `gorm:"type:text"`
	Metadata  DiaryMetadata   `gorm:"type:jsonb"`
	Embedding pgvector.Vector `gorm:"type:vector(1536)"`
}

type DiaryMetadata struct {
	Topic      []string `json:"topic"`
	Keywords   []string `json:"keywords"`
	Mood       string   `json:"mood"`
	Importance int      `json:"importance"`
	Summary    string   `json:"summary"`
}

type SearchResult struct {
	DiaryMessage
	Score float32 `json:"score"`
}
