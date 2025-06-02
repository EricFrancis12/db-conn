package pkg

import "log"

type Logger interface {
	Printf(format string, v ...any)
	Println(v ...any)
}

func NewLogger() Logger {
	return &logger{}
}

type logger struct{}

func (l logger) Printf(format string, v ...any) {
	log.Printf(format, v...)
}

func (l logger) Println(v ...any) {
	log.Println(v...)
}

type mockLogger struct{}

func (mockLogger) Printf(_ string, _ ...any) {}

func (mockLogger) Println(_ ...any) {}
