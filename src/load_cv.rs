use std::{error::Error, fmt, fs, io};

use crate::types::CVData;

const CV_FILE_PATH: &str = "cv_data.yaml";

pub fn load_cv_data() -> Result<CVData, CvLoadError> {
    let config_cv_str = fs::read_to_string(CV_FILE_PATH)?;
    let config = serde_yaml::from_str(&config_cv_str)?;
    Ok(config)
}

#[derive(Debug)]
pub enum CvLoadError {
    Io(io::Error),
    Parse(serde_yaml::Error),
}

impl fmt::Display for CvLoadError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            CvLoadError::Io(e) => match e.kind() {
                io::ErrorKind::NotFound => {
                    write!(
                        f,
                        "Error: Configuration file '{}' was not found.",
                        CV_FILE_PATH
                    )
                }
                io::ErrorKind::PermissionDenied => {
                    write!(
                        f,
                        "Error: Insufficient permissions to read '{}'.",
                        CV_FILE_PATH
                    )
                }
                _ => write!(f, "Error reading file '{}': {}", CV_FILE_PATH, e),
            },
            CvLoadError::Parse(e) => {
                write!(
                    f,
                    "Error parsing '{}'. Please check for valid YAML syntax.\n  Details: {}",
                    CV_FILE_PATH, e
                )
            }
        }
    }
}

impl Error for CvLoadError {}

impl From<io::Error> for CvLoadError {
    fn from(err: io::Error) -> CvLoadError {
        CvLoadError::Io(err)
    }
}

impl From<serde_yaml::Error> for CvLoadError {
    fn from(err: serde_yaml::Error) -> Self {
        CvLoadError::Parse(err)
    }
}
