package generator

import (
	"fmt"
	"regexp"
	"resume/internals/resume"
	"strings"
)

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

func GenerateHighlights(points []string) string {
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
func GenerateHeaderMainCV(data *resume.CVData) string {
	info := data.General.PersonalInfo
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

func GenerateHeaderBwCV(data *resume.CVData) string {

	info := data.General.PersonalInfo
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
func GenerateExperienceMainCV(data *resume.CVData) string {

	var builder strings.Builder
	builder.WriteString("\\documentclass[../main.tex]{subfiles}\n\\begin{document}\n\\section{Experience}\n")
	for i, exp := range data.Experiences {
		roleType := fmt.Sprintf("%s (%s)", escapeLatex(exp.Role), escapeLatex(exp.JobType))
		builder.WriteString(fmt.Sprintf(`
\begin{twocolentry}{%s}
    \textbf{%s \color{Gray}- %s}\\
    \textit{%s}
\end{twocolentry}

\vspace{0.10 cm}
\begin{onecolentry}
%s
\end{onecolentry}
`, escapeLatex(exp.Dates), escapeLatex(exp.Company), escapeLatex(exp.Location), roleType, GenerateHighlights(exp.Points)))
		if i < len(data.Experiences)-1 {
			builder.WriteString("\n\\vspace{0.40 cm}\n")
		}
	}
	builder.WriteString("\n\\end{document}")
	return builder.String()
}

func GenerateExperienceBwCV(data *resume.CVData) string {

	var builder strings.Builder
	builder.WriteString("\\documentclass[../main.tex]{subfiles}\n\\begin{document}\n\\section{\\sectiontitle[\\Large]{Experience}}\n")
	for i, exp := range data.Experiences {
		builder.WriteString(fmt.Sprintf(`
\begin{twocolentry}{%s}
    \sectiontitle{%s} \location{-- %s}\\
    \vspace{0.10 cm}
    \textbf{\textit{%s}} (\textit{%s})
\end{twocolentry}

\vspace{0.10 cm}
\begin{onecolentry}
%s
\end{onecolentry}
`, escapeLatex(exp.Dates), escapeLatex(exp.Role), escapeLatex(exp.JobType), escapeLatex(exp.Company), escapeLatex(exp.Location), GenerateHighlights(exp.Points)))
		if i < len(data.Experiences)-1 {
			builder.WriteString("\n\\vspace{0.4 cm}\n")
		}
	}
	builder.WriteString("\n\\end{document}")
	return builder.String()
}

// == Education ==
func GenerateEducationMainCV(data *resume.CVData) string {

	var builder strings.Builder
	builder.WriteString("\\documentclass[../main.tex]{subfiles}\n\\begin{document}\n\\section{Education}\n")
	for _, edu := range data.General.Education {
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

func GenerateEducationBwCV(data *resume.CVData) string {

	var builder strings.Builder
	builder.WriteString("\\documentclass[../main.tex]{subfiles}\n\\begin{document}\n\\section{\\sectiontitle[\\Large]{Education}}\n")
	for _, edu := range data.General.Education {
		boldEduText := fmt.Sprintf(`\textbf{%s, %s}`, escapeLatex(edu.Institution), escapeLatex(edu.Degree))
		builder.WriteString(fmt.Sprintf(`
\begin{twocolentry}{
    %s
}
    %s
\end{twocolentry}
`, escapeLatex(edu.Dates), boldEduText))
		if len(edu.Details) > 0 {
			builder.WriteString(fmt.Sprintf("\n\\vspace{0.10 cm}\n\\begin{onecolentry}\n%s\n\\end{onecolentry}\n", GenerateHighlights(edu.Details)))
		}
	}
	builder.WriteString("\n\\end{document}")
	return builder.String()
}

// == Awards ==
func GenerateAwardsMainCV(data *resume.CVData) string {

	var builder strings.Builder
	builder.WriteString("\\documentclass[../main.tex]{subfiles}\n\\begin{document}\n\\section{Awards}\n")
	for i, award := range data.General.Awards {
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
`, escapeLatex(award.Date), escapeLatex(award.Title), escapeLatex(award.Organization), GenerateHighlights(award.Points)))
		if i < len(data.General.Awards)-1 {
			builder.WriteString("\n\\vspace{0.40 cm}\n")
		}
	}
	builder.WriteString("\n\\end{document}")
	return builder.String()
}

func GenerateAwardsBwCV(data *resume.CVData) string {

	var builder strings.Builder
	builder.WriteString("\\documentclass[../main.tex]{subfiles}\n\\begin{document}\n\\section{\\sectiontitle[\\Large]{Awards}}\n")
	for i, award := range data.General.Awards {
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
`, escapeLatex(award.Date), escapeLatex(award.Title), escapeLatex(award.Organization), GenerateHighlights(award.Points)))
		if i < len(data.General.Awards)-1 {
			builder.WriteString("\n\\vspace{0.10 cm}\n")
		}
	}
	builder.WriteString("\n\\end{document}")
	return builder.String()
}

// == Projects ==
func GenerateProjectsMainCV(data *resume.CVData) string {

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
`, githubLink, escapeLatex(proj.Name), GenerateHighlights(proj.Points)))
		if i < len(data.Projects)-1 {
			builder.WriteString("\n\\vspace{0.2 cm}\n")
		}
	}
	builder.WriteString("\n\\end{document}")
	return builder.String()
}

func GenerateProjectsBwCV(data *resume.CVData) string {

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
`, githubLink, escapeLatex(proj.Name), GenerateHighlights(proj.Points)))
		if i < len(data.Projects)-1 {
			builder.WriteString("\n\\vspace{0.2 cm}\n")
		}
	}
	builder.WriteString("\n\\end{document}")
	return builder.String()
}

// == Achievements & Skills ==
func GenerateSkillsAchievements(data *resume.CVData) string {
	skills := data.General.SkillsAchievements
	var builder strings.Builder
	builder.WriteString("\\begin{highlights}\n")
	builder.WriteString(fmt.Sprintf("    \\item \\textbf{Hard Skill:} %s\n", escapeLatex(strings.Join(skills.HardSkills, ", "))))
	builder.WriteString(fmt.Sprintf("    \\item \\textbf{Soft Skill:} %s\n", escapeLatex(strings.Join(skills.SoftSkills, ", "))))
	builder.WriteString(fmt.Sprintf("    \\item \\textbf{Programming Languages:} %s\n", escapeLatex(strings.Join(skills.ProgrammingLanguages, ", "))))
	builder.WriteString(fmt.Sprintf("    \\item \\textbf{Database:} %s\n", escapeLatex(strings.Join(skills.DatabaseLanguages, ", "))))
	builder.WriteString(fmt.Sprintf("    \\item \\textbf{Frameworks \\& Misc:} %s\n", escapeLatex(strings.Join(skills.Misc, ", "))))
	for _, cert := range skills.Certificates {
		builder.WriteString(fmt.Sprintf("    \\item \\textbf{Certificate} (%d): %s\n", cert.Year, escapeLatex(cert.Name)))
	}
	builder.WriteString("\\end{highlights}\n")
	return builder.String()
}

func GenerateSkillsMainCV(data *resume.CVData) string {

	sectionContent := GenerateSkillsAchievements(data)
	return fmt.Sprintf("\\documentclass[../main.tex]{subfiles}\n\\begin{document}\n\\section{Skills}\n%s\\end{document}", sectionContent)
}

func GenerateSkillsBwCV(data *resume.CVData) string {

	sectionContent := GenerateSkillsAchievements(data)
	return fmt.Sprintf("\\documentclass[../main.tex]{subfiles}\n\\begin{document}\n\\section{\\sectiontitle[\\Large]{{Skills}}}\n%s\\end{document}", sectionContent)
}
