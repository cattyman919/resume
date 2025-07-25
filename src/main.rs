mod types;
use std::{io, process};

use crate::types::CVData;

const CONFIG_CV_FILE_PATH: &str = "cv_data.yaml";

#[tokio::main]
async fn main() {
    let cv_data: CVData = match parse_get_cv().await {
        Ok(data) => data,
        Err(e) => {
            eprintln!("{e}");
            process::exit(1)
        }
    };
    println!("{}", cv_data.personal_info.name)
}

async fn parse_get_cv() -> Result<types::CVData, String> {
    let config_cv_str = match tokio::fs::read_to_string(CONFIG_CV_FILE_PATH).await {
        Ok(content) => content,
        Err(e) => {
            let error_message = match e.kind() {
                io::ErrorKind::NotFound => {
                    format!("Error: Configuration file '{CONFIG_CV_FILE_PATH}' does not exist")
                }
                io::ErrorKind::PermissionDenied => {
                    format!("Error: Insufficient permission to read '{CONFIG_CV_FILE_PATH}'")
                }
                _ => format!("An unexpected error reading '{CONFIG_CV_FILE_PATH}'"),
            };
            return Err(error_message);
        }
    };
    let config: types::CVData = match serde_yaml::from_str(&config_cv_str) {
        Ok(parsed_config) => parsed_config,
        Err(e) => {
            let error_message = format!(
                "Error: Failed to parse '{CONFIG_CV_FILE_PATH}'\nPlease check thhe file for correct YAML syntax.\nDetails:\n{e}"
            );
            return Err(error_message);
        }
    };
    Ok(config)
}
