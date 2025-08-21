#include "parse/YAMLProcessor.h"
#include "model/experience.h"
#include "model/general.h"
#include "yaml-cpp/yaml.h"
#include <thread>
#include <variant>
#include <vector>

void setPersonalInfo(YAMLProcessor &yaml_processor, YAML::Node &general_node) {
  yaml_processor.general.personal_info.name =
      general_node["personal_info"]["name"].as<std::string>();
  yaml_processor.general.personal_info.location =
      general_node["personal_info"]["location"].as<std::string>();
  yaml_processor.general.personal_info.email =
      general_node["personal_info"]["email"].as<std::string>();
  yaml_processor.general.personal_info.phone =
      general_node["personal_info"]["phone"].as<std::string>();
  yaml_processor.general.personal_info.website =
      general_node["personal_info"]["website"].as<std::string>();
  yaml_processor.general.personal_info.linkedin =
      general_node["personal_info"]["linkedin"].as<std::string>();
  yaml_processor.general.personal_info.linkedin_handle =
      general_node["personal_info"]["linkedin_handle"].as<std::string>();
  yaml_processor.general.personal_info.github =
      general_node["personal_info"]["github"].as<std::string>();
  yaml_processor.general.personal_info.github_handle =
      general_node["personal_info"]["github_handle"].as<std::string>();
  yaml_processor.general.personal_info.profile_pic =
      general_node["personal_info"]["profile_pic"].as<std::string>();
}

void YAMLProcessor::parseGeneral() {
  std::thread::id thread_id = std::this_thread::get_id();
  std::stringstream ss;
  ss << "Thread (" << thread_id << ") :  Parsing general.yaml...\n";
  std::cout << ss.str();

  YAML::Node general_node = YAML::LoadFile("config/general.yaml");

  if (general_node) {
    setPersonalInfo(*this, general_node);

    if (auto skills_achievements_node = general_node["skills_achievements"];
        skills_achievements_node && skills_achievements_node.IsMap()) {

      for (const auto &pair : skills_achievements_node) {
        std::string category_name = pair.first.as<std::string>();
        const YAML::Node &category_node = pair.second;

        if (category_name == "Certificates" && category_node.IsSequence()) {
          for (const auto &item : category_node) {
            if (item.IsMap() && item["year"] && item["name"]) {
              this->general.skills_achievements.emplace_back(
                  std::in_place_type<Certificate>,
                  item["year"].as<unsigned int>(),
                  item["name"].as<std::string>());
            }
          }
        } else if (category_node.IsSequence()) {

          this->general.skills_achievements.emplace_back(
              std::in_place_type<Skills>, category_name,
              category_node.as<std::vector<std::string>>());
        }
      }
    }

    if (auto education_node = general_node["education"];
        education_node && education_node.IsSequence()) {
      for (const auto &edu_node : education_node) {
        this->general.educations.emplace_back(
            edu_node["institution"].as<std::string>(),
            edu_node["degree"].as<std::string>(),
            edu_node["dates"].as<std::string>(),
            edu_node["gpa"].as<std::string>(),
            edu_node["details"] && edu_node["details"].IsSequence()
                ? edu_node["details"].as<std::vector<std::string>>()
                : std::vector<std::string>{});
      }
    }

    if (auto awards_node = general_node["awards"];
        awards_node && awards_node.IsSequence()) {

      for (const auto &award_node : awards_node) {
        this->general.awards.emplace_back(
            award_node["title"].as<std::string>(),
            award_node["organization"].as<std::string>(),
            award_node["date"].as<std::string>(),
            award_node["points"] && award_node["points"].IsSequence()
                ? award_node["points"].as<std::vector<std::string>>()
                : std::vector<std::string>{});
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

    for (const YAML::Node &project_node : projects_node) {
      std::vector<std::string> cv_types;
      if (project_node["cv_type"] && project_node["cv_type"].IsSequence()) {
        cv_types = project_node["cv_type"].as<std::vector<std::string>>();
      }

      projects.emplace_back(
          project_node["name"].as<std::string>(),
          project_node["github"].as<std::string>(),
          project_node["github_handle"].as<std::string>(), std::move(cv_types),
          project_node["description"].as<std::string>(),
          project_node["points"].as<std::vector<std::string>>());
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

    for (const YAML::Node &experience_node : experiences_node) {

      std::vector<std::string> cv_types;
      if (experience_node["cv_type"] &&
          experience_node["cv_type"].IsSequence()) {
        cv_types = experience_node["cv_type"].as<std::vector<std::string>>();
      }

      experiences.emplace_back(
          experience_node["company"].as<std::string>(),
          experience_node["location"].as<std::string>(),
          experience_node["role"].as<std::string>(),
          experience_node["dates"].as<std::string>(),
          experience_node["job_type"].as<std::string>(), std::move(cv_types),
          experience_node["points"].as<std::vector<std::string>>());
    }
  }
}
