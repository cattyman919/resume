#include "model/project.h"
#include <iostream>

std::ostream &operator<<(std::ostream &os, const Project &p) {
  os << "Project {\nname: " << p.name << '\n'
     << "github: " << p.github << '\n'
     << "github_handle: " << p.github_handle << '\n'
     << "cv_type: [ ";
  for (const auto &type : p.cv_type) {
    os << type << ", ";
  }
  os << " ]\n"
     << "description: " << p.description << '\n'
     << "points: [\n";
  for (const auto &point : p.points) {
    os << "- " << point << "\n";
  }
  os << "]" << " }";
  return os;
}
