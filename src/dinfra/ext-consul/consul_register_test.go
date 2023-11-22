package extconsul

import (
	"testing"

	extlogrus "github.com/AlgerDu/go-dream/src/dinfra/ext-logrus"
)

func newTestRegister() *ConsulRegister {

	register, _ := New(extlogrus.New(nil), &ConsulRegisterOptions{
		Address: "localhost:18500",
		Token:   "123456",
		Env:     "test",
	})

	return register
}

func TestGetService(t *testing.T) {
	register := newTestRegister()
	service, err := register.Get("testService")

	if err != nil {
		t.Errorf("%v", err)
	}

	if service.ID == "" {
		t.Error("empty id")
	}
}
