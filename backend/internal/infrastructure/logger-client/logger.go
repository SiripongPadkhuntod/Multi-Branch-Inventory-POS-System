package loggerclient

import (
	"context"
	"log"
)

type Logger struct{}

func InitGlobalLogger(_, _ bool) {}

func WithContext(ctx context.Context, _ ...any) context.Context {
	return ctx
}

func FromContext(context.Context) *Logger {
	return &Logger{}
}

func GetLogger(context.Context, ...any) *Logger {
	return &Logger{}
}

func (l *Logger) Infow(message string, keysAndValues ...any) {
	log.Println(append([]any{message}, keysAndValues...)...)
}

func (l *Logger) Errorw(message string, keysAndValues ...any) {
	log.Println(append([]any{message}, keysAndValues...)...)
}

func (l *Logger) Sync() {}
