use crate::types::{Award, Education, Experience, PersonalInfo, Project, Skills};
use lazy_static::lazy_static;
use regex::Regex;

// --- Helper Functions ---

lazy_static! {
    static ref BOLD_REGEX: Regex = Regex::new(r"\*\*(.*?)\*\*").unwrap();
}

/// Escapes special LaTeX characters in a string.
/// The order of replacement is crucial, especially for the backslash.
fn escape_latex(text: &str) -> String {
    text
        // FIX: Backslash must be escaped FIRST.
        .replace('\\', r"\textbackslash{}")
        .replace('&', r"\&")
        .replace('%', r"\%")
        .replace('$', r"\$")
        .replace('#', r"\#")
        .replace('_', r"\_")
        .replace('{', r"\{")
        .replace('}', r"\}")
        .replace('~', r"\textasciitilde{}")
        .replace('^', r"\textasciicircum{}")
}

/// Finds markdown-style bold text (**text**) and converts it to LaTeX \textbf{text}.
fn format_latex_bold(text: &str) -> String {
    let escaped_text = escape_latex(text);
    BOLD_REGEX
        .replace_all(&escaped_text, r"\textbf{$1}")
        .to_string()
}

/// Generates a \begin{highlights} ... \end{highlights} block.
fn generate_highlights(points: &[String]) -> String {
    if points.is_empty() {
        return String::new();
    }
    let mut builder = String::from("\\begin{highlights}\n");
    for point in points {
        builder.push_str(&format!(" \\item {}\n", format_latex_bold(point)));
    }
    builder.push_str("\\end{highlights}");
    builder
}

// --- Section Generators ---

// == Header ==
pub fn generate_header_main_cv(data: &PersonalInfo) -> String {
    let info = data;
    let website_handle = info.website.replace("https://", "");
    format!(
        r#"\documentclass[../main.tex]{{subfiles}}
\begin{{document}}
\begin{{header}}
    \setcolumnwidth{{3.5cm, \fill}}
    \columnsep=0pt
    \raggedright
    \noindent %% Prevent paragraph indentation
    \begin{{paracol}}{{2}}
        \begin{{leftcolumn}}
            \vspace*{{\fill}} %% Vertical alignment
            \includegraphics[width=3cm]{{{}}}
            \vspace*{{\fill}} %% Vertical alignment
        \end{{leftcolumn}}
        \begin{{rightcolumn}}
            \textbf{{\fontsize{{18 pt}}{{24 pt}}\selectfont \color{{NavyBlue}} {}}} \\
            \normalsize
            \kern 0.25 cm%%
            \mbox{{\hrefWithoutArrow{{mailto:{}}}{{\color{{black}}{{\footnotesize\faEnvelope[regular]}}\hspace*{{0.13cm}}{}}}}}%%
            \kern 0.25 cm%%
            \AND%%
            \kern 0.25 cm%%
            \mbox{{\hrefWithoutArrow{{tel:{}}}{{\color{{black}}{{\footnotesize\faPhone*}}\hspace*{{0.13cm}}{}}}}}%%
            \kern 0.25 cm%%
            \AND%%
            \kern 0.25 cm%%
            \mbox{{\hrefWithoutArrow{{{}}}{{\color{{black}}{{\footnotesize\faLink}}\hspace*{{0.13cm}}{}}}}}%%
            \kern 0.25 cm%%
            \AND%%
            \kern 0.25 cm%%
            \mbox{{\hrefWithoutArrow{{{}}}{{\color{{black}}{{\footnotesize\faLinkedinIn}}\hspace*{{0.13cm}}{}}}}}%%
            \kern 0.25 cm%%
            \AND%%
            \kern 0.25 cm%%
            \mbox{{\hrefWithoutArrow{{{}}}{{\color{{black}}{{\footnotesize\faGithub}}\hspace*{{0.13cm}}{}}}}}%%
        \end{{rightcolumn}}
    \end{{paracol}}
\end{{header}}
\end{{document}}
"#,
        escape_latex(&info.profile_pic),
        escape_latex(&info.name),
        escape_latex(&info.email),
        escape_latex(&info.email),
        escape_latex(&info.phone),
        escape_latex(&info.phone),
        escape_latex(&info.website),
        escape_latex(&website_handle),
        escape_latex(&info.linkedin),
        escape_latex(&info.linkedin_handle),
        escape_latex(&info.github),
        escape_latex(&info.github_handle)
    )
}

pub fn generate_header_bw_cv(data: &PersonalInfo) -> String {
    let info = data;
    let phone_handle = info.phone.replacen("+62-", "0", 1);
    let website_handle = info.website.replace("https://", "");
    format!(
        r#"\documentclass[../main.tex]{{subfiles}}

\begin{{document}}
\begin{{header}}
\fontsize{{25 pt}}{{25 pt}}\selectfont {}

    \vspace{{5 pt}}

    \normalsize
    \mbox{{\color{{black}}{{\footnotesize\faHome}}\hspace*{{0.13cm}}{}}}%%
    \kern 5.0 pt%%
    \AND%%
    \kern 5.0 pt%%
    \mbox{{\hrefWithoutArrow{{mailto:{}}}{{\color{{black}}{{\footnotesize\faEnvelope[regular]}}\hspace*{{0.13cm}}{}}}}}%%
    \kern 5.0 pt%%
    \AND%%
    \kern 5.0 pt%%
    \mbox{{\hrefWithoutArrow{{tel:{}}}{{\color{{black}}{{\footnotesize\faPhone*}}\hspace*{{0.13cm}}{}}}}}
    \kern 5.0 pt%%
    \AND%%
    \kern 5.0 pt%%
    \mbox{{\hrefWithoutArrow{{{}}}{{\color{{black}}{{\footnotesize\faLink}}\hspace*{{0.13cm}}{}}}}}%%
    \kern 5.0 pt%%
    \AND%%
    \kern 5.0 pt%%
    \mbox{{\hrefWithoutArrow{{{}}}{{\color{{black}}{{\footnotesize\faLinkedinIn}}\hspace*{{0.13cm}}{}}}}}%%
    \kern 5.0 pt%%
    \AND%%
    \kern 5.0 pt%%
    \mbox{{\hrefWithoutArrow{{{}}}{{\color{{black}}{{\footnotesize\faGithub}}\hspace*{{0.13cm}}{}}}}}%%

\end{{header}}
\end{{document}}
"#,
        escape_latex(&info.name),
        escape_latex(&info.location),
        escape_latex(&info.email),
        escape_latex(&info.email),
        escape_latex(&info.phone),
        escape_latex(&phone_handle),
        escape_latex(&info.website),
        escape_latex(&website_handle),
        escape_latex(&info.linkedin),
        escape_latex(&info.linkedin_handle),
        escape_latex(&info.github),
        escape_latex(&info.github_handle)
    )
}

// == Experience ==
pub fn generate_experience_main_cv(data: &Vec<&Experience>) -> String {
    let mut builder = String::from(
        "\\documentclass[../main.tex]{subfiles}\n\\begin{document}\n\\section{Experience}\n",
    );
    for (i, exp) in data.iter().enumerate() {
        let role_type = format!(
            "{} ({})",
            escape_latex(&exp.role),
            escape_latex(&exp.job_type)
        );
        builder.push_str(&format!(
            r#"
\begin{{twocolentry}}{{{}}}
\textbf{{{} \color{{Gray}}- {}}}\\
\textit{{{}}}
\end{{twocolentry}}

\vspace{{0.10 cm}}
\begin{{onecolentry}}
{}
\end{{onecolentry}}
"#,
            escape_latex(&exp.dates),
            escape_latex(&exp.company),
            escape_latex(&exp.location),
            role_type,
            generate_highlights(&exp.points)
        ));
        if i < data.len() - 1 {
            builder.push_str("\n\\vspace{0.40 cm}\n");
        }
    }
    builder.push_str("\n\\end{document}");
    builder
}

pub fn generate_experience_bw_cv(data: &Vec<&Experience>) -> String {
    let mut builder = String::from(
        "\\documentclass[../main.tex]{subfiles}\n\\begin{document}\n\\section{\\sectiontitle[\\Large]{Experience}}\n",
    );
    for (i, exp) in data.iter().enumerate() {
        builder.push_str(&format!(
            r#"
\begin{{twocolentry}}{{{}}}
\sectiontitle{{{}}} \location{{-- {}}}\\
\vspace{{0.10 cm}}
\textit{{{}}}
\end{{twocolentry}}

\vspace{{0.10 cm}}
\begin{{onecolentry}}
{}
\end{{onecolentry}}
"#,
            escape_latex(&exp.dates),
            escape_latex(&exp.company),
            escape_latex(&exp.location),
            escape_latex(&exp.role),
            generate_highlights(&exp.points)
        ));
        if i < data.len() - 1 {
            builder.push_str("\n\\vspace{0.4 cm}\n");
        }
    }
    builder.push_str("\n\\end{document}");
    builder
}

// == Education ==
pub fn generate_education_main_cv(data: &Vec<&Education>) -> String {
    let mut builder = String::from(
        "\\documentclass[../main.tex]{subfiles}\n\\begin{document}\n\\section{Education}\n",
    );
    for edu in data {
        builder.push_str(&format!(
            r#"
\begin{{twocolentry}}{{
\small
\textit{{{}}}
}}
\textbf{{{}}} \\
\textit{{{}, {}}}
\end{{twocolentry}}
"#,
            escape_latex(&edu.dates),
            escape_latex(&edu.institution),
            escape_latex(&edu.degree),
            escape_latex(&edu.gpa)
        ));
    }
    builder.push_str("\n\\end{document}");
    builder
}

pub fn generate_education_bw_cv(data: &Vec<&Education>) -> String {
    let mut builder = String::from(
        "\\documentclass[../main.tex]{subfiles}\n\\begin{document}\n\\section{\\sectiontitle[\\Large]{Education}}\n",
    );
    for edu in data {
        let bold_edu_text = format!(
            r#"\textbf{{{}, {}}}"#,
            escape_latex(&edu.institution),
            escape_latex(&edu.degree)
        );
        builder.push_str(&format!(
            r#"
\begin{{twocolentry}}{{
{}
}}
{}
\end{{twocolentry}}
"#,
            escape_latex(&edu.dates),
            bold_edu_text
        ));
        if !edu.details.is_empty() {
            builder.push_str(&format!(
                "\n\\vspace{{0.10 cm}}\n\\begin{{onecolentry}}\n{}\n\\end{{onecolentry}}\n",
                generate_highlights(&edu.details)
            ));
        }
    }
    builder.push_str("\n\\end{document}");
    builder
}

// == Awards ==
pub fn generate_awards_main_cv(data: &Vec<&Award>) -> String {
    let mut builder = String::from(
        "\\documentclass[../main.tex]{subfiles}\n\\begin{document}\n\\section{Awards}\n",
    );
    for (i, award) in data.iter().enumerate() {
        builder.push_str(&format!(
            r#"
\begin{{twocolentry}}{{
\textit{{{}}}
}}
\textbf{{{}}}\\
{}
\end{{twocolentry}}
\vspace{{0.10 cm}}
\begin{{onecolentry}}
{}
\end{{onecolentry}}
"#,
            escape_latex(&award.date),
            escape_latex(&award.title),
            escape_latex(&award.organization),
            generate_highlights(&award.points)
        ));
        if i < data.len() - 1 {
            builder.push_str("\n\\vspace{0.40 cm}\n");
        }
    }
    builder.push_str("\n\\end{document}");
    builder
}

pub fn generate_awards_bw_cv(data: &Vec<&Award>) -> String {
    let mut builder = String::from(
        "\\documentclass[../main.tex]{subfiles}\n\\begin{document}\n\\section{\\sectiontitle[\\Large]{Awards}}\n",
    );
    for (i, award) in data.iter().enumerate() {
        builder.push_str(&format!(
            r#"
\begin{{twocolentry}}{{
{}
}}
\textbf{{{}}}\\
{}
\end{{twocolentry}}
\vspace{{0.10 cm}}
\begin{{onecolentry}}
{}
\end{{onecolentry}}
"#,
            escape_latex(&award.date),
            escape_latex(&award.title),
            escape_latex(&award.organization),
            generate_highlights(&award.points)
        ));
        if i < data.len() - 1 {
            builder.push_str("\n\\vspace{0.10 cm}\n");
        }
    }
    builder.push_str("\n\\end{document}");
    builder
}

// == Projects ==
pub fn generate_projects_main_cv(data: &Vec<&Project>) -> String {
    let mut builder = String::from(
        "\\documentclass[../main.tex]{subfiles}\n\\begin{document}\n\\section{Projects}\n",
    );
    for (i, proj) in data.iter().enumerate() {
        let github_link = format!(
            r#"\href{{{}}}{{{}}}"#,
            escape_latex(&proj.github),
            escape_latex(&proj.github_handle)
        );
        builder.push_str(&format!(
            r#"
\begin{{twocolentry}}{{
\small
\textit{{{}}}
}}
\textbf{{{}}}
\end{{twocolentry}}

\vspace{{0.10 cm}}
\begin{{onecolentry}}
{}
\end{{onecolentry}}
"#,
            github_link,
            escape_latex(&proj.name),
            generate_highlights(&proj.points)
        ));
        if i < data.len() - 1 {
            builder.push_str("\n\\vspace{0.2 cm}\n");
        }
    }
    builder.push_str("\n\\end{document}");
    builder
}

pub fn generate_projects_bw_cv(data: &Vec<&Project>) -> String {
    let mut builder = String::from(
        "\\documentclass[../main.tex]{subfiles}\n\\begin{document}\n\\section{\\sectiontitle[\\Large]{Projects}}\n",
    );
    for (i, proj) in data.iter().enumerate() {
        let github_link = format!(
            r#"\href{{{}}}{{{}}}"#,
            escape_latex(&proj.github),
            escape_latex(&proj.github_handle)
        );
        builder.push_str(&format!(
            r#"
\begin{{twocolentry}}{{
\small
\textit{{{}}}
}}
\sectiontitle{{{}}}
\end{{twocolentry}}

\vspace{{0.10 cm}}
\begin{{onecolentry}}
{}
\end{{onecolentry}}
"#,
            github_link,
            escape_latex(&proj.name),
            generate_highlights(&proj.points)
        ));
        if i < data.len() - 1 {
            builder.push_str("\n\\vspace{0.2 cm}\n");
        }
    }
    builder.push_str("\n\\end{document}");
    builder
}

// == Skills & Achievements ==
fn generate_skills_achievements(skills: &Skills) -> String {
    let mut builder = String::from("\\begin{highlights}\n");
    if !skills.hard_skills.is_empty() {
        builder.push_str(&format!(
            " \\item \\textbf{{Hard Skill:}} {}\n",
            escape_latex(&skills.hard_skills.join(", "))
        ));
    }
    if !skills.soft_skills.is_empty() {
        builder.push_str(&format!(
            " \\item \\textbf{{Soft Skill:}} {}\n",
            escape_latex(&skills.soft_skills.join(", "))
        ));
    }
    if !skills.programming_languages.is_empty() {
        builder.push_str(&format!(
            " \\item \\textbf{{Programming Languages:}} {}\n",
            escape_latex(&skills.programming_languages.join(", "))
        ));
    }
    if !skills.databases.is_empty() {
        builder.push_str(&format!(
            " \\item \\textbf{{Database:}} {}\n",
            escape_latex(&skills.databases.join(", "))
        ));
    }
    if !skills.misc.is_empty() {
        builder.push_str(&format!(
            " \\item \\textbf{{Frameworks \\& Misc:}} {}\n",
            escape_latex(&skills.misc.join(", "))
        ));
    }
    for cert in &skills.certificates {
        builder.push_str(&format!(
            " \\item \\textbf{{Certificate}} ({}): {}\n",
            cert.year,
            escape_latex(&cert.name)
        ));
    }
    builder.push_str("\\end{highlights}\n");
    builder
}

pub fn generate_skills_main_cv(data: &Skills) -> String {
    let section_content = generate_skills_achievements(data);
    format!(
        "\\documentclass[../main.tex]{{subfiles}}\n\\begin{{document}}\n\\section{{Skills}}\n{}\\end{{document}}",
        section_content
    )
}

pub fn generate_skills_bw_cv(data: &Skills) -> String {
    let section_content = generate_skills_achievements(data);
    format!(
        "\\documentclass[../main.tex]{{subfiles}}\n\\begin{{document}}\n\\section{{\\sectiontitle[\\Large]{{Skills}}}}\n{}\\end{{document}}",
        section_content
    )
}
