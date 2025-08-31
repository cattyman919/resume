package model

Experience :: struct {
	Company:  string `yaml: "company"`,
	Location: string `yaml: "location"`,
	Role:     string `yaml: "role"`,
	Dates:    string `yaml: "dates"`,
	Job_Type: string `yaml: "job_type"`,
	CV_Types: [dynamic]string `yaml: "cv_type"`,
	Points:   [dynamic]string `yaml: "points"`,
}

make_experience :: proc() -> Experience {
	return Experience{CV_Types = make([dynamic]string), Points = make([dynamic]string)}
}

