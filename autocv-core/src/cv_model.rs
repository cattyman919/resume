use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize, Clone, Default)]
pub struct PersonalInfo {
    pub name: String,
    pub email: String,
    pub phone: String,
    pub website: String,
    pub linkedin: String,
    pub linkedin_handle: String,
    pub github: String,
    pub github_handle: String,
    pub profile_pic: String,
    pub location: String,
}

#[derive(Debug, Serialize, Deserialize, Clone, Default)]
pub struct Experience {
    pub role: String,
    pub job_type: String,
    pub company: String,
    pub location: String,
    pub dates: String,
    pub points: Vec<String>,
    pub cv_type: Vec<String>,
}

#[derive(Debug, Serialize, Deserialize, Clone, Default)]
pub struct Education {
    pub institution: String,
    pub degree: String,
    pub dates: String,
    pub gpa: String,
    pub details: Vec<String>,
}

#[derive(Debug, Serialize, Deserialize, Clone, Default)]
pub struct Award {
    pub title: String,
    pub organization: String,
    pub date: String,
    pub points: Vec<String>,
}

#[derive(Debug, Serialize, Deserialize, Clone, Default)]
pub struct Project {
    pub name: String,
    pub github: String,
    pub github_handle: String,
    pub cv_type: Vec<String>,
    pub points: Vec<String>,
}

#[derive(Debug, Serialize, Deserialize, Clone, Default)]
pub struct Certificate {
    pub name: String,
    pub year: String,
}

#[derive(Debug, Serialize, Deserialize, Clone, Default)]
pub struct SkillsAchievements {
    #[serde(rename = "Hard Skills")]
    pub hard_skills: Vec<String>,

    #[serde(rename = "Soft Skills")]
    pub soft_skills: Vec<String>,

    #[serde(rename = "Programming Languages")]
    pub programming_languages: Vec<String>,

    #[serde(rename = "Databases")]
    pub databases: Vec<String>,

    #[serde(rename = "Misc")]
    pub misc: Vec<String>,

    #[serde(rename = "Certificates")]
    pub certificates: Vec<Certificate>,
}

#[derive(Debug, Serialize, Deserialize, Clone, Default)]
pub struct GeneralCVData {
    pub personal_info: PersonalInfo,
    pub education: Vec<Education>,
    pub awards: Vec<Award>,
    pub skills_achievements: SkillsAchievements,
}

pub type ProjectsCVData = Vec<Project>;

pub type ExperiencesCVData = Vec<Experience>;

pub trait HasCvTypes {
    fn cv_types(&self) -> &Vec<String>;
}

impl HasCvTypes for Experience {
    fn cv_types(&self) -> &Vec<String> {
        &self.cv_type
    }
}

impl HasCvTypes for Project {
    fn cv_types(&self) -> &Vec<String> {
        &self.cv_type
    }
}
