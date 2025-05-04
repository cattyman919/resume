import json
import os
import re

# --- Configuration ---
DATA_FILE = 'cv_data.json'
MAIN_CV_SECTIONS_DIR = 'main_cv/sections'
BW_CV_SECTIONS_DIR = 'bw_cv/sections'

# --- Helper Functions ---

def escape_latex(text):
    """Escapes LaTeX special characters in a string."""
    if not isinstance(text, str):
        return text
    chars = {
        '&': r'\&',
        '%': r'\%',
        '$': r'\$',
        '#': r'\#',
        '_': r'\_',
        '{': r'\{',
        '}': r'\}',
        '~': r'\textasciitilde{}',
        '^': r'\textasciicircum{}',
        '\\': r'\textbackslash{}',
    }
    # Regex to apply escape except for already escaped sequences
    regex = re.compile('|'.join(re.escape(key) for key in chars.keys()))
    return regex.sub(lambda match: chars[match.group(0)], text)

def format_latex_bold(text):
    """Converts **bold** markdown to \textbf{bold} LaTeX."""
    # Escape text first to avoid escaping the command itself
    escaped_text = escape_latex(text)
    # Then apply bold formatting
    return re.sub(r'\*\*(.*?)\*\*', r'\\textbf{\1}', escaped_text)

def write_tex_file(filepath, content):
    """Writes content to a .tex file, creating directories if needed."""
    os.makedirs(os.path.dirname(filepath), exist_ok=True)
    try:
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(content)
        #print(f"Successfully wrote {filepath}")
    except IOError as e:
        print(f"Error writing {filepath}: {e}")

def generate_highlights(points):
    """Generates the LaTeX highlights list."""
    items = "".join([f"    \\item {format_latex_bold(point)}\n" for point in points])
    return f"\\begin{{highlights}}\n{items}\\end{{highlights}}"

# --- Section Generators ---

# == Header ==
def generate_header_main_cv(data):
    info = data['personal_info']
    content = f"""\\documentclass[../main.tex]{{subfiles}}
\\begin{{document}}
\\begin{{header}}
    \\setcolumnwidth{{3.5cm, \\fill}}
    \\columnsep=0pt
    \\raggedright
    \\noindent % Prevent paragraph indentation
    \\begin{{paracol}}{{2}}
        \\begin{{leftcolumn}}
            \\vspace*{{\\fill}} % Vertical alignment
            \\includegraphics[width=3cm]{{{escape_latex(info['profile_pic'])}}}
            \\vspace*{{\\fill}} % Vertical alignment
        \\end{{leftcolumn}}
        \\begin{{rightcolumn}}
            \\textbf{{\\fontsize{{18 pt}}{{24 pt}}\\selectfont \\color{{NavyBlue}} {escape_latex(info['name'])}}} \\\\
            \\normalsize
            \\kern 0.25 cm%
            \\mbox{{\\hrefWithoutArrow{{mailto:{escape_latex(info['email'])}}}{{\\color{{black}}{{\\footnotesize\\faEnvelope[regular]}}\\hspace*{{0.13cm}}{escape_latex(info['email'])}}}}}%
            \\kern 0.25 cm%
            \\AND%
            \\kern 0.25 cm%
            \\mbox{{\\hrefWithoutArrow{{tel:{escape_latex(info['phone'])}}}{{\\color{{black}}{{\\footnotesize\\faPhone*}}\\hspace*{{0.13cm}}{escape_latex(info['phone'])}}}}}%
            \\kern 0.25 cm%
            \\AND%
            \\kern 0.25 cm%
            \\mbox{{\\hrefWithoutArrow{{{escape_latex(info['website'])}}}{{\\color{{black}}{{\\footnotesize\\faLink}}\\hspace*{{0.13cm}}{escape_latex(info['website'].replace('https://',''))}}}}}%
            \\kern 0.25 cm%
            \\AND%
            \\kern 0.25 cm%
            \\mbox{{\\hrefWithoutArrow{{{escape_latex(info['linkedin'])}}}{{\\color{{black}}{{\\footnotesize\\faLinkedinIn}}\\hspace*{{0.13cm}}{escape_latex(info['linkedin_handle'])}}}}}%
            \\kern 0.25 cm%
            \\AND%
            \\kern 0.25 cm%
            \\mbox{{\\hrefWithoutArrow{{{escape_latex(info['github'])}}}{{\\color{{black}}{{\\footnotesize\\faGithub}}\\hspace*{{0.13cm}}{escape_latex(info['github_handle'])}}}}}%
        \\end{{rightcolumn}}
    \\end{{paracol}}
\\end{{header}}
\\end{{document}}
"""
    return content

def generate_header_bw_cv(data):
    info = data['personal_info']
    content = f"""\\documentclass[../main.tex]{{subfiles}}
\\begin{{document}}
\\begin{{header}}
    \\fontsize{{25 pt}}{{25 pt}}\\selectfont {escape_latex(info['name'])}

    \\vspace{{5 pt}}

    \\normalsize
    \\mbox{{\\color{{black}}{{\\footnotesize\\faHome}}\\hspace*{{0.13cm}}{escape_latex(info['location'])}}}%
    \\kern 5.0 pt%
    \\AND%
    \\kern 5.0 pt%
    \\mbox{{\\hrefWithoutArrow{{mailto:{escape_latex(info['email'])}}}{{\\color{{black}}{{\\footnotesize\\faEnvelope[regular]}}\\hspace*{{0.13cm}}{escape_latex(info['email'])}}}}}%
    \\kern 5.0 pt%
    \\AND%
    \\kern 5.0 pt%
    \\mbox{{\\hrefWithoutArrow{{tel:{escape_latex(info['phone'])}}}{{\\color{{black}}{{\\footnotesize\\faPhone*}}\\hspace*{{0.13cm}}{escape_latex(info['phone'].replace('+62-','0'))}}}}}
    \\kern 5.0 pt%
    \\AND%
    \\kern 5.0 pt%
    \\mbox{{\\hrefWithoutArrow{{{escape_latex(info['website'])}}}{{\\color{{black}}{{\\footnotesize\\faLink}}\\hspace*{{0.13cm}}{escape_latex(info['website'].replace('https://',''))}}}}}%
    \\kern 5.0 pt%
    \\AND%
    \\kern 5.0 pt%
    \\mbox{{\\hrefWithoutArrow{{{escape_latex(info['linkedin'])}}}{{\\color{{black}}{{\\footnotesize\\faLinkedinIn}}\\hspace*{{0.13cm}}{escape_latex(info['linkedin_handle'])}}}}}%
    \\kern 5.0 pt%
    \\AND%
    \\kern 5.0 pt%
    \\mbox{{\\hrefWithoutArrow{{{escape_latex(info['github'])}}}{{\\color{{black}}{{\\footnotesize\\faGithub}}\\hspace*{{0.13cm}}{escape_latex(info['github_handle'])}}}}}%
\\end{{header}}
\\end{{document}}
"""
    return content

# == Experience ==
def generate_experience_main_cv(data):
    content = "\\documentclass[../main.tex]{subfiles}\n\\begin{document}\n\\section{Experience}\n"
    for i, exp in enumerate(data['experiences']):
        role_type = f"{escape_latex(exp['role'])} ({escape_latex(exp['type'])})"
        content += f"""
\\begin{{twocolentry}}{{{escape_latex(exp['dates'])}}}
    \\textbf{{{escape_latex(exp['company'])} \\color{{Gray}}- {escape_latex(exp['location'])}}}\\\\
    \\textit{{{role_type}}}
\\end{{twocolentry}}

\\vspace{{0.10 cm}}
\\begin{{onecolentry}}
{generate_highlights(exp['points'])}
\\end{{onecolentry}}
"""
        if i < len(data['experiences']) - 1:
            content += "\n\\vspace{0.40 cm}\n" # Add space between entries
    content += "\n\\end{document}"
    return content

def generate_experience_bw_cv(data):
    content = "\\documentclass[../main.tex]{subfiles}\n\\begin{document}\n\\section{\\sectiontitle[\\Large]{Experience}}\n"
    for i, exp in enumerate(data['experiences']):
        content += f"""
\\begin{{twocolentry}}{{{escape_latex(exp['dates'])}}}
    \\sectiontitle{{{escape_latex(exp['company'])}}} \\location{{-- {escape_latex(exp['location'])}}}\\\\
    \\vspace{{0.10 cm}}
    \\textit{{{escape_latex(exp['role'])}}}
\\end{{twocolentry}}

\\vspace{{0.10 cm}}
\\begin{{onecolentry}}
{generate_highlights(exp['points'])}
\\end{{onecolentry}}
"""
        if i < len(data['experiences']) - 1:
            content += "\n\\vspace{0.4 cm}\n" # Add space between entries
    content += "\n\\end{document}"
    return content

# == Education ==
def generate_education_main_cv(data):
    content = "\\documentclass[../main.tex]{subfiles}\n\\begin{document}\n\\section{Education}\n"
    for edu in data['education']:
        content += f"""
\\begin{{twocolentry}}{{
    \\small
    \\textit{{{escape_latex(edu['dates'])}}}
}}
    \\textbf{{{escape_latex(edu['institution'])}}} \\\\
    \\textit{{{escape_latex(edu['degree'])}, {escape_latex(edu['gpa'])}}}
\\end{{twocolentry}}
"""
        # If details exist, add them (optional based on your main_cv style)
        # if 'details' in edu and edu['details']:
        #    content += f"\\vspace{{0.10 cm}}\n\\begin{{onecolentry}}\n{generate_highlights(edu['details'])}\n\\end{{onecolentry}}\n"

    content += "\n\\end{document}"
    return content

def generate_education_bw_cv(data):
    content = "\\documentclass[../main.tex]{subfiles}\n\\begin{document}\n\\section{\\sectiontitle[\\Large]{Education}}\n"
    for edu in data['education']:
        # --- Start Change ---
        # Pre-format the institution and degree string
        inst = escape_latex(edu['institution'])
        deg = escape_latex(edu['degree'])
        # Create the \textbf{...} part separately using an f-string
        bold_edu_text = f"\\textbf{{{inst}, {deg}}}"
        # --- End Change ---

        content += f"""
\\begin{{twocolentry}}{{
    {escape_latex(edu['dates'])}
}}
    {bold_edu_text}
\\end{{twocolentry}}
"""
        # The rest of the function remains the same
        if 'details' in edu and edu['details']:
            content += f"\n\\vspace{{0.10 cm}}\n\\begin{{onecolentry}}\n{generate_highlights(edu['details'])}\n\\end{{onecolentry}}\n"

    content += "\n\\end{document}"
    return content

# == Awards ==
def generate_awards_main_cv(data):
    content = "\\documentclass[../main.tex]{subfiles}\n\\begin{document}\n\\section{Awards}\n"
    for i, award in enumerate(data['awards']):
        content += f"""
\\begin{{twocolentry}}{{
    \\textit{{{escape_latex(award['date'])}}}
}}
    \\textbf{{{escape_latex(award['title'])}}}\\\\
    {escape_latex(award['organization'])}
\\end{{twocolentry}}
\\vspace{{0.10 cm}}
\\begin{{onecolentry}}
{generate_highlights(award['points'])}
\\end{{onecolentry}}
"""
        if i < len(data['awards']) - 1:
             content += "\n\\vspace{0.40 cm}\n" # Add space if multiple awards
    content += "\n\\end{document}"
    return content

def generate_awards_bw_cv(data):
    content = "\\documentclass[../main.tex]{subfiles}\n\\begin{document}\n\\section{\\sectiontitle[\\Large]{Awards}}\n"
    for i, award in enumerate(data['awards']):
        content += f"""
\\begin{{twocolentry}}{{
    {escape_latex(award['date'])}
}}
    \\textbf{{{escape_latex(award['title'])}}}\\\\
    {escape_latex(award['organization'])}
\\end{{twocolentry}}
\\vspace{{0.10 cm}}
\\begin{{onecolentry}}
{generate_highlights(award['points'])}
\\end{{onecolentry}}
"""
        if i < len(data['awards']) - 1:
             content += "\n\\vspace{0.10 cm}\n" # Add space if multiple awards
    content += "\n\\end{document}"
    return content

# == Projects ==
def generate_projects_main_cv(data):
    content = "\\documentclass[../main.tex]{subfiles}\n\\begin{document}\n\\section{Projects}\n"
    for i, proj in enumerate(data['projects']):
        github_link = f"\\href{{{escape_latex(proj['github'])}}}{{{escape_latex(proj['github_handle'])}}}"
        content += f"""
\\begin{{twocolentry}}{{
    \\small
    \\textit{{{github_link}}}
}}
    \\textbf{{{escape_latex(proj['name'])}}}
\\end{{twocolentry}}

\\vspace{{0.10 cm}}
\\begin{{onecolentry}}
{generate_highlights(proj['points'])}
\\end{{onecolentry}}
"""
        if i < len(data['projects']) - 1:
            content += "\n\\vspace{0.2 cm}\n" # Add space between entries
    content += "\n\\end{document}"
    return content

def generate_projects_bw_cv(data):
    content = "\\documentclass[../main.tex]{subfiles}\n\\begin{document}\n\\section{\\sectiontitle[\\Large]{Projects}}\n"
    for i, proj in enumerate(data['projects']):
        github_link = f"\\href{{{escape_latex(proj['github'])}}}{{{escape_latex(proj['github_handle'])}}}"
        content += f"""
\\begin{{twocolentry}}{{
    \\small
    \\textit{{{github_link}}}
}}
    \\sectiontitle{{{escape_latex(proj['name'])}}}
\\end{{twocolentry}}

\\vspace{{0.10 cm}}
\\begin{{onecolentry}}
{generate_highlights(proj['points'])}
\\end{{onecolentry}}
"""
        if i < len(data['projects']) - 1:
            content += "\n\\vspace{0.2 cm}\n" # Add space between entries
    content += "\n\\end{document}"
    return content

# == Achievements & Skills ==
def generate_skills_achievements(data):
    skills = data['skills_achievements']
    content = "\\begin{highlights}\n"
    content += f"    \\item \\textbf{{Hard Skill:}} {escape_latex(', '.join(skills['Hard Skills']))}\n"
    content += f"    \\item \\textbf{{Soft Skill:}} {escape_latex(', '.join(skills['Soft Skills']))}\n"
    content += f"    \\item \\textbf{{Programming Languages:}} {escape_latex(', '.join(skills['Programming Languages']))}\n"
    content += f"    \\item \\textbf{{Database Languages:}} {escape_latex(', '.join(skills['Database Languages']))}\n"
    content += f"    \\item \\textbf{{Misc:}} {escape_latex(', '.join(skills['Misc']))}\n"
    for cert in skills['Certificates']:
        content += f"    \\item \\textbf{{Certificate}} ({cert['year']}): {escape_latex(cert['name'])}\n"
    content += "\\end{highlights}\n"
    return content

def generate_skills_main_cv(data):
    section_content = generate_skills_achievements(data)
    return f"\\documentclass[../main.tex]{{subfiles}}\n\\begin{{document}}\n\\section{{Achievement and Skills}}\n{section_content}\\end{{document}}"

def generate_skills_bw_cv(data):
    section_content = generate_skills_achievements(data) # Content is the same, only section title differs
    return f"\\documentclass[../main.tex]{{subfiles}}\n\\begin{{document}}\n\\section{{\\sectiontitle[\\Large]{{Achievement and Skills}}}}\n{section_content}\\end{{document}}"


# --- Main Script Logic ---
if __name__ == "__main__":
    print("Starting LaTeX section generation...")

    # Check if data file exists
    if not os.path.exists(DATA_FILE):
        print(f"Error: Data file not found at {DATA_FILE}")
        exit(1)

    # Load data from JSON file
    try:
        with open(DATA_FILE, 'r', encoding='utf-8') as f:
            cv_data = json.load(f)
    except json.JSONDecodeError as e:
        print(f"Error decoding JSON from {DATA_FILE}: {e}")
        exit(1)
    except Exception as e:
        print(f"Error reading {DATA_FILE}: {e}")
        exit(1)

    # Generate and write files for main_cv
    print("Generating sections for main_cv...")
    write_tex_file(os.path.join(MAIN_CV_SECTIONS_DIR, 'Header.tex'), generate_header_main_cv(cv_data))
    write_tex_file(os.path.join(MAIN_CV_SECTIONS_DIR, 'Experience.tex'), generate_experience_main_cv(cv_data))
    write_tex_file(os.path.join(MAIN_CV_SECTIONS_DIR, 'Education.tex'), generate_education_main_cv(cv_data))
    write_tex_file(os.path.join(MAIN_CV_SECTIONS_DIR, 'Awards.tex'), generate_awards_main_cv(cv_data))
    write_tex_file(os.path.join(MAIN_CV_SECTIONS_DIR, 'Projects.tex'), generate_projects_main_cv(cv_data))
    write_tex_file(os.path.join(MAIN_CV_SECTIONS_DIR, 'Achivements_Skills.tex'), generate_skills_main_cv(cv_data)) # Note spelling mistake in original filename

    # Generate and write files for bw_cv
    print("Generating sections for bw_cv...")
    write_tex_file(os.path.join(BW_CV_SECTIONS_DIR, 'Header.tex'), generate_header_bw_cv(cv_data))
    write_tex_file(os.path.join(BW_CV_SECTIONS_DIR, 'Experience.tex'), generate_experience_bw_cv(cv_data))
    write_tex_file(os.path.join(BW_CV_SECTIONS_DIR, 'Education.tex'), generate_education_bw_cv(cv_data))
    write_tex_file(os.path.join(BW_CV_SECTIONS_DIR, 'Awards.tex'), generate_awards_bw_cv(cv_data))
    write_tex_file(os.path.join(BW_CV_SECTIONS_DIR, 'Projects.tex'), generate_projects_bw_cv(cv_data))
    write_tex_file(os.path.join(BW_CV_SECTIONS_DIR, 'Achivements_Skills.tex'), generate_skills_bw_cv(cv_data)) # Note spelling mistake in original filename

    print("LaTeX section generation complete.")
