package dinfra

type (
	IdProvider interface {
		Provide() ID
	}
)
