package internals

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

type Experience struct {
	Role     string   `yaml:"role"`
	Type     string   `yaml:"type"`
	Company  string   `yaml:"company"`
	Location string   `yaml:"location"`
	Dates    string   `yaml:"dates"`
	Points   []string `yaml:"points"`
	Types    []string `yaml:"types"`
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

type Project struct {
	Name         string   `yaml:"name"`
	Github       string   `yaml:"github"`
	GithubHandle string   `yaml:"github_handle"`
	Points       []string `yaml:"points"`
	Types        []string `yaml:"types"`
}

type Certificate struct {
	Name string `yaml:"name"`
	Year int    `yaml:"year"`
}

type SkillsAchievements struct {
	HardSkills           []string      `yaml:"Hard Skills"`
	SoftSkills           []string      `yaml:"Soft Skills"`
	ProgrammingLanguages []string      `yaml:"Programming Languages"`
	DatabaseLanguages    []string      `yaml:"Database Languages"`
	Misc                 []string      `yaml:"Misc"`
	Certificates         []Certificate `yaml:"Certificates"`
}

type CVData struct {
	PersonalInfo       PersonalInfo       `yaml:"personal_info"`
	Experiences        []Experience       `yaml:"experiences"`
	Education          []Education        `yaml:"education"`
	Awards             []Award            `yaml:"awards"`
	Projects           []Project          `yaml:"projects"`
	SkillsAchievements SkillsAchievements `yaml:"skills_achievements"`
}
