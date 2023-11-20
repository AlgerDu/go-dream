package dinfra

type (
	App interface {
		Run() error
	}

	AppBuilder interface {
		Build() (App, error)
	}
)
