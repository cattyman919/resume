mod config;
mod cv_processor;
mod latex;
mod types;

use futures::future::join_all;
use std::{env, io::ErrorKind, process, sync::Arc};
use tokio::fs;
use types::CVData;

#[tokio::main]
async fn main() {
    println!("\n==== Generating All LaTeX CV ====\n");

    let debug_mode = env::args().any(|arg| arg == "--debug");

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

    let cv_data = Arc::new(cv_data);

    println!("Getting Total CV Types...");
    let all_cv_types = match cv_processor::get_all_cv_types(&cv_data).await {
        Ok(data) => data,
        Err(e) => {
            eprintln!("{e}");
            process::exit(1)
        }
    };

    println!("All CV Types: {:?}", all_cv_types);
    let handles = all_cv_types
        .into_iter()
        .map(|cv_type| {
            let cv_data_clone = Arc::clone(&cv_data);
            tokio::spawn(async move {
                println!("Processing CV type: {cv_type}");
                match cv_processor::write_cv(cv_data_clone, cv_type, debug_mode).await {
                    Ok(_res) => (),
                    Err(e) => {
                        eprintln!("{e}");
                        process::exit(1)
                    }
                };
            })
        })
        .collect::<Vec<_>>();
    let _results = join_all(handles).await;

    let _ = cv_processor::move_aux_files().await;

    println!("\n==== All LaTeX CV Generation Complete ====\n")
}
