use crate::{
    latex,
    types::{Award, CVData, Education, Experience, HasCvTypes, Project},
};
use futures::future::join_all;
use std::{collections::HashSet, path::Path, process::Stdio, sync::Arc};
use tokio::{fs, io, process::Command, sync::Mutex};

pub async fn get_all_cv_types(cv_data: &CVData) -> Result<HashSet<String>, &'static str> {
    let total_types = Arc::new(Mutex::new(HashSet::<String>::new()));

    let experience_task = get_types_per_section(&cv_data.experiences, Arc::clone(&total_types));
    let project_task = get_types_per_section(&cv_data.projects, Arc::clone(&total_types));

    tokio::join!(experience_task, project_task);

    let final_set = Arc::try_unwrap(total_types)
        .expect("Error: Arc should be uniquely owned here")
        .into_inner();

    Ok(final_set)
}

async fn get_types_per_section<T>(section: &[T], total_types: Arc<Mutex<HashSet<String>>>)
where
    T: HasCvTypes,
{
    let mut local_types = HashSet::new();
    for item in section {
        for cv_type in item.cv_types() {
            local_types.insert(cv_type.into());
        }
    }
    total_types.lock().await.extend(local_types);
}

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
}
pub async fn write_cv(cv_data: Arc<CVData>, cv_type: String, debug_mode: bool) -> io::Result<()> {
    let main_path_dir = format!("cv/{cv_type}/main_cv");
    let bw_path_dir = format!("cv/{cv_type}/bw_cv");

    let main_path = Path::new(&main_path_dir);
    let bw_path = Path::new(&bw_path_dir);

    tokio::try_join!(
        copy_dir_recursively(Path::new("template_cv/main_cv"), main_path),
        copy_dir_recursively(Path::new("template_cv/bw_cv"), bw_path)
    )?;

    let main_sections_path = main_path.join("sections");
    let bw_sections_path = bw_path.join("sections");

    tokio::try_join!(
        fs::create_dir_all(&main_sections_path),
        fs::create_dir_all(&bw_sections_path)
    )?;

    let mut handles = Vec::new();
    let sections_to_generate = [
        Section::Header,
        Section::Experience,
        Section::Education,
        Section::Awards,
        Section::Projects,
        Section::Skills,
    ];

    let cv_type = Arc::new(cv_type);

    for section in &sections_to_generate {
        let main_path = Path::new(&main_sections_path).join(section.filename());
        let bw_path = Path::new(&bw_sections_path).join(section.filename());
        let data_main = Arc::clone(&cv_data);
        let data_bw = Arc::clone(&cv_data);
        let type_clone_main = Arc::clone(&cv_type);
        let type_clone_bw = Arc::clone(&cv_type);

        match section {
            Section::Header => {
                handles.push(tokio::spawn(async move {
                    let content = latex::generate_header_main_cv(&data_main.personal_info);
                    write_tex_file(&main_path, content).await
                }));
                handles.push(tokio::spawn(async move {
                    let content = latex::generate_header_bw_cv(&data_bw.personal_info);
                    write_tex_file(&bw_path, content).await
                }));
            }
            Section::Experience => {
                handles.push(tokio::spawn(async move {
                    let filtered: Vec<&Experience> = data_main
                        .experiences
                        .iter()
                        .filter(|e| e.cv_types.contains(&type_clone_main))
                        .collect();
                    let content = latex::generate_experience_main_cv(&filtered);
                    write_tex_file(&main_path, content).await
                }));
                handles.push(tokio::spawn(async move {
                    let filtered: Vec<&Experience> = data_bw
                        .experiences
                        .iter()
                        .filter(|e| e.cv_types.contains(&type_clone_bw))
                        .collect();
                    let content = latex::generate_experience_bw_cv(&filtered);
                    write_tex_file(&bw_path, content).await
                }));
            }
            Section::Education => {
                handles.push(tokio::spawn(async move {
                    let edu_refs: Vec<&Education> = data_main.education.iter().collect();
                    let content = latex::generate_education_main_cv(&edu_refs);
                    write_tex_file(&main_path, content).await
                }));
                handles.push(tokio::spawn(async move {
                    let edu_refs: Vec<&Education> = data_bw.education.iter().collect();
                    let content = latex::generate_education_bw_cv(&edu_refs);
                    write_tex_file(&bw_path, content).await
                }));
            }
            Section::Awards => {
                handles.push(tokio::spawn(async move {
                    let award_refs: Vec<&Award> = data_main.awards.iter().collect();
                    let content = latex::generate_awards_main_cv(&award_refs);
                    write_tex_file(&main_path, content).await
                }));
                handles.push(tokio::spawn(async move {
                    let award_refs: Vec<&Award> = data_bw.awards.iter().collect();
                    let content = latex::generate_awards_bw_cv(&award_refs);
                    write_tex_file(&bw_path, content).await
                }));
            }
            Section::Projects => {
                handles.push(tokio::spawn(async move {
                    let filtered: Vec<&Project> = data_main
                        .projects
                        .iter()
                        .filter(|p| p.cv_types.contains(&type_clone_main))
                        .collect();
                    let content = latex::generate_projects_main_cv(&filtered);
                    write_tex_file(&main_path, content).await
                }));
                handles.push(tokio::spawn(async move {
                    let filtered: Vec<&Project> = data_bw
                        .projects
                        .iter()
                        .filter(|p| p.cv_types.contains(&type_clone_bw))
                        .collect();
                    let content = latex::generate_projects_bw_cv(&filtered);
                    write_tex_file(&bw_path, content).await
                }));
            }
            Section::Skills => {
                handles.push(tokio::spawn(async move {
                    let content = latex::generate_skills_main_cv(&data_main.skills);
                    write_tex_file(&main_path, content).await
                }));
                handles.push(tokio::spawn(async move {
                    let content = latex::generate_skills_bw_cv(&data_bw.skills);
                    write_tex_file(&bw_path, content).await
                }));
            }
        }
    }

    for result in join_all(handles).await {
        result.unwrap()?;
    }

    println!(
        "Finished writing all .tex sections for CV type: {}",
        cv_type
    );

    let pdf_main_handle = write_pdf(&cv_type, "main", debug_mode);
    let pdf_bw_handle = write_pdf(&cv_type, "bw", debug_mode);

    tokio::try_join!(pdf_main_handle, pdf_bw_handle)?;

    Ok(())
}

/// Runs pdflatex to generate the final PDF.
async fn write_pdf(cv_type: &str, style: &str, debug_mode: bool) -> io::Result<()> {
    let target_pdf = format!(
        "Seno Pamungkas Rahman - CV ({}) ({})",
        cv_type,
        style.to_uppercase()
    );
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
        return Err(io::Error::new(io::ErrorKind::Other, error_message));
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
        if path.is_file() {
            if let Some(ext) = path.extension().and_then(|s| s.to_str()) {
                if ext == "log" || ext == "aux" {
                    if let Some(file_name) = path.file_name() {
                        let new_path = aux_dir.join(file_name);
                        fs::rename(&path, &new_path).await?;
                        // println!("Moved {:?} to {:?}", path, new_path);
                    }
                }
            }
        }
    }
    Ok(())
}

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
