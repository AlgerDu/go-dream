package dinfra

type (
	Result struct {
		Code int    `json:"code"`
		Data any    `json:"data,omitempty"`
		Msg  string `json:"msg,omitempty"`
	}
)
