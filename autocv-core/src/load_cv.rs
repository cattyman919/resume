use serde::de::DeserializeOwned;
use std::{error::Error, fmt, io};
use tokio::fs;

use crate::cv_model::{ExperiencesCVData, GeneralCVData, ProjectsCVData};

const GENERAL_CV_PATH: &str = "config/general.yaml";
const PROJECTS_CV_PATH: &str = "config/projects.yaml";
const EXPERIENCES_CV_PATH: &str = "config/experiences.yaml";

pub async fn load_cv_data()
-> Result<(GeneralCVData, ProjectsCVData, ExperiencesCVData), CvLoadError> {
    let (general_cv, projects_cv, experiences_cv) = tokio::try_join!(
        load_parse_yaml(GENERAL_CV_PATH),
        load_parse_yaml(PROJECTS_CV_PATH),
        load_parse_yaml(EXPERIENCES_CV_PATH)
    )?;

    Ok((general_cv, projects_cv, experiences_cv))
}

async fn load_parse_yaml<T: DeserializeOwned>(path: &str) -> Result<T, CvLoadError> {
    let content = fs::read_to_string(path)
        .await
        .map_err(|e| CvLoadError::Io {
            path: path.to_string(),
            source: e,
        })?;

    let data = serde_yaml::from_str(&content).map_err(|e| CvLoadError::Parse {
        path: path.to_string(),
        source: e,
    })?;

    return Ok(data);
}

#[derive(Debug)]
pub enum CvLoadError {
    Io {
        path: String,
        source: io::Error,
    },
    Parse {
        path: String,
        source: serde_yaml::Error,
    },
}

impl fmt::Display for CvLoadError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            CvLoadError::Io { path, source } => match source.kind() {
                io::ErrorKind::NotFound => {
                    write!(f, "Error: Configuration file '{}' was not found.", path)
                }
                io::ErrorKind::PermissionDenied => {
                    write!(f, "Error: Insufficient permissions to read '{}'.", path)
                }
                _ => write!(f, "Error reading file '{}': {}", path, source),
            },
            CvLoadError::Parse { path, source } => {
                write!(
                    f,
                    "Error parsing '{}'. Please check for valid YAML syntax.\n  Details: {}",
                    path, source
                )
            }
        }
    }
}

impl Error for CvLoadError {}

#[cfg(test)]
mod tests {
    use super::*;

    #[tokio::test]
    async fn load_cv_data_test() {
        let result = load_cv_data().await;
        if let Err(r) = &result {
            eprintln!("Error details: {:?}", r);
        }
        assert!(result.is_ok(), "Failed to load Projects CV data");
    }
}
