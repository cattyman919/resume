use serde::de::DeserializeOwned;
use std::{error::Error, fmt, io, path::PathBuf, sync::Arc};
use tokio::fs;

use super::model::{ExperiencesCVConfig, GeneralCVConfig, ProjectsCVConfig};

const CONFIG_PATH: &str = "config";
const GENERAL_CV_PATH: [&str; 2] = [CONFIG_PATH, "general.yaml"];
const PROJECTS_CV_PATH: [&str; 2] = [CONFIG_PATH, "projects.yaml"];
const EXPERIENCES_CV_PATH: [&str; 2] = [CONFIG_PATH, "experiences.yaml"];

pub async fn load_cv_config()
-> Result<(GeneralCVConfig, ProjectsCVConfig, ExperiencesCVConfig), CvLoadError> {
    let (general_cv, projects_cv, experiences_cv) = tokio::try_join!(
        load_parse_yaml(GENERAL_CV_PATH),
        load_parse_yaml(PROJECTS_CV_PATH),
        load_parse_yaml(EXPERIENCES_CV_PATH)
    )?;

    Ok((general_cv, projects_cv, experiences_cv))
}

async fn load_parse_yaml<T: DeserializeOwned>(path_segments: [&str; 2]) -> Result<T, CvLoadError> {
    let path_buf: PathBuf = path_segments.iter().collect();
    let path_arc = Arc::new(path_buf);

    let content = fs::read_to_string(path_arc.as_path())
        .await
        .map_err(|e| CvLoadError::Io {
            path: path_arc.clone(),
            source: e,
        })?;

    let data = serde_yaml::from_str(&content).map_err(|e| CvLoadError::Parse {
        path: path_arc.clone(),
        source: e,
    })?;

    Ok(data)
}

#[derive(Debug)]
pub enum CvLoadError {
    Io {
        path: Arc<PathBuf>,
        source: io::Error,
    },
    Parse {
        path: Arc<PathBuf>,
        source: serde_yaml::Error,
    },
}

impl Error for CvLoadError {}

impl fmt::Display for CvLoadError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        // Helper to extract path for display
        let path_display = match self {
            CvLoadError::Io { path, .. } => path.display(),
            CvLoadError::Parse { path, .. } => path.display(),
        };

        match self {
            CvLoadError::Io { source, .. } => match source.kind() {
                io::ErrorKind::NotFound => {
                    write!(
                        f,
                        "Error: Configuration file '{}' was not found.",
                        path_display
                    )
                }
                io::ErrorKind::PermissionDenied => {
                    write!(
                        f,
                        "Error: Insufficient permissions to read '{}'.",
                        path_display
                    )
                }
                _ => write!(f, "Error reading file '{}': {}", path_display, source),
            },
            CvLoadError::Parse { source, .. } => {
                write!(
                    f,
                    "Error parsing '{}'. Please check for valid YAML syntax.\n Details: {}",
                    path_display, source
                )
            }
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[tokio::test]
    async fn load_cv_data_test() {
        // FIX: Temporarily move the working directory up to the workspace root
        // so the test can find "config/..."
        let _ = std::env::set_current_dir("..");

        let result = load_cv_config().await;
        if let Err(r) = &result {
            eprintln!("Error details: {:?}", r);
        }
        assert!(result.is_ok(), "Failed to load Projects CV data");
    }
}
