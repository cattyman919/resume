use std::{
    collections::HashMap,
    fs::{self},
    path::PathBuf,
};

use lazy_static::lazy_static;
use minijinja::{Environment, syntax::SyntaxConfig};
use regex::Regex;

macro_rules! template_path {
    ($filename:expr) => {
        // env!("CARGO_MANIFEST_DIR") gives the absolute path to your crate folder at compile time
        [
            env!("CARGO_MANIFEST_DIR"),
            "..",
            "template_cv",
            "type1_cv",
            "sections",
            $filename,
        ]
        .iter()
        .collect::<PathBuf>()
    };
}

lazy_static! {
static ref BOLD_REGEX: Regex = Regex::new(r"\*\*(.*?)\*\*").unwrap();
pub static ref TEMPLATE_PATHS: HashMap<&'static str, PathBuf> = {
    let mut m = HashMap::new();

    m.insert(
        "Achievement_Skills",
        template_path!("Achievement_Skills.tex"),
    );
    m.insert("Awards", template_path!("Awards.tex"));
    m.insert("Education", template_path!("Education.tex"));
    m.insert("Experience", template_path!("Experience.tex"));
    m.insert("Header", template_path!("Header.tex"));
    m.insert("Projects", template_path!("Projects.tex"));
    m
};
pub static ref TMPL_ENV: Environment<'static> = {

// Configure it to use (( )) instead of {{ }}
let mut environment = Environment::new();

environment.add_filter("latex", |value: String| -> String {
    format_latex_bold(&value)
    });

environment.set_syntax(SyntaxConfig::builder()
    .block_delimiters("((*", "*))")
    .variable_delimiters("((", "))")
    // .comment_delimiters("\\#{", "}")
    .build()
    .unwrap()
);
    for (key,val) in TEMPLATE_PATHS.iter() {

        let content = fs::read_to_string(val).expect("Failed to read file");
        environment.add_template_owned(*key, content).unwrap();
    }
    environment
    };
}

/// Escapes special LaTeX characters in a string.
/// The order of replacement is crucial, especially for the backslash.
pub fn escape_latex(text: &str) -> String {
    text
        // FIX: Backslash must be escaped FIRST.
        .replace('\\', r"\textbackslash{}")
        .replace('&', r"\&")
        .replace('%', r"\%")
        .replace('$', r"\$")
        .replace('#', r"\#")
        .replace('_', r"\_")
        .replace('{', r"\{")
        .replace('}', r"\}")
        .replace('~', r"\textasciitilde{}")
        .replace('^', r"\textasciicircum{}")
}

/// Finds markdown-style bold text (**text**) and converts it to LaTeX \textbf{text}.
fn format_latex_bold(text: &str) -> String {
    let escaped_text = escape_latex(text);
    BOLD_REGEX
        .replace_all(&escaped_text, r"\textbf{$1}")
        .to_string()
}

#[cfg(test)]
mod tests {
    use tokio::fs::File;

    use crate::cv::model::Experience;
    use crate::cv::model::PersonalInfo;
    use minijinja::context;

    use super::*;

    #[tokio::test]
    async fn generate_header_template() {
        let _ = std::env::set_current_dir("..");
        let mut personal_info = PersonalInfo::default();
        personal_info.name = "Seno Pamungas".into();
        personal_info.email = "skibidi".into();
        let output_file = File::create("target/test.tex")
            .await
            .unwrap()
            .into_std()
            .await;

        let tmpl = TMPL_ENV
            .get_template("Header")
            .expect("Template 'Header' not found in TMPL_ENV");
        let result = tmpl.render_to_write(&personal_info, output_file);
        if let Err(ref e) = result {
            eprintln!("Rendering failed: {:#?}", e);
        }
        assert!(result.is_ok(), "Failed to render template to file");
    }

    #[tokio::test]
    async fn generate_experience_template() {
        // 1. Prepare Data
        let exp1 = Experience {
            role: "Senior Rustacean".to_string(),
            company: "Crab Corp".to_string(),
            job_type: "Full-time".to_string(),
            location: "The Sea".to_string(),
            dates: "Jan 2024 -- Present".to_string(),
            points: vec![
                "Rewrote the entire backend in **Rust**".to_string(),
                "Reduced memory usage by 90%".to_string(),
            ],
            cv_type: vec![], // Optional
        };

        let exp2 = Experience {
            role: "Go Developer".to_string(),
            company: "Gopher Inc".to_string(),
            job_type: "Contract".to_string(),
            location: "Mountain View".to_string(),
            dates: "Jan 2020 -- Dec 2023".to_string(),
            points: vec!["Built microservices".to_string()],
            cv_type: vec![],
        };

        let experiences = vec![exp1, exp2];

        // 2. Setup Output File
        let output_file = File::create("target/test_experience.tex")
            .await
            .unwrap()
            .into_std()
            .await;

        // 3. Render
        // We use context! to map the "experiences" key in the template
        // to our rust vector.
        let tmpl = TMPL_ENV
            .get_template("Experience")
            .expect("Template not found");

        let result = tmpl.render_to_write(context! { experiences => experiences }, output_file);

        // 4. Assertions
        if let Err(ref e) = result {
            eprintln!("Rendering failed: {:#?}", e);
        }
        assert!(result.is_ok());

        // Optional: Verify content content
        let content = std::fs::read_to_string("target/test_experience.tex").unwrap();
        assert!(content.contains("Crab Corp"));
        assert!(content.contains(r"Rewrote the entire backend in \textbf{Rust}"));
    }
}
