use anyhow::Result;
use futures::future::join_all;
use std::{env, error::Error, process, sync::Arc, time::Instant};

use autocv_core::{self, cv_processor, load_cv};

struct AppConfig {
    is_benchmark_mode: bool,
    is_debug_mode: bool,
}

impl AppConfig {
    fn new() -> Self {
        let args: Vec<String> = env::args().collect();
        Self {
            is_benchmark_mode: args.contains(&"--benchmark".to_string()),
            is_debug_mode: args.contains(&"--debug".to_string()),
        }
    }
}

async fn run() -> Result<(), Box<dyn Error>> {
    let config = AppConfig::new();
    let start_time = if config.is_benchmark_mode {
        Some(Instant::now())
    } else {
        None
    };

    println!("\n==== Generating All LaTeX CV ====\n");

    cv_processor::setup_directories().await?;

    println!("Loading YAML Data...");
    let (general_cv, projects_cv, experiences_cv) = load_cv::load_cv_data().await?;

    let general_cv = Arc::new(general_cv);
    let projects_cv = Arc::new(projects_cv);
    let experiences_cv = Arc::new(experiences_cv);

    println!("Getting Total CV Types...");
    let all_cv_types = cv_processor::get_all_cv_types(&projects_cv, &experiences_cv).await?;
    println!("All CV Types: {:?}", all_cv_types);

    let processing_tasks = all_cv_types.into_iter().map(|cv_type| {
        let general_cv_clone = Arc::clone(&general_cv);
        let projects_cv_clone = projects_cv.as_ref().clone();
        let experiences_cv_clone = experiences_cv.as_ref().clone();
        tokio::spawn(async move {
            println!("Processing CV type: {}", cv_type);
            cv_processor::write_cv(
                general_cv_clone,
                projects_cv_clone,
                experiences_cv_clone,
                cv_type,
                config.is_debug_mode,
            )
            .await
        })
    });

    let results = join_all(processing_tasks).await;

    // --- Error Handling for Concurrent Tasks ---
    let mut errors = Vec::new();
    for result in results {
        match result {
            // The task itself panicked (a serious bug).
            Err(join_error) => errors.push(format!("A task panicked: {}", join_error)),
            // The task completed but returned an error.
            Ok(Err(task_error)) => errors.push(format!("A task failed: {}", task_error)),
            Ok(Ok(_)) => (),
        }
    }

    if !errors.is_empty() {
        eprintln!("\nErrors occurred during CV generation:");
        for e in errors {
            eprintln!("- {}", e);
        }
        return Err("CV generation failed due to one or more task errors.".into());
    }

    cv_processor::move_aux_files().await?;

    println!("\n==== All LaTeX CV Generation Complete ====");

    if let Some(start) = start_time {
        println!("Total time taken: {:?}", start.elapsed());
    }

    Ok(())
}

#[tokio::main]
async fn main() {
    if let Err(e) = run().await {
        eprintln!("\n{}", e);
        process::exit(1);
    }
}
