mod config;
mod cv_processor;
mod types;

use crate::types::CVData;
use std::{io::ErrorKind, process};
use tokio::fs;

#[tokio::main]
async fn main() {
    println!("\n==== Generating All LaTeX CV ====\n");

    println!("Creating out & cv directory...");

    let out_directory_task = async {
        if let Err(e) = fs::create_dir("out").await {
            if e.kind() != ErrorKind::AlreadyExists {
                panic!("Failed to create 'out' directory: {}", e);
            }
        }
    };

    let cv_directory_task = async {
        if let Err(e) = fs::create_dir("cv").await {
            if e.kind() != ErrorKind::AlreadyExists {
                panic!("Failed to create 'cv' directory: {}", e);
            }
        }
    };

    tokio::join!(out_directory_task, cv_directory_task);
    println!("Directories are ready.");

    println!("Loading YAML Data...");
    let cv_data: CVData = match config::parse_get_cv() {
        Ok(data) => data,
        Err(e) => {
            eprintln!("{e}");
            process::exit(1)
        }
    };

    println!("Getting Total CV Types...");
    let all_cv_types = match cv_processor::get_all_cv_types(&cv_data).await {
        Ok(data) => data,
        Err(e) => {
            eprintln!("{e}");
            process::exit(1)
        }
    };

    println!("All CV Types: {:?}", all_cv_types);

    for cv_type in all_cv_types.iter() {
        println!("Processing CV type: {cv_type}");
        match cv_processor::write_cv(cv_type).await {
            Ok(_res) => (),
            Err(e) => {
                eprintln!("{e}");
                process::exit(1)
            }
        };
    }
}
