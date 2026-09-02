package engine

import (
	"github.com/gdamore/tcell/v2"
)

type LogMessage struct {
	Text  string
	Color tcell.Color
}

type MessageLog struct {
	Messages []LogMessage
	MaxLines int
}

func NewMessageLog(maxLines int) *MessageLog {
	return &MessageLog{
		Messages: make([]LogMessage, 0),
		MaxLines: maxLines,
	}
}

func (l *MessageLog) Add(text string, color tcell.Color) {
	l.Messages = append(l.Messages, LogMessage{
		Text:  text,
		Color: color,
	})
	if len(l.Messages) > l.MaxLines {
		l.Messages = l.Messages[len(l.Messages)-l.MaxLines:]
	}
}

func (l *MessageLog) GetRecent(count int) []LogMessage {
	if len(l.Messages) <= count {
		return l.Messages
	}
	return l.Messages[len(l.Messages)-count:]
}
