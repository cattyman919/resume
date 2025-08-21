#include "parse/YAMLProcessor.h"
#include "model/experience.h"
#include "model/general.h"
#include "yaml-cpp/yaml.h"
#include <thread>
#include <vector>

void YAMLProcessor::parseGeneral() {
  std::thread::id thread_id = std::this_thread::get_id();
  std::stringstream ss;
  ss << "Thread (" << thread_id << ") :  Parsing general.yaml...\n";
  std::cout << ss.str();
  YAML::Node general_node = YAML::LoadFile("config/general.yaml");

  if (general_node) {
    this->general.personal_info.name =
        general_node["personal_info"]["name"].as<std::string>();
    this->general.personal_info.location =
        general_node["personal_info"]["location"].as<std::string>();
    this->general.personal_info.email =
        general_node["personal_info"]["email"].as<std::string>();
    this->general.personal_info.phone =
        general_node["personal_info"]["phone"].as<std::string>();
    this->general.personal_info.website =
        general_node["personal_info"]["website"].as<std::string>();
    this->general.personal_info.linkedin =
        general_node["personal_info"]["linkedin"].as<std::string>();
    this->general.personal_info.linkedin_handle =
        general_node["personal_info"]["linkedin_handle"].as<std::string>();
    this->general.personal_info.github =
        general_node["personal_info"]["github"].as<std::string>();
    this->general.personal_info.github_handle =
        general_node["personal_info"]["github_handle"].as<std::string>();
    this->general.personal_info.profile_pic =
        general_node["personal_info"]["profile_pic"].as<std::string>();

    if (auto skills_achievements_node = general_node["skills_achievements"];
        skills_achievements_node && skills_achievements_node.IsMap()) {
      for (const auto &pair : skills_achievements_node) {
        std::string category_name = pair.first.as<std::string>();
        const YAML::Node &category_node = pair.second;

        if (category_name == "Certificates" && category_node.IsSequence()) {
          for (const auto &item : category_node) {
            if (item.IsMap() && item["year"] && item["name"]) {
              Certificate cert{item["year"].as<unsigned int>(),
                               item["name"].as<std::string>()};
              this->general.skills_achievements.emplace_back(cert);
            }
          }
        } else if (category_node.IsSequence()) {
          Skills skills_item;
          skills_item.name = category_name;

          skills_item.skills = category_node.as<std::vector<std::string>>();

          this->general.skills_achievements.emplace_back(skills_item);
        }
      }
    }

    if (auto education_node = general_node["education"];
        education_node && education_node.IsSequence()) {
      for (const auto &edu_node : education_node) {
        Education education{
            edu_node["institution"].as<std::string>(),
            edu_node["degree"].as<std::string>(),
            edu_node["dates"].as<std::string>(),
            edu_node["gpa"].as<std::string>(),
            edu_node["details"] && edu_node["details"].IsSequence()
                ? edu_node["details"].as<std::vector<std::string>>()
                : std::vector<std::string>{}};
        this->general.educations.emplace_back(education);
      }
    }

    if (auto awards_node = general_node["awards"];
        awards_node && awards_node.IsSequence()) {
      for (const auto &award_node : awards_node) {
        Award award{award_node["title"].as<std::string>(),
                    award_node["organization"].as<std::string>(),
                    award_node["date"].as<std::string>(),
                    award_node["points"] && award_node["points"].IsSequence()
                        ? award_node["points"].as<std::vector<std::string>>()
                        : std::vector<std::string>{}};
        this->general.awards.emplace_back(award);
      }
    }
  }
}

void YAMLProcessor::parseProject() {
  std::thread::id thread_id = std::this_thread::get_id();
  std::stringstream ss;
  ss << "Thread (" << thread_id << ") :  Parsing project.yaml...\n";
  std::cout << ss.str();

  YAML::Node projects_node = YAML::LoadFile("config/projects.yaml");

  if (projects_node && projects_node.IsSequence()) {
    Project project;
    for (const YAML::Node &project_node : projects_node) {

      std::vector<std::string> cv_types;
      if (project_node["cv_type"] && project_node["cv_type"].IsSequence()) {
        cv_types = project_node["cv_type"].as<std::vector<std::string>>();
      }

      project.name = project_node["name"].as<std::string>();
      project.github = project_node["github"].as<std::string>();
      project.github_handle = project_node["github_handle"].as<std::string>();
      project.cv_type = std::move(cv_types);
      project.description = project_node["description"].as<std::string>();
      project.points = project_node["points"].as<std::vector<std::string>>();

#ifdef DEBUG
      std::cout << "- " << project << "\n\n";
#endif // DEBUG

      projects.emplace_back(std::move(project));
    }
  }
}

void YAMLProcessor::parseExperience() {
  std::thread::id thread_id = std::this_thread::get_id();
  std::stringstream ss;
  ss << "Thread (" << thread_id << ") :  Parsing experiences.yaml...\n";
  std::cout << ss.str();
  YAML::Node experiences_node = YAML::LoadFile("config/experiences.yaml");

  if (experiences_node && experiences_node.IsSequence()) {
    Experience experience;
    for (const YAML::Node &experience_node : experiences_node) {

      std::vector<std::string> cv_types;
      if (experience_node["cv_type"] &&
          experience_node["cv_type"].IsSequence()) {
        cv_types = experience_node["cv_type"].as<std::vector<std::string>>();
      }

      experience.company = experience_node["company"].as<std::string>();
      experience.location = experience_node["location"].as<std::string>();
      experience.role = experience_node["role"].as<std::string>();
      experience.dates = experience_node["dates"].as<std::string>();
      experience.job_type = experience_node["job_type"].as<std::string>();
      experience.cv_type = std::move(cv_types);
      experience.points =
          experience_node["points"].as<std::vector<std::string>>();
#ifdef DEBUG
      std::cout << "- " << experience << "\n\n";
#endif // DEBUG
      experiences.emplace_back(std::move(experience));
    }
  }
}
