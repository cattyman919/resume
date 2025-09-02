use crate::{
    cv_model::{
        Award, Education, Experience, ExperiencesCVData, GeneralCVData, HasCvTypes, Project,
        ProjectsCVData,
    },
    latex,
};
use anyhow::{Context, Result};
use futures::future::join_all;
use log::info;
use std::{
    collections::HashSet,
    path::{Path, PathBuf},
    process::Stdio,
    sync::Arc,
};
use tokio::{fs, io, process::Command};

pub async fn get_all_cv_types(
    projects_cv: &ProjectsCVData,
    experiences_cv: &ExperiencesCVData,
) -> Result<HashSet<String>> {
    let experience_task = get_types_per_section(projects_cv);
    let project_task = get_types_per_section(experiences_cv);

    let (mut project_cv_types, experiences_cv_types) = tokio::join!(experience_task, project_task);
    project_cv_types.extend(experiences_cv_types);

    Ok(project_cv_types)
}

async fn get_types_per_section<T>(section: &[T]) -> HashSet<String>
where
    T: HasCvTypes,
{
    let mut local_types = HashSet::new();
    for item in section {
        for cv_type in item.cv_types() {
            local_types.insert(cv_type.into());
        }
    }
    local_types
}

pub async fn setup_directories() -> Result<()> {
    info!("Creating out & cv directory...");
    let create_out = fs::create_dir("out");
    let create_cv = fs::create_dir("cv");

    if let Err(e) = tokio::try_join!(create_out, create_cv)
        && e.kind() != std::io::ErrorKind::AlreadyExists
    {
        return Err(e).context("Failed to create initial directories");
    }

    info!("Directories are ready.");
    Ok(())
}

#[derive(Clone, Copy)]
enum Section {
    Header,
    Experience,
    Education,
    Awards,
    Projects,
    Skills,
}

impl Section {
    fn filename(&self) -> &'static str {
        match self {
            Section::Header => "Header.tex",
            Section::Experience => "Experience.tex",
            Section::Education => "Education.tex",
            Section::Awards => "Awards.tex",
            Section::Projects => "Projects.tex",
            Section::Skills => "Achivements_Skills.tex",
        }
    }
    fn generate_content(&self, style: &str, context: &CvContext) -> String {
        match self {
            Section::Header => {
                let data = &context.general.personal_info;
                match style {
                    "main" => latex::generate_header_main_cv(data),
                    "bw" => latex::generate_header_bw_cv(data),
                    _ => String::new(),
                }
            }
            Section::Experience => {
                let filtered: Vec<&Experience> = context
                    .experiences
                    .iter()
                    .filter(|e| e.cv_types().contains(context.cv_type.as_ref()))
                    .collect();
                match style {
                    "main" => latex::generate_experience_main_cv(&filtered),
                    "bw" => latex::generate_experience_bw_cv(&filtered),
                    _ => String::new(),
                }
            }
            Section::Education => {
                let data: Vec<&Education> = context.general.education.iter().collect();
                match style {
                    "main" => latex::generate_education_main_cv(&data),
                    "bw" => latex::generate_education_bw_cv(&data),
                    _ => String::new(),
                }
            }
            Section::Awards => {
                let data: Vec<&Award> = context.general.awards.iter().collect();
                match style {
                    "main" => latex::generate_awards_main_cv(&data),
                    "bw" => latex::generate_awards_bw_cv(&data),
                    _ => String::new(),
                }
            }
            Section::Projects => {
                let filtered: Vec<&Project> = context
                    .projects
                    .iter()
                    .filter(|p| p.cv_type.contains(context.cv_type.as_ref()))
                    .collect();
                match style {
                    "main" => latex::generate_projects_main_cv(&filtered),
                    "bw" => latex::generate_projects_bw_cv(&filtered),
                    _ => String::new(),
                }
            }
            Section::Skills => {
                let data = &context.general.skills_achievements;
                match style {
                    "main" => latex::generate_skills_main_cv(data),
                    "bw" => latex::generate_skills_bw_cv(data),
                    _ => String::new(),
                }
            }
        }
    }
}

struct CvContext {
    general: Arc<GeneralCVData>,
    projects: Arc<ProjectsCVData>,
    experiences: Arc<ExperiencesCVData>,
    cv_type: Arc<String>,
}

// Generates the .tex files for each section and compiles the final PDFs.
pub async fn write_cv(
    general_cv: Arc<GeneralCVData>,
    projects_cv: ProjectsCVData,
    experiences_cv: ExperiencesCVData,
    cv_type: String,
    debug_mode: bool,
) -> io::Result<()> {
    let main_path: PathBuf = format!("cv/{}/main_cv", cv_type).into();
    let bw_path: PathBuf = format!("cv/{}/bw_cv", cv_type).into();

    tokio::try_join!(
        copy_dir_recursively(Path::new("template_cv/main_cv"), &main_path),
        copy_dir_recursively(Path::new("template_cv/bw_cv"), &bw_path)
    )?;

    let main_sections_path = main_path.join("sections");
    let bw_sections_path = bw_path.join("sections");

    tokio::try_join!(
        fs::create_dir_all(&main_sections_path),
        fs::create_dir_all(&bw_sections_path)
    )?;

    let context = Arc::new(CvContext {
        general: general_cv,
        projects: Arc::new(projects_cv),
        experiences: Arc::new(experiences_cv),
        cv_type: Arc::new(cv_type.clone()),
    });

    let sections_to_generate = [
        Section::Header,
        Section::Experience,
        Section::Education,
        Section::Awards,
        Section::Projects,
        Section::Skills,
    ];

    let mut join_handles = Vec::new();
    for section in sections_to_generate {
        let main_path = Path::new(&main_sections_path).join(section.filename());
        let bw_path = Path::new(&bw_sections_path).join(section.filename());

        let context_clone = Arc::clone(&context);

        join_handles.push(tokio::spawn(async move {
            let content = section.generate_content("main", &context_clone);
            write_tex_file(&main_path, content).await
        }));

        let context_clone = Arc::clone(&context);

        join_handles.push(tokio::spawn(async move {
            let content = section.generate_content("bw", &context_clone);
            write_tex_file(&bw_path, content).await
        }));
    }

    for result in join_all(join_handles).await {
        result.unwrap()?;
    }

    println!(
        "Finished writing all .tex sections for CV type: {}",
        cv_type
    );

    let pdf_main_handle = write_pdf(
        &context.general.personal_info.name,
        &cv_type,
        "main",
        debug_mode,
    );
    let pdf_bw_handle = write_pdf(
        &context.general.personal_info.name,
        &cv_type,
        "bw",
        debug_mode,
    );

    tokio::try_join!(pdf_main_handle, pdf_bw_handle)?;

    Ok(())
}

/// Runs pdflatex to generate the final PDF.
async fn write_pdf(name: &str, cv_type: &str, style: &str, debug_mode: bool) -> io::Result<()> {
    let target_pdf = format!("{} - CV ({}) ({})", name, cv_type, style.to_uppercase());
    let working_dir = match style {
        "main" => format!("cv/{cv_type}/main_cv"),
        "bw" => format!("cv/{cv_type}/bw_cv"),
        _ => {
            return Err(io::Error::new(
                io::ErrorKind::InvalidInput,
                format!("Invalid style for PDF generation: {}", style),
            ));
        }
    };

    println!("Running pdflatex for {}...", target_pdf);

    let mut cmd = Command::new("pdflatex");
    cmd.current_dir(working_dir)
        .arg("-output-directory=../../../out")
        .arg("-output-format=pdf")
        .arg("-jobname")
        .arg(&target_pdf)
        .arg("main.tex");

    if debug_mode {
        cmd.stdout(Stdio::inherit()).stderr(Stdio::inherit());
    } else {
        cmd.stdout(Stdio::null()).stderr(Stdio::piped());
    }

    let output = cmd.spawn()?.wait_with_output().await?;

    if !output.status.success() {
        let stderr_info = if debug_mode {
            "See terminal output for details.".to_string()
        } else {
            String::from_utf8_lossy(&output.stderr).to_string()
        };
        let error_message = format!(
            "pdflatex failed for {}.pdf with exit code: {}\nError: {}",
            target_pdf, output.status, stderr_info
        );
        return Err(io::Error::other(error_message));
    }

    println!("Generated {}.pdf", target_pdf);
    Ok(())
}

/// Moves .log and .aux files to an 'out/aux' directory.
pub async fn move_aux_files() -> io::Result<()> {
    let out_dir = Path::new("out");
    let aux_dir = out_dir.join("aux");

    fs::create_dir_all(&aux_dir).await?;

    let mut entries = fs::read_dir(out_dir).await?;
    while let Some(entry) = entries.next_entry().await? {
        let path = entry.path();

        if path.is_file()
            && let Some(ext) = path.extension().and_then(|s| s.to_str())
            && (ext == "log" || ext == "aux")
            && let Some(file_name) = path.file_name()
        {
            let new_path = aux_dir.join(file_name);
            fs::rename(&path, &new_path).await?;
        }
    }
    Ok(())
}

#[inline]
async fn write_tex_file(path: &Path, content: String) -> io::Result<()> {
    if content.is_empty() {
        return Ok(());
    }
    fs::write(path, content).await?;
    println!("Successfully wrote file: {:?}", path);
    Ok(())
}

async fn copy_dir_recursively(from: &Path, to: &Path) -> io::Result<()> {
    fs::create_dir_all(to).await?;

    let mut reader = fs::read_dir(from).await?;

    while let Some(entry) = reader.next_entry().await? {
        let file_type = entry.file_type().await?;

        let to_path_entry = to.join(entry.file_name());

        if file_type.is_dir() {
            Box::pin(copy_dir_recursively(&entry.path(), &to_path_entry)).await?;
        } else {
            fs::copy(entry.path(), &to_path_entry).await?;
        }
    }

    Ok(())
}
