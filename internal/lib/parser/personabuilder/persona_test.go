package personabuilder

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vipulvpatil/candidate-tracker-go/internal/clients/openai"
)

func Test_BuildSuccess(t *testing.T) {

	// tests := []struct {
	// 	name           string
	// 	input          string
	// 	expectedOutput *Persona
	// }{
	// 	{
	// 		name: "sample response 1",
	// 		input: `{
	// 			"Name": "Vipul Patil",
	// 			"Email": "vipulvpatil@gmail.com",
	// 			"Phone": "",
	// 			"City": "",
	// 			"State": "",
	// 			"Country": "",
	// 			"YoE": 9,
	// 			"Tech Skills": ["Go", "NodeJS", "Ruby on Rails", "Xcode", "ObjC"],
	// 			"Soft Skills": [],
	// 			"Recommended Roles": ["Senior Software Engineer", "Backend Developer", "iOS Developer"],
	// 			"Education": [
	// 				{
	// 					"Qualification": "B.E., Computer Engineering",
	// 					"CompletionYear": "2008",
	// 					"Institute": "Ramrao Adik Institute of Technology"
	// 				},
	// 				{
	// 					"Qualification": "HSC, Science",
	// 					"CompletionYear": "2004",
	// 					"Institute": "Ramnivas Ruia Junior College"
	// 				},
	// 				{
	// 					"Qualification": "SSC",
	// 					"CompletionYear": "2002",
	// 					"Institute": "St. John The Baptist High School"
	// 				}
	// 			]
	// 		}`,
	// 		expectedOutput: &Persona{
	// 			Name:             "Vipul Patil",
	// 			Email:            "vipulvpatil@gmail.com",
	// 			Phone:            "",
	// 			City:             "",
	// 			State:            "",
	// 			Country:          "",
	// 			YoE:              9,
	// 			TechSkills:       []string{"Go", "NodeJS", "Ruby on Rails", "Xcode", "ObjC"},
	// 			SoftSkills:       []string{},
	// 			RecommendedRoles: []string{"Senior Software Engineer", "Backend Developer", "iOS Developer"},
	// 			Education: []Education{
	// 				{
	// 					Institute:      "Ramrao Adik Institute of Technology",
	// 					Qualification:  "B.E., Computer Engineering",
	// 					CompletionYear: "2008",
	// 				},
	// 				{
	// 					Institute:      "Ramnivas Ruia Junior College",
	// 					Qualification:  "HSC, Science",
	// 					CompletionYear: "2004",
	// 				},
	// 				{
	// 					Institute:      "St. John The Baptist High School",
	// 					Qualification:  "SSC",
	// 					CompletionYear: "2002",
	// 				},
	// 			},
	// 			Certifications: nil,
	// 			BuilderVersion: "1.0.0",
	// 		},
	// 	},
	// }

	t.Run("run examples", func(t *testing.T) {
		fileContent, err := os.ReadFile("persona_testcase_inputs.txt")
		assert.NoError(t, err)

		testInputs := strings.Split(string(fileContent), "***")
		testExpectedOutputs := []*Persona{
			{Name: "Person", Email: "someemail@example.com", Phone: "+91-1234567890", City: "", State: "", Country: "", YoE: 9, TechSkills: []string{"Go", "NodeJS", "Ruby on Rails", "Xcode", "ObjC"}, SoftSkills: []string{}, RecommendedRoles: []string{"Senior Software Engineer", "Backend Developer", "iOS Developer"}, Education: []Education{{Institute: "Ramrao Adik Institute of Technology", Qualification: "B.E., Computer Engineering", CompletionYear: "2008"}, {Institute: "Ramnivas Ruia Junior College", Qualification: "HSC, Science", CompletionYear: "2004"}, {Institute: "St. John The Baptist High School", Qualification: "SSC", CompletionYear: "2002"}}, Certifications: nil, BuilderVersion: "1.0.0"},
			{Name: "Person", Email: "someemail@example.com", Phone: "+91-1234567890", City: "", State: "", Country: "", YoE: 7, TechSkills: []string{"React JS", "React Native", "TypeScript", "Next JS", "Graphql"}, SoftSkills: []string{"Leadership", "Collaboration", "Problem-solving", "Communication", "Creativity"}, RecommendedRoles: []string{"Senior Software Engineer", "Product Engineer", "Lead Front End Developer"}, Education: []Education{{Institute: "Dehradun Institute of Technology, Dehradun", Qualification: "Bachelor of Technology", CompletionYear: "May 2016"}}, Certifications: []string(nil), BuilderVersion: "1.0.0"},
			{Name: "Person", Email: "someemail@example.com", Phone: "+91-1234567890", City: "San Francisco", State: "California", Country: "United States", YoE: 10, TechSkills: []string{"AWS", "JavaScript", "NodeJS", "ReactJS", "Java"}, SoftSkills: []string{"Team Management", "Technical Leadership", "Hiring", "Sprint Planning", "Scrum"}, RecommendedRoles: []string{"Software Engineering Manager", "Engineering Manager", "Senior Software Engineer"}, Education: []Education{{Institute: "San Francisco State University", Qualification: "Masters in Computer Science", CompletionYear: "December 2010"}, {Institute: "University of Mumbai, India", Qualification: "Bachelors in Computer Engineering", CompletionYear: "July 2008"}}, Certifications: []string(nil), BuilderVersion: "1.0.0"},
			{Name: "Person", Email: "someemail@example.com", Phone: "+91-1234567890", City: "Mumbai", State: "Maharashtra", Country: "India", YoE: 14, TechSkills: []string{"Delphi", "C#", "Elixir", "Azure DevOps", "JavaScript"}, SoftSkills: []string{"Team player", "Problem-solving", "Communication", "Leadership", "Time management"}, RecommendedRoles: []string{"Software Engineer", "Solution Architect", "Technical Specialist"}, Education: []Education{{Institute: "Konkan Gyanpeeth College Of Engineering", Qualification: "Bachelor's Degree in Information Technology (Data Warehousing and Mining)", CompletionYear: "2008"}, {Institute: "Distance Open Learning (Mumbai University)", Qualification: "Master in Information Technology", CompletionYear: "Ongoing"}}, Certifications: []string(nil), BuilderVersion: "1.0.0"},
			{Name: "Person", Email: "someemail@example.com", Phone: "+91-1234567890", City: "Mumbai", State: "Maharashtra", Country: "India", YoE: 11, TechSkills: []string{"C#", ".NET", "ASP.Net MVC", "Angular 9", "Web API"}, SoftSkills: []string{"Good Communication", "Interpersonal Skills", "Flexible", "Self-Motivated", "Team Player"}, RecommendedRoles: []string{"Sr. Technical Lead", "Sr. Manager", "Sr. Software Engineer"}, Education: []Education{{Institute: "Ramrao Adik Institute of Technology (R.A.I.T)", Qualification: "B. E. Computer Engineering", CompletionYear: "June 2011"}, {Institute: "Shreeram Polytechnic", Qualification: "Diploma in Computer Technology", CompletionYear: "June 2008"}}, Certifications: []string(nil), BuilderVersion: "1.0.0"},
			{Name: "Person", Email: "someemail@example.com", Phone: "+91-1234567890", City: "", State: "", Country: "", YoE: 14, TechSkills: []string{"Account Management", "Business Development", "Market Analysis", "Client Relationship Management", "Advertising Strategies"}, SoftSkills: []string{"Communication", "Leadership", "Negotiation", "Team Management", "Problem Solving"}, RecommendedRoles: []string{"Business Development Manager", "Key Accounts Manager", "Senior Manager, Marketing"}, Education: []Education{{Institute: "VIT (Wadala)", Qualification: "MMS in Marketing Management", CompletionYear: "2011"}, {Institute: "Vivekanand Education Society’s Institute of Technology", Qualification: "B.E. in Engineering", CompletionYear: "2008"}, {Institute: "St. John the Baptist Junior College", Qualification: "H.S.C", CompletionYear: ""}, {Institute: "St. John the Baptist High School", Qualification: "SSC", CompletionYear: ""}}, Certifications: []string(nil), BuilderVersion: "1.0.0"},
			{Name: "Person", Email: "someemail@example.com", Phone: "+91-1234567890", City: "Berkeley", State: "California", Country: "United States", YoE: 10, TechSkills: []string{"Evaluation and assessment", "Therapeutic physical therapy interventions", "Documentation in EMR", "Skilled therapeutic exercises", "Electrotherapy treatments"}, SoftSkills: []string{"Team player", "Management", "Communication", "Flexibility", "Multidisciplinary collaboration"}, RecommendedRoles: []string{"Cardiopulmonary Physical Therapist", "Orthopedic Physical Therapist", "Rehabilitation Program Director"}, Education: []Education{{Institute: "Arcadia University, Pennsylvania", Qualification: "Doctorate of Physical Therapy", CompletionYear: "2021"}, {Institute: "San Francisco State University, San Francisco", Qualification: "MSc Kinesiology", CompletionYear: "2011"}, {Institute: "Laxmi Memorial College, India", Qualification: "BSc Physiotherapy", CompletionYear: "2008"}}, Certifications: []string{"CPR certified", "Spine, Elbow & Ankle Manipulation, C Institute, India", "Diploma in Yoga Therapy", "Balance training and Swiss ball Thera TheraAcademy"}, BuilderVersion: "1.0.0"},
			{Name: "Person", Email: "someemail@example.com", Phone: "+91-1234567890", City: "Thane", State: "Maharashtra", Country: "India", YoE: 10, TechSkills: []string{"AutoCAD 2D", "MS Office", "SketchUp 3D", "Photoshop", "Revit"}, SoftSkills: []string{"Communication", "Project Management", "Team Collaboration", "Problem Solving", "Attention to Detail"}, RecommendedRoles: []string{"Senior Architect", "Junior Architect", "Intern"}, Education: []Education{{Institute: "Bharati Vidyapeeth College Of Architecture", Qualification: "Bachelor of Architecture", CompletionYear: "2012"}, {Institute: "South Indian Education Society", Qualification: "Higher Secondary Certificate", CompletionYear: "2004"}, {Institute: "St.John The Baptist High School", Qualification: "Secondary School Certificate", CompletionYear: "2002"}}, Certifications: []string(nil), BuilderVersion: "1.0.0"},
			{Name: "Person", Email: "someemail@example.com", Phone: "+91-1234567890", City: "Navi Mumbai", State: "Maharashtra", Country: "India", YoE: 10, TechSkills: []string{"Financial Accounting", "Taxation", "Planning", "Reporting", "ERP"}, SoftSkills: []string{"Strategic Thinking", "Problem Solving", "Data Collection", "Communication", "Analytical Skills"}, RecommendedRoles: []string{"Financial Accountant", "Tax Consultant", "ERP Specialist"}, Education: []Education{{Institute: "ICWAI", Qualification: "CMA", CompletionYear: "2010"}, {Institute: "Pune University", Qualification: "M.Com (Costing)", CompletionYear: "2011"}, {Institute: "Mumbai University", Qualification: "B.Com (Accounts and Finance)", CompletionYear: "2005"}}, Certifications: []string(nil), BuilderVersion: "1.0.0"},
			{Name: "Person", Email: "someemail@example.com", Phone: "+91-1234567890", City: "Bangalore", State: "Karnataka", Country: "India", YoE: 11, TechSkills: []string{"Brand Development", "Marketing Strategy", "Product Launch", "Digital Marketing", "Content Development"}, SoftSkills: []string{"Leadership", "Team Management", "Strategic Thinking", "Communication", "Analytical Skills"}, RecommendedRoles: []string{"Chief Marketing Officer", "Brand Manager", "Senior Marketing Manager"}, Education: []Education{{Institute: "Mudra Institute of Communications, Ahmedabad", Qualification: "PGDM(C)", CompletionYear: "2013"}, {Institute: "Mumbai University", Qualification: "B.E. (Electronics)", CompletionYear: "2009"}, {Institute: "S.I.E.S. College of Arts, Science and Commerce, Mumbai", Qualification: "Class XII (Maharashtra State Board)", CompletionYear: "2004"}, {Institute: "St. John the Baptist High School, Thane", Qualification: "Class X (Maharashtra State Board)", CompletionYear: "2002"}}, Certifications: []string(nil), BuilderVersion: "1.0.0"},
		}

		for i, testInput := range testInputs {
			openAiMockClient := openai.MockClientSuccess{
				Text: testInput,
			}
			persona, err := Build(testInput, &openAiMockClient)
			assert.NoError(t, err)

			var expectedOutput *Persona = nil

			if i < len(testExpectedOutputs) {
				expectedOutput = testExpectedOutputs[i]
			}

			assert.Equal(t, expectedOutput, persona)
		}
	})
}
