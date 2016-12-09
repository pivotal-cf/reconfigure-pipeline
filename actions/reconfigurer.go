package actions

//go:generate counterfeiter . Reconfigurer

type Reconfigurer interface {
	Reconfigure(target, pipeline, configPath, variablesPath string) error
}
