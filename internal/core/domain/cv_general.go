package domain

// --- Data Structures for YAML Unmarshaling ---
type PersonalInfo struct {
	Name           string `yaml:"name"`
	Email          string `yaml:"email"`
	Phone          string `yaml:"phone"`
	Website        string `yaml:"website"`
	Linkedin       string `yaml:"linkedin"`
	LinkedinHandle string `yaml:"linkedin_handle"`
	Github         string `yaml:"github"`
	GithubHandle   string `yaml:"github_handle"`
	ProfilePic     string `yaml:"profile_pic"`
	Location       string `yaml:"location"`
}

type Education struct {
	Institution string   `yaml:"institution"`
	Degree      string   `yaml:"degree"`
	Dates       string   `yaml:"dates"`
	Gpa         string   `yaml:"gpa"`
	Details     []string `yaml:"details"`
}

type Award struct {
	Title        string   `yaml:"title"`
	Organization string   `yaml:"organization"`
	Date         string   `yaml:"date"`
	Points       []string `yaml:"points"`
}

type Certificate struct {
	Name string `yaml:"name"`
	Year string `yaml:"year"`
}

type SkillGroup struct {
	Title string   `yaml:"title"`
	Items []string `yaml:"items"`
}

type CVGeneral struct {
	PersonalInfo      PersonalInfo  `yaml:"personal_info"`
	Description       string        `yaml:"description"`
	ResearchInterests []string      `yaml:"research_interests"`
	Skills            []SkillGroup  `yaml:"skills"`
	Certifications    []Certificate `yaml:"certificates"`
	Educations        []Education   `yaml:"education"`
	Awards            []Award       `yaml:"awards"`
}
