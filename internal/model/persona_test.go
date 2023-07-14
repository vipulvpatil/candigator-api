package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Persona_IsValid(t *testing.T) {
	t.Run("returns True", func(t *testing.T) {
		persona := &Persona{
			Name: "awesome persona",
		}
		assert.True(t, persona.IsValid())
	})

	t.Run("returns False", func(t *testing.T) {
		persona := &Persona{
			Id: "some_id",
		}
		assert.False(t, persona.IsValid())
	})
}
