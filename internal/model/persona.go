package model

type Persona struct {
	Id               string
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
}

type Education struct {
	Id             string
	Institute      string `json:"Institute"`
	Qualification  string `json:"Qualification"`
	CompletionYear string `json:"CompletionYear"`
}
