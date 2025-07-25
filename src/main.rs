mod config;
mod cv_processor;
mod types;

use crate::types::CVData;
use std::process;

#[tokio::main]
async fn main() {
    println!("\n==== Generating All LaTeX CV ====\n");

    println!("Loading YAML Data...\n");
    let cv_data: CVData = match config::parse_get_cv().await {
        Ok(data) => data,
        Err(e) => {
            eprintln!("{e}");
            process::exit(1)
        }
    };

    println!("Getting Total CV Types...\n");
    let all_types = match cv_processor::get_all_cv_types(&cv_data).await {
        Ok(data) => data,
        Err(e) => {
            eprintln!("{e}");
            process::exit(1)
        }
    };

    println!("All CV Types: {:?}", all_types);
}
