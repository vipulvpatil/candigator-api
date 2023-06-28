package candidatebuilder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getCandidateDataFromOpenAiResponse(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedOutput *Candidate
		errorExpected  bool
		errorString    string
	}{
		// 		{
		// 			name: "sample response 1",
		// 			input: `
		// Name: Aamir Madan

		// Email: aamir.madan@gmail.com

		// Phone: +91 9876543210

		// City: Not mentioned

		// State: Not mentioned

		// Country: Not mentioned

		// Years of Experience: 7

		// Tech Skills: React JS, React Native, TypeScript, Next JS, Graphql

		// Soft Skills: Leadership, Team Management, Problem Solving, Creativity, Communication

		// Recommended Job Positions: Senior Front-end Developer, Technical Lead, Full Stack Developer`,
		// 			expectedOutput: &Candidate{
		// 				Name:            "Aamir Madan",
		// 				Email:           "aamir.madan@gmail.com",
		// 				Phone:           "+91 9876543210",
		// 				City:            "",
		// 				State:           "",
		// 				Country:         "",
		// 				YoE:             "7",
		// 				TechSkills:      []string{"React JS", "React Native", "TypeScript", "Next JS", "Graphql"},
		// 				SoftSkills:      []string{"Leadership", "Team Management", "Problem Solving", "Creativity", "Communication"},
		// 				RecommendedJobs: []string{"Senior Front-end Developer", "Technical Lead", "Full Stack Developer"},
		// 			},

		// 			errorExpected: false,
		// 			errorString:   "",
		// 		},
		// 		{
		// 			name: "sample response 2",
		// 			input: `
		// Full Name: Vijay Pant
		// Email: vijaypant@gmail.com
		// Phone: N/A
		// City: N/A
		// State: N/A
		// Country: N/A

		// Years of Experience: 9+

		// Tech Skills:
		// 1. Go
		// 2. NodeJS
		// 3. Ruby on Rails
		// 4. iOS (ObjC, Swift 4.0, Xcode, UIKit, Cocoa framework)
		// 5. Docker, AWS

		// Soft Skills:
		// 1. Leadership
		// 2. Problem-solving
		// 3. Teamwork
		// 4. Communication
		// 5. Mentoring

		// Recommended Job Positions:
		// 1. Senior Software Engineer
		// 2. Backend Developer
		// 3. iOS Developer`,
		// 			expectedOutput: &Candidate{
		// 				Name:            "Vijay Pant",
		// 				Email:           "vijaypant@gmail.com",
		// 				Phone:           "",
		// 				City:            "",
		// 				State:           "",
		// 				Country:         "",
		// 				YoE:             "7",
		// 				TechSkills:      []string{"React JS", "React Native", "TypeScript", "Next JS", "Graphql"},
		// 				SoftSkills:      []string{"Leadership", "Team Management", "Problem Solving", "Creativity", "Communication"},
		// 				RecommendedJobs: []string{"Senior Front-end Developer", "Technical Lead", "Full Stack Developer"},
		// 			},

		// 			errorExpected: false,
		// 			errorString:   "",
		// 		},
		// 		{
		// 			name: "sample response 3",
		// 			input: `
		// Name: Heman Universe
		// Email: heman.universe@gmail.com
		// Phone: +1-555-666-9999
		// Location: San Francisco, CA

		// Years of Experience: 10+ years

		// Tech Skills:
		// 1. AWS
		// 2. Java
		// 3. JavaScript
		// 4. ReactJS
		// 5. NodeJS

		// Soft Skills:
		// 1. Team Management
		// 2. Technical Leadership
		// 3. Hiring
		// 4. Sprint Planning
		// 5. Scrum

		// Recommended Job Positions:
		// 1. Software Engineering Manager
		// 2. Senior Software Engineer - Full Stack
		// 3. Engineering Manager`,
		// 			expectedOutput: &Candidate{
		// 				Name:            "Heman Universe",
		// 				Email:           "heman.universe@gmail.com",
		// 				Phone:           "+1-555-666-9999",
		// 				City:            "",
		// 				State:           "",
		// 				Country:         "",
		// 				YoE:             "7",
		// 				TechSkills:      []string{"React JS", "React Native", "TypeScript", "Next JS", "Graphql"},
		// 				SoftSkills:      []string{"Leadership", "Team Management", "Problem Solving", "Creativity", "Communication"},
		// 				RecommendedJobs: []string{"Senior Front-end Developer", "Technical Lead", "Full Stack Developer"},
		// 			},

		// 			errorExpected: false,
		// 			errorString:   "",
		// 		},
		{
			name: "sample response 4",
			input: `
Name: Sonal Yash Matter
Email: sonal555@gmail.com
Phone: +91-8812349099
City: Mumbai
State: Maharashtra
Country: India
Years of Experience: 11 years

Tech Skills:
1. C#
2. .NET
3. ASP.Net MVC
4. SQL Server
5. Angular

Soft Skills:
1. Good communication
2. Interpersonal skills
3. Adaptability
4. Self-motivated
5. Team player

Recommended Job Positions:
1. Senior Technical Lead
2. Senior Manager
3. Software Engineer`,
			expectedOutput: &Candidate{
				Name:            "Sonal Yash Matter",
				Email:           "sonal555@gmail.com",
				Phone:           "+91-8812349099",
				City:            "",
				State:           "",
				Country:         "",
				YoE:             "7",
				TechSkills:      []string{"React JS", "React Native", "TypeScript", "Next JS", "Graphql"},
				SoftSkills:      []string{"Leadership", "Team Management", "Problem Solving", "Creativity", "Communication"},
				RecommendedJobs: []string{"Senior Front-end Developer", "Technical Lead", "Full Stack Developer"},
			},

			errorExpected: false,
			errorString:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			candidate, err := getCandidateDataFromOpenAiResponse(tt.input)
			if !tt.errorExpected {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedOutput, candidate)
			} else {
				assert.NotEmpty(t, tt.errorString)
				assert.ErrorContains(t, err, tt.errorString)
			}
		})
	}
}
