package model

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/vipulvpatil/candidate-tracker-go/internal/utilities"
)

type Persona struct {
	Name             string       `json:"Name,omitempty"`
	Email            string       `json:"Email,omitempty"`
	Phone            string       `json:"Phone,omitempty"`
	City             string       `json:"City,omitempty"`
	State            string       `json:"State,omitempty"`
	Country          string       `json:"Country,omitempty"`
	YoE              int          `json:"YoE,omitempty"`
	TechSkills       []string     `json:"Tech Skills,omitempty"`
	SoftSkills       []string     `json:"Soft Skills,omitempty"`
	RecommendedRoles []string     `json:"Recommended Roles,omitempty"`
	Education        []Education  `json:"Education,omitempty"`
	Experience       []Experience `json:"Experience,omitempty"`
	Certifications   []string     `json:"Certifications,omitempty"`
	BuilderVersion   string       `json:"BuilderVersion,omitempty"`
	BuiltBy          string       `json:"BuiltBy,omitempty"`
	FileUploadId     string       `json:"FileUploadId,omitempty"`
}

type Education struct {
	Institute      string `json:"Institute,omitempty"`
	Qualification  string `json:"Qualification,omitempty"`
	CompletionYear string `json:"CompletionYear,omitempty"`
}

type Experience struct {
	Title        string `json:"Title,omitempty"`
	CompanyName  string `json:"Company Name,omitempty"`
	StartingYear string `json:"Starting Year,omitempty"`
	EndingYear   string `json:"Ending Year,omitempty"`
	Ongoing      bool   `json:"Ongoing,omitempty"`
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

func (p *Persona) IsEqual(other *Persona) bool {
	if p == nil {
		return other == nil
	}
	if other == nil {
		return p == nil
	}
	return p.Name == other.Name &&
		p.Email == other.Email &&
		p.Phone == other.Phone &&
		p.City == other.City &&
		p.State == other.State &&
		p.Country == other.Country &&
		p.YoE == other.YoE &&
		EqualPersonaAttributeArray(p.TechSkills, other.TechSkills) &&
		EqualPersonaAttributeArray(p.SoftSkills, other.SoftSkills) &&
		EqualPersonaAttributeArray(p.RecommendedRoles, other.RecommendedRoles) &&
		EqualPersonaAttributeArray(p.Education, other.Education) &&
		EqualPersonaAttributeArray(p.Experience, other.Experience) &&
		EqualPersonaAttributeArray(p.Certifications, other.Certifications) &&
		p.BuilderVersion == other.BuilderVersion &&
		p.BuiltBy == other.BuiltBy &&
		p.FileUploadId == other.FileUploadId
}

func EqualPersonaAttributeArray[A comparable](first, second []A) bool {
	if len(first) != len(second) {
		return false
	}
	for i := range first {
		if first[i] != second[i] {
			return false
		}
	}

	return true
}

func (p *Persona) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *Persona) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &p)
}
