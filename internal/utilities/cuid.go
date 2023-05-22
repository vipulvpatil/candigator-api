package utilities

import "github.com/lucsky/cuid"

type CuidGenerator interface {
	Generate() string
}

type RandomIdGenerator struct{}

func (r *RandomIdGenerator) Generate() string {
	return cuid.New()
}
