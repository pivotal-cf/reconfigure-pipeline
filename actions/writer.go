package actions

//go:generate counterfeiter . Writer

type Writer interface {
	Write(content string) (string, error)
}
