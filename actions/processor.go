package actions

//go:generate counterfeiter . Processor

type Processor interface {
	Process(config string) string
}
