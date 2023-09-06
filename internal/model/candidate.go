package model

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/vipulvpatil/candidate-tracker-go/internal/utilities"
)

type Candidate struct {
	id                     string
	createdAt              time.Time
	updatedAt              time.Time
	aiGeneratedPersona     *Persona
	manuallyCreatedPersona *Persona
	team                   *Team
	fileUploadId           string
}

type CandidateOptions struct {
	Id                     string
	CreatedAt              time.Time
	UpdatedAt              time.Time
	AiGeneratedPersona     *Persona
	ManuallyCreatedPersona *Persona
	Team                   *Team
	FileUploadId           string
}

func NewCandidate(opts CandidateOptions) (*Candidate, error) {
	if utilities.IsBlank(opts.Id) {
		return nil, errors.New("cannot create Candidate with an empty id")
	}

	if opts.Team == nil {
		return nil, errors.New("cannot create Candidate with a nil Team")
	}

	var aiGeneratedPersona, manuallyCreatedPersona *Persona

	if opts.AiGeneratedPersona.IsValid() {
		aiGeneratedPersona = opts.AiGeneratedPersona
	}
	if opts.ManuallyCreatedPersona.IsValid() {
		manuallyCreatedPersona = opts.ManuallyCreatedPersona
	}

	if aiGeneratedPersona == nil && manuallyCreatedPersona == nil {
		return nil, errors.New("cannot create Candidate without a valid persona")
	}

	candidate := Candidate{
		id:                     opts.Id,
		createdAt:              opts.CreatedAt,
		updatedAt:              opts.UpdatedAt,
		aiGeneratedPersona:     aiGeneratedPersona,
		manuallyCreatedPersona: manuallyCreatedPersona,
		team:                   opts.Team,
		fileUploadId:           opts.FileUploadId,
	}
	return &candidate, nil
}

func (c *Candidate) Id() string {
	return c.id
}

func (c *Candidate) AiGeneratedPersonaAsJsonString() string {
	if c.aiGeneratedPersona == nil {
		return ""
	}

	jsonString, err := json.Marshal(c.aiGeneratedPersona)
	if err != nil {
		return ""
	}

	return string(jsonString)
}

func (c *Candidate) ManuallyCreatedPersonaAsJsonString() string {
	if c.manuallyCreatedPersona == nil {
		return ""
	}

	jsonString, err := json.Marshal(c.manuallyCreatedPersona)
	if err != nil {
		return ""
	}

	return string(jsonString)
}

func (c *Candidate) FileUploadId() string {
	return c.fileUploadId
}

func (c *Candidate) IsEqual(other *Candidate) bool {
	fmt.Println(c.manuallyCreatedPersona)
	fmt.Println(other.manuallyCreatedPersona)
	return c.id == other.id &&
		c.aiGeneratedPersona.IsEqual(other.aiGeneratedPersona) &&
		c.manuallyCreatedPersona.IsEqual(other.manuallyCreatedPersona) &&
		c.team == other.team &&
		c.fileUploadId == other.fileUploadId
}
