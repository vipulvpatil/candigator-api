package model

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/vipulvpatil/candidate-tracker-go/internal/utilities"
)

type Persona struct {
	Name             string      `json:"Name"`
	Email            string      `json:"Email"`
	Phone            string      `json:"Phone"`
	City             string      `json:"City"`
	State            string      `json:"State"`
	Country          string      `json:"Country"`
	YoE              int         `json:"YoE"`
	TechSkills       []string    `json:"Tech Skills"`
	SoftSkills       []string    `json:"Soft Skills"`
	RecommendedRoles []string    `json:"Recommended Roles"`
	Education        []Education `json:"Education"`
	Certifications   []string    `json:"Certifications"`
	BuilderVersion   string
	BuiltBy          string
	FileUploadId     string
}

type Education struct {
	Institute      string `json:"Institute"`
	Qualification  string `json:"Qualification"`
	CompletionYear string `json:"CompletionYear"`
}

func (p *Persona) IsValid() bool {
	if p == nil {
		return false
	}

	if utilities.IsBlank(p.Name) {
		return false
	}

	if p.BuiltBy == "AI" {
		return !utilities.IsBlank(p.FileUploadId)
	}

	return true
}

func (p *Persona) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *Persona) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &p)
}
