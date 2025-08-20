#include <iostream>
#include "yaml-cpp/yaml.h"

struct Project
{
    std::string name {};
    std::string github {};
    std::string github_handle {};
    std::vector<std::string> cv_type {};
    std::string description {};
    std::vector<std::string> points {};

    Project() = default;

    // Constructor
    explicit Project(std::string name,
        std::string github,
        std::string github_handle,
        std::vector<std::string> cv_type,
        std::string description,
        std::vector<std::string> points) :
      name(name),
      github(github),
      github_handle(github_handle),
      cv_type(cv_type),
      description(description),
      points(points) {}

    friend std::ostream& operator<<(std::ostream& os, const Project& p);
};

std::ostream& operator<<(std::ostream& os, const Project& p) {
    os << "Project { name: " << p.name << '\n'
       << "github: " << p.github << '\n'
       << "github_handle: " << p.github_handle << '\n'
       << "cv_type: [ " ;
        for (const auto& type : p.cv_type){
          os << type << ", ";
        }
       os << " ]\n"
       << "description: " << p.description << '\n'
       << "points: [ " ;
        for (const auto& point : p.points){
          os << point << "\n";
        }
       os << "]"  << " }";
    return os;
}

int main(){
  std::vector<Project> projects;

    try {
        YAML::Node projects_node = YAML::LoadFile("config/projects.yaml");

        if (projects_node && projects_node.IsSequence()) {
            Project project;
            for (const YAML::Node& project_node : projects_node) {
                std::vector<std::string> cv_types ;
                if (project_node["cv_type"] && project_node["cv_type"].IsSequence()) {
                    cv_types = project_node["cv_type"].as<std::vector<std::string>>();
                }
                project.name = project_node["name"].as<std::string>();
                project.github =  project_node["github"].as<std::string>();
                project.github_handle =  project_node["github_handle"].as<std::string>();
                project.cv_type =  std::move(cv_types);
                project.description =  project_node["description"].as<std::string>();
                project.points =  project_node["points"].as<std::vector<std::string>>();

                std::cout << "- " << project << "\n\n";

                projects.emplace_back(std::move(project));
            }
        }
        std::cout << "Total Project: " << projects.size() << '\n';

    } catch (const YAML::Exception& e) {
        std::cerr << "Error parsing YAML: " << e.what() << std::endl;
        std::exit(1);
    }
  return 0;
}
