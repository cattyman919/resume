package resume

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
	JobType  string   `yaml:"job_type"`
	Company  string   `yaml:"company"`
	Location string   `yaml:"location"`
	Dates    string   `yaml:"dates"`
	Points   []string `yaml:"points"`
	CVType   []string `yaml:"cv_type"`
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
	CVType       []string `yaml:"cv_type"`
}

type Certificate struct {
	Name string `yaml:"name"`
	Year string `yaml:"year"`
}

type SkillsAchievements struct {
	HardSkills           []string      `yaml:"Hard Skills"`
	SoftSkills           []string      `yaml:"Soft Skills"`
	ProgrammingLanguages []string      `yaml:"Programming Languages"`
	DatabaseLanguages    []string      `yaml:"Database Languages"`
	Misc                 []string      `yaml:"Misc"`
	Certificates         []Certificate `yaml:"Certificates"`
}

type CV_General struct {
	PersonalInfo       PersonalInfo       `yaml:"personal_info"`
	SkillsAchievements SkillsAchievements `yaml:"skills_achievements"`
	Education          []Education        `yaml:"education"`
	Awards             []Award            `yaml:"awards"`
}

type CVData struct {
	General     *CV_General
	Projects    []Project
	Experiences []Experience
}
