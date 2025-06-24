package es

type Term struct {
	Name  string
	Value interface{}
}

type Terms struct {
	Name   string
	Values []interface{}
}

type BenTerms struct {
	Name   string
	Values []interface{}
}

type Range struct {
	Name string
	Gte  *uint64 `json:"gte,omitempty"`
	Lte  *uint64 `json:"lte,omitempty"`
	Gt   *uint64 `json:"gt,omitempty"`
	Lt   *uint64 `json:"lt,omitempty"`
}
