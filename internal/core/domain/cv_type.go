package domain

type Layout string

const (
	LayoutExperiences  Layout = "experiences"
	LayoutDescription  Layout = "description"
	LayoutEducations   Layout = "educations"
	LayoutSkills       Layout = "skills"
	LayoutCertificates Layout = "certificates"
	LayoutAwards       Layout = "awards"
	LayoutProjects     Layout = "projects"
)

type Experience struct {
	Role     string   `yaml:"role"`
	JobType  string   `yaml:"job_type"`
	Company  string   `yaml:"company"`
	Location string   `yaml:"location"`
	Dates    string   `yaml:"dates"`
	Points   []string `yaml:"points"`
}

type Project struct {
	Name         string   `yaml:"name"`
	Github       string   `yaml:"github"`
	GithubHandle string   `yaml:"github_handle"`
	Points       []string `yaml:"points"`
}

type CVType struct {
	// Type name taken from the filename
	Type string
	// Will give a warning if the generated PDF has more than MaxPages
	MaxPages    *int         `yaml:"max_pages,omitempty"`
	Layouts     []Layout     `yaml:"layout"`
	Description string       `yaml:"description"`
	Projects    []Project    `yaml:"projects"`
	Experiences []Experience `yaml:"experiences"`
}
