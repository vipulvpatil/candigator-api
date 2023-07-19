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

	t.Run("returns False for manually generated personas with on name", func(t *testing.T) {
		persona := &Persona{}
		assert.False(t, persona.IsValid())
	})

	t.Run("returns False for AI generated personas if non file upload Id exists", func(t *testing.T) {
		persona := &Persona{
			Name:    "awesome persona",
			BuiltBy: "AI",
		}
		assert.False(t, persona.IsValid())
	})

	t.Run("returns True for AI generated personas", func(t *testing.T) {
		persona := &Persona{
			Name:         "awesome persona",
			BuiltBy:      "AI",
			FileUploadId: "fp_id1",
		}
		assert.True(t, persona.IsValid())
	})
}
