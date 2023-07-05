package logger

import "go.uber.org/zap"

type ZapLooger struct {
	l *zap.Logger
}

func NewZapLogger() (*ZapLooger, error) {
	zaplogger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	zaplogger.Sync()
	return &ZapLooger{zaplogger}, nil
}

func (zl *ZapLooger) Info(s string) {
	zl.l.Info(s)
}

func (zl *ZapLooger) Error(s string) {
	zl.l.Error(s)
}

func (zl *ZapLooger) Panic(s string) {
	zl.l.Panic(s)
}
