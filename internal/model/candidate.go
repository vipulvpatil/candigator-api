package model

import (
	"time"

	"github.com/pkg/errors"
	"github.com/vipulvpatil/candidate-tracker-go/internal/utilities"
)

type Candidate struct {
	id                     string
	createdAt              time.Time
	aiGeneratedPersona     *Persona
	manuallyCreatedPersona *Persona
}

type CandidateOptions struct {
	Id                     string
	CreatedAt              time.Time
	AiGeneratedPersona     *Persona
	ManuallyCreatedPersona *Persona
}

func NewCandidate(opts CandidateOptions) (*Candidate, error) {
	if utilities.IsBlank(opts.Id) {
		return nil, errors.New("cannot create candidate with an empty id")
	}

	var aiGeneratedPersona, manuallyCreatedPersona *Persona

	if opts.AiGeneratedPersona.IsValid() {
		aiGeneratedPersona = opts.AiGeneratedPersona
	}
	if opts.ManuallyCreatedPersona.IsValid() {
		manuallyCreatedPersona = opts.ManuallyCreatedPersona
	}

	if aiGeneratedPersona == nil && manuallyCreatedPersona == nil {
		return nil, errors.New("cannot create candidate without a valid persona")
	}

	candidate := Candidate{
		id:                     opts.Id,
		createdAt:              opts.CreatedAt,
		aiGeneratedPersona:     aiGeneratedPersona,
		manuallyCreatedPersona: manuallyCreatedPersona,
	}
	return &candidate, nil
}
