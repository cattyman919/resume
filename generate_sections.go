package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// --- Configuration ---
const (
	dataFile          = "cv_data.json"
	mainCVSectionsDir = "main_cv/sections"
	bwCVSectionsDir   = "bw_cv/sections"
)

// --- Data Structures for JSON Unmarshaling ---

type PersonalInfo struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	Website        string `json:"website"`
	Linkedin       string `json:"linkedin"`
	LinkedinHandle string `json:"linkedin_handle"`
	Github         string `json:"github"`
	GithubHandle   string `json:"github_handle"`
	ProfilePic     string `json:"profile_pic"`
	Location       string `json:"location"`
}

type Experience struct {
	Role     string   `json:"role"`
	Type     string   `json:"type"`
	Company  string   `json:"company"`
	Location string   `json:"location"`
	Dates    string   `json:"dates"`
	Points   []string `json:"points"`
}

type Education struct {
	Institution string   `json:"institution"`
	Degree      string   `json:"degree"`
	Dates       string   `json:"dates"`
	Gpa         string   `json:"gpa"`
	Details     []string `json:"details"`
}

type Award struct {
	Title        string   `json:"title"`
	Organization string   `json:"organization"`
	Date         string   `json:"date"`
	Points       []string `json:"points"`
}

type Project struct {
	Name         string   `json:"name"`
	Github       string   `json:"github"`
	GithubHandle string   `json:"github_handle"`
	Points       []string `json:"points"`
}

type SkillsAchievements struct {
	HardSkills           []string `json:"Hard Skills"`
	SoftSkills           []string `json:"Soft Skills"`
	ProgrammingLanguages []string `json:"Programming Languages"`
	DatabaseLanguages    []string `json:"Database Languages"`
	Misc                 []string `json:"Misc"`
	Certificates         []struct {
		Name string `json:"name"`
		Year uint16 `json:"year"`
	} `json:"Certificates"`
}

type CVData struct {
	PersonalInfo       PersonalInfo       `json:"personal_info"`
	Experiences        []Experience       `json:"experiences"`
	Education          []Education        `json:"education"`
	Awards             []Award            `json:"awards"`
	Projects           []Project          `json:"projects"`
	SkillsAchievements SkillsAchievements `json:"skills_achievements"`
}

// --- Helper Functions ---

var latexEscaper = strings.NewReplacer(
	`&`, `\&`,
	`%`, `\%`,
	`$`, `\$`,
	`#`, `\#`,
	`_`, `\_`,
	`{`, `\{`,
	`}`, `\}`,
	`~`, `\textasciitilde{}`,
	`^`, `\textasciicircum{}`,
	`\`, `\textbackslash{}`,
)

func escapeLatex(text string) string {
	return latexEscaper.Replace(text)
}

var boldRegex = regexp.MustCompile(`\*\*(.*?)\*\*`)

func formatLatexBold(text string) string {
	escapedText := escapeLatex(text)
	return boldRegex.ReplaceAllString(escapedText, `\textbf{$1}`)
}

func writeTexFile(path, content string) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("Error creating directory %s: %v", dir, err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		log.Fatalf("Error writing %s: %v", path, err)
	}
}

func generateHighlights(points []string) string {
	var builder strings.Builder
	builder.WriteString("\\begin{highlights}\n")
	for _, point := range points {
		builder.WriteString(fmt.Sprintf("    \\item %s\n", formatLatexBold(point)))
	}
	builder.WriteString("\\end{highlights}")
	return builder.String()
}

// --- Section Generators ---

// == Header ==
func generateHeaderMainCV(data CVData) string {
	info := data.PersonalInfo
	return fmt.Sprintf(`\documentclass[../main.tex]{subfiles}
\begin{document}
\begin{header}
    \setcolumnwidth{3.5cm, \fill}
    \columnsep=0pt
    \raggedright
    \noindent %% Prevent paragraph indentation
    \begin{paracol}{2}
        \begin{leftcolumn}
            \vspace*{\fill} %% Vertical alignment
            \includegraphics[width=3cm]{%s}
            \vspace*{\fill} %% Vertical alignment
        \end{leftcolumn}
        \begin{rightcolumn}
            \textbf{\fontsize{18 pt}{24 pt}\selectfont \color{NavyBlue} %s} \\
            \normalsize
            \kern 0.25 cm%%
            \mbox{\hrefWithoutArrow{mailto:%s}{\color{black}{\footnotesize\faEnvelope[regular]}\hspace*{0.13cm}%s}}%%
            \kern 0.25 cm%%
            \AND%%
            \kern 0.25 cm%%
            \mbox{\hrefWithoutArrow{tel:%s}{\color{black}{\footnotesize\faPhone*}\hspace*{0.13cm}%s}}%%
            \kern 0.25 cm%%
            \AND%%
            \kern 0.25 cm%%
            \mbox{\hrefWithoutArrow{%s}{\color{black}{\footnotesize\faLink}\hspace*{0.13cm}%s}}%%
            \kern 0.25 cm%%
            \AND%%
            \kern 0.25 cm%%
            \mbox{\hrefWithoutArrow{%s}{\color{black}{\footnotesize\faLinkedinIn}\hspace*{0.13cm}%s}}%%
            \kern 0.25 cm%%
            \AND%%
            \kern 0.25 cm%%
            \mbox{\hrefWithoutArrow{%s}{\color{black}{\footnotesize\faGithub}\hspace*{0.13cm}%s}}%%
        \end{rightcolumn}
    \end{paracol}
\end{header}
\end{document}
`, escapeLatex(info.ProfilePic), escapeLatex(info.Name), escapeLatex(info.Email), escapeLatex(info.Email), escapeLatex(info.Phone), escapeLatex(info.Phone), escapeLatex(info.Website), escapeLatex(strings.Replace(info.Website, "https://", "", 1)), escapeLatex(info.Linkedin), escapeLatex(info.LinkedinHandle), escapeLatex(info.Github), escapeLatex(info.GithubHandle))
}

func generateHeaderBwCV(data CVData) string {
	info := data.PersonalInfo
	return fmt.Sprintf(`\documentclass[../main.tex]{subfiles}
\begin{document}
\begin{header}
    \fontsize{25 pt}{25 pt}\selectfont %s

    \vspace{5 pt}

    \normalsize
    \mbox{\color{black}{\footnotesize\faHome}\hspace*{0.13cm}%s}%%
    \kern 5.0 pt%%
    \AND%%
    \kern 5.0 pt%%
    \mbox{\hrefWithoutArrow{mailto:%s}{\color{black}{\footnotesize\faEnvelope[regular]}\hspace*{0.13cm}%s}}%%
    \kern 5.0 pt%%
    \AND%%
    \kern 5.0 pt%%
    \mbox{\hrefWithoutArrow{tel:%s}{\color{black}{\footnotesize\faPhone*}\hspace*{0.13cm}%s}}
    \kern 5.0 pt%%
    \AND%%
    \kern 5.0 pt%%
    \mbox{\hrefWithoutArrow{%s}{\color{black}{\footnotesize\faLink}\hspace*{0.13cm}%s}}%%
    \kern 5.0 pt%%
    \AND%%
    \kern 5.0 pt%%
    \mbox{\hrefWithoutArrow{%s}{\color{black}{\footnotesize\faLinkedinIn}\hspace*{0.13cm}%s}}%%
    \kern 5.0 pt%%
    \AND%%
    \kern 5.0 pt%%
    \mbox{\hrefWithoutArrow{%s}{\color{black}{\footnotesize\faGithub}\hspace*{0.13cm}%s}}%%
\end{header}
\end{document}
`, escapeLatex(info.Name), escapeLatex(info.Location), escapeLatex(info.Email), escapeLatex(info.Email), escapeLatex(info.Phone), escapeLatex(strings.Replace(info.Phone, "+62-", "0", 1)), escapeLatex(info.Website), escapeLatex(strings.Replace(info.Website, "https://", "", 1)), escapeLatex(info.Linkedin), escapeLatex(info.LinkedinHandle), escapeLatex(info.Github), escapeLatex(info.GithubHandle))
}

// == Experience ==
func generateExperienceMainCV(data CVData) string {
	var builder strings.Builder
	builder.WriteString("\\documentclass[../main.tex]{subfiles}\n\\begin{document}\n\\section{Experience}\n")
	for i, exp := range data.Experiences {
		roleType := fmt.Sprintf("%s (%s)", escapeLatex(exp.Role), escapeLatex(exp.Type))
		builder.WriteString(fmt.Sprintf(`
\begin{twocolentry}{%s}
    \textbf{%s \color{Gray}- %s}\\
    \textit{%s}
\end{twocolentry}

\vspace{0.10 cm}
\begin{onecolentry}
%s
\end{onecolentry}
`, escapeLatex(exp.Dates), escapeLatex(exp.Company), escapeLatex(exp.Location), roleType, generateHighlights(exp.Points)))
		if i < len(data.Experiences)-1 {
			builder.WriteString("\n\\vspace{0.40 cm}\n")
		}
	}
	builder.WriteString("\n\\end{document}")
	return builder.String()
}

func generateExperienceBwCV(data CVData) string {
	var builder strings.Builder
	builder.WriteString("\\documentclass[../main.tex]{subfiles}\n\\begin{document}\n\\section{\\sectiontitle[\\Large]{Experience}}\n")
	for i, exp := range data.Experiences {
		builder.WriteString(fmt.Sprintf(`
\begin{twocolentry}{%s}
    \sectiontitle{%s} \location{-- %s}\\
    \vspace{0.10 cm}
    \textit{%s}
\end{twocolentry}

\vspace{0.10 cm}
\begin{onecolentry}
%s
\end{onecolentry}
`, escapeLatex(exp.Dates), escapeLatex(exp.Company), escapeLatex(exp.Location), escapeLatex(exp.Role), generateHighlights(exp.Points)))
		if i < len(data.Experiences)-1 {
			builder.WriteString("\n\\vspace{0.4 cm}\n")
		}
	}
	builder.WriteString("\n\\end{document}")
	return builder.String()
}

// == Education ==
func generateEducationMainCV(data CVData) string {
	var builder strings.Builder
	builder.WriteString("\\documentclass[../main.tex]{subfiles}\n\\begin{document}\n\\section{Education}\n")
	for _, edu := range data.Education {
		builder.WriteString(fmt.Sprintf(`
\begin{twocolentry}{
    \small
    \textit{%s}
}
    \textbf{%s} \\
    \textit{%s, %s}
\end{twocolentry}
`, escapeLatex(edu.Dates), escapeLatex(edu.Institution), escapeLatex(edu.Degree), escapeLatex(edu.Gpa)))
	}
	builder.WriteString("\n\\end{document}")
	return builder.String()
}

func generateEducationBwCV(data CVData) string {
	var builder strings.Builder
	builder.WriteString("\\documentclass[../main.tex]{subfiles}\n\\begin{document}\n\\section{\\sectiontitle[\\Large]{Education}}\n")
	for _, edu := range data.Education {
		boldEduText := fmt.Sprintf(`\textbf{%s, %s}`, escapeLatex(edu.Institution), escapeLatex(edu.Degree))
		builder.WriteString(fmt.Sprintf(`
\begin{twocolentry}{
    %s
}
    %s
\end{twocolentry}
`, escapeLatex(edu.Dates), boldEduText))
		if len(edu.Details) > 0 {
			builder.WriteString(fmt.Sprintf("\n\\vspace{0.10 cm}\n\\begin{onecolentry}\n%s\n\\end{onecolentry}\n", generateHighlights(edu.Details)))
		}
	}
	builder.WriteString("\n\\end{document}")
	return builder.String()
}

// == Awards ==
func generateAwardsMainCV(data CVData) string {
	var builder strings.Builder
	builder.WriteString("\\documentclass[../main.tex]{subfiles}\n\\begin{document}\n\\section{Awards}\n")
	for i, award := range data.Awards {
		builder.WriteString(fmt.Sprintf(`
\begin{twocolentry}{
    \textit{%s}
}
    \textbf{%s}\\
    %s
\end{twocolentry}
\vspace{0.10 cm}
\begin{onecolentry}
%s
\end{onecolentry}
`, escapeLatex(award.Date), escapeLatex(award.Title), escapeLatex(award.Organization), generateHighlights(award.Points)))
		if i < len(data.Awards)-1 {
			builder.WriteString("\n\\vspace{0.40 cm}\n")
		}
	}
	builder.WriteString("\n\\end{document}")
	return builder.String()
}

func generateAwardsBwCV(data CVData) string {
	var builder strings.Builder
	builder.WriteString("\\documentclass[../main.tex]{subfiles}\n\\begin{document}\n\\section{\\sectiontitle[\\Large]{Awards}}\n")
	for i, award := range data.Awards {
		builder.WriteString(fmt.Sprintf(`
\begin{twocolentry}{
    %s
}
    \textbf{%s}\\
    %s
\end{twocolentry}
\vspace{0.10 cm}
\begin{onecolentry}
%s
\end{onecolentry}
`, escapeLatex(award.Date), escapeLatex(award.Title), escapeLatex(award.Organization), generateHighlights(award.Points)))
		if i < len(data.Awards)-1 {
			builder.WriteString("\n\\vspace{0.10 cm}\n")
		}
	}
	builder.WriteString("\n\\end{document}")
	return builder.String()
}

// == Projects ==
func generateProjectsMainCV(data CVData) string {
	var builder strings.Builder
	builder.WriteString("\\documentclass[../main.tex]{subfiles}\n\\begin{document}\n\\section{Projects}\n")
	for i, proj := range data.Projects {
		githubLink := fmt.Sprintf(`\href{%s}{%s}`, escapeLatex(proj.Github), escapeLatex(proj.GithubHandle))
		builder.WriteString(fmt.Sprintf(`
\begin{twocolentry}{
    \small
    \textit{%s}
}
    \textbf{%s}
\end{twocolentry}

\vspace{0.10 cm}
\begin{onecolentry}
%s
\end{onecolentry}
`, githubLink, escapeLatex(proj.Name), generateHighlights(proj.Points)))
		if i < len(data.Projects)-1 {
			builder.WriteString("\n\\vspace{0.2 cm}\n")
		}
	}
	builder.WriteString("\n\\end{document}")
	return builder.String()
}

func generateProjectsBwCV(data CVData) string {
	var builder strings.Builder
	builder.WriteString("\\documentclass[../main.tex]{subfiles}\n\\begin{document}\n\\section{\\sectiontitle[\\Large]{Projects}}\n")
	for i, proj := range data.Projects {
		githubLink := fmt.Sprintf(`\href{%s}{%s}`, escapeLatex(proj.Github), escapeLatex(proj.GithubHandle))
		builder.WriteString(fmt.Sprintf(`
\begin{twocolentry}{
    \small
    \textit{%s}
}
    \sectiontitle{%s}
\end{twocolentry}

\vspace{0.10 cm}
\begin{onecolentry}
%s
\end{onecolentry}
`, githubLink, escapeLatex(proj.Name), generateHighlights(proj.Points)))
		if i < len(data.Projects)-1 {
			builder.WriteString("\n\\vspace{0.2 cm}\n")
		}
	}
	builder.WriteString("\n\\end{document}")
	return builder.String()
}

// == Achievements & Skills ==
func generateSkillsAchievements(data CVData) string {
	skills := data.SkillsAchievements
	var builder strings.Builder
	builder.WriteString("\\begin{highlights}\n")
	builder.WriteString(fmt.Sprintf("    \\item \\textbf{Hard Skill:} %s\n", escapeLatex(strings.Join(skills.HardSkills, ", "))))
	builder.WriteString(fmt.Sprintf("    \\item \\textbf{Soft Skill:} %s\n", escapeLatex(strings.Join(skills.SoftSkills, ", "))))
	builder.WriteString(fmt.Sprintf("    \\item \\textbf{Programming Languages:} %s\n", escapeLatex(strings.Join(skills.ProgrammingLanguages, ", "))))
	builder.WriteString(fmt.Sprintf("    \\item \\textbf{Database Languages:} %s\n", escapeLatex(strings.Join(skills.DatabaseLanguages, ", "))))
	builder.WriteString(fmt.Sprintf("    \\item \\textbf{Misc:} %s\n", escapeLatex(strings.Join(skills.Misc, ", "))))
	for _, cert := range skills.Certificates {
		builder.WriteString(fmt.Sprintf("    \\item \\textbf{Certificate} (%d): %s\n", cert.Year, escapeLatex(cert.Name)))
	}
	builder.WriteString("\\end{highlights}\n")
	return builder.String()
}

func generateSkillsMainCV(data CVData) string {
	sectionContent := generateSkillsAchievements(data)
	return fmt.Sprintf("\\documentclass[../main.tex]{subfiles}\n\\begin{document}\n\\section{Achievement and Skills}\n%s\\end{document}", sectionContent)
}

func generateSkillsBwCV(data CVData) string {
	sectionContent := generateSkillsAchievements(data)
	return fmt.Sprintf("\\documentclass[../main.tex]{subfiles}\n\\begin{document}\n\\section{\\sectiontitle[\\Large]{{Achievement and Skills}}}\n%s\\end{document}", sectionContent)
}

// --- Main Script Logic ---
func main() {
	fmt.Println("Starting LaTeX section generation...")

	// Load data from JSON file
	byteValue, err := os.ReadFile(dataFile)
	if err != nil {
		log.Fatalf("Error: Data file not found at %s: %v", dataFile, err)
	}

	var cvData CVData
	if err := json.Unmarshal(byteValue, &cvData); err != nil {
		log.Fatalf("Error decoding JSON from %s: %v", dataFile, err)
	}

	// Generate and write files for main_cv
	fmt.Println("Generating sections for main_cv...")
	writeTexFile(filepath.Join(mainCVSectionsDir, "Header.tex"), generateHeaderMainCV(cvData))
	writeTexFile(filepath.Join(mainCVSectionsDir, "Experience.tex"), generateExperienceMainCV(cvData))
	writeTexFile(filepath.Join(mainCVSectionsDir, "Education.tex"), generateEducationMainCV(cvData))
	writeTexFile(filepath.Join(mainCVSectionsDir, "Awards.tex"), generateAwardsMainCV(cvData))
	writeTexFile(filepath.Join(mainCVSectionsDir, "Projects.tex"), generateProjectsMainCV(cvData))
	writeTexFile(filepath.Join(mainCVSectionsDir, "Achivements_Skills.tex"), generateSkillsMainCV(cvData))

	// Generate and write files for bw_cv
	fmt.Println("Generating sections for bw_cv...")
	writeTexFile(filepath.Join(bwCVSectionsDir, "Header.tex"), generateHeaderBwCV(cvData))
	writeTexFile(filepath.Join(bwCVSectionsDir, "Experience.tex"), generateExperienceBwCV(cvData))
	writeTexFile(filepath.Join(bwCVSectionsDir, "Education.tex"), generateEducationBwCV(cvData))
	writeTexFile(filepath.Join(bwCVSectionsDir, "Awards.tex"), generateAwardsBwCV(cvData))
	writeTexFile(filepath.Join(bwCVSectionsDir, "Projects.tex"), generateProjectsBwCV(cvData))
	writeTexFile(filepath.Join(bwCVSectionsDir, "Achivements_Skills.tex"), generateSkillsBwCV(cvData))

	fmt.Println("LaTeX section generation complete.")
}
