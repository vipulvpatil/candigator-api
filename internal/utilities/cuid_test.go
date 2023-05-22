package utilities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Generate(t *testing.T) {
	generator := &RandomIdGenerator{}
	result := generator.Generate()
	assert.NotEmpty(t, result)
}
