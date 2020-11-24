package model

type TokenWA interface {
	GetTokenWA() string
	GetExpired() float64
	SetTokenWA(token string)
	SetExpired(exp float64)
}

type OTP struct {
	Code    string
	Expired float64
}

func (o OTP) GetTokenWA() string {
	return o.Code
}

func (e OTP) GetExpired() float64 {
	return e.Expired
}

func (o *OTP) SetTokenWA(code string) {
	o.Code = code
}

func (e *OTP) SetExpired(expired float64) {
	e.Expired = expired
}
