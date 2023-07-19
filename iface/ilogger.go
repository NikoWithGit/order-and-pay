package iface

type Ilogger interface {
	Info(string)
	Error(string)
	Panic(string)
	Infof(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Panicf(template string, args ...interface{})
}
