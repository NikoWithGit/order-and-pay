package intrface

type Ilogger interface {
	Info(string)
	Error(string)
	Panic(string)
}
