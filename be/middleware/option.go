package middleware

type Option struct {
	noLogin bool
}

func DefaultOption() Option {
	return Option{}
}

func (o Option) WithNoLogin() Option {
	o.noLogin = true
	return o
}

func (o Option) WithLogin() Option {
	o.noLogin = false
	return o
}

func (o Option) IsNoLogin() bool {
	return o.noLogin
}
