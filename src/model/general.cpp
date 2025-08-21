#include "model/general.h"
#include <iostream>
#include <type_traits>

std::ostream &operator<<(std::ostream &os, const PersonalInfo &p) {
  os << "Personal Info:\n";
  os << "Name: " << p.name << "\n";
  os << "Location: " << p.location << "\n";
  os << "Email: " << p.email << "\n";
  os << "Phone: " << p.phone << "\n";
  os << "Website: " << p.website << "\n";
  os << "LinkedIn: " << p.linkedin << "\n";
  os << "LinkedIn Handle: " << p.linkedin_handle << "\n";
  os << "GitHub: " << p.github << "\n";
  os << "GitHub Handle: " << p.github_handle << "\n";
  os << "Profile Pic: " << p.profile_pic << "\n";
  return os;
}

std::ostream &operator<<(std::ostream &os, const Education &e) {
  os << "Education:\n";
  os << "Institution: " << e.institution << "\n";
  os << "Degree: " << e.degree << "\n";
  os << "GPA: " << e.gpa << "\n";
  os << "Details: \n";
  for (const auto &detail : e.details) {
    os << "- " << detail << "\n";
  }
  return os;
}

std::ostream &operator<<(std::ostream &os, const Award &a) {
  os << "Award:\n";
  os << "Title: " << a.title << "\n";
  os << "Organization: " << a.organization << "\n";
  os << "Date: " << a.date << "\n";
  os << "Points: \n";
  for (const auto &point : a.points) {
    os << "- " << point << "\n";
  }
  return os;
}

std::ostream &operator<<(std::ostream &os, const Certificate &c) {
  os << "Certificate:\n";
  os << "Year: " << c.year << "\n";
  os << "Name: " << c.name << "\n";
  return os;
}

std::ostream &operator<<(std::ostream &os, const Skills &s) {
  os << "Skill: " << s.name << "\n";
  for (const auto &skill : s.skills) {
    os << "- " << skill << '\n';
  }
  return os;
}

std::ostream &operator<<(std::ostream &os,
                         const std::vector<General::AchievementItem> &sa) {

  os << "Skills & Achievements:\n\n";
  for (const auto &item : sa) {
    std::visit([&os](const auto &value) { os << value << "\n"; }, item);
  }
  return os;
}

std::ostream &operator<<(std::ostream &os, const std::vector<Education> &eds) {
  for (const auto &edu : eds) {
    os << edu << "\n";
  }
  return os;
}

std::ostream &operator<<(std::ostream &os, const std::vector<Award> &awa) {
  for (const auto &award : awa) {
    os << award << "\n";
  }
  return os;
}

std::ostream &operator<<(std::ostream &os, const General &g) {
  os << "General Information:\n";
  os << g.personal_info << "\n";
  os << g.skills_achievements << "\n";
  os << g.educations << "\n";
  os << g.awards << "\n";
  return os;
}

void General::print_skills_achivements() const {
  std::cout << "--- Achievements ---\n";
  for (const auto &item : this->skills_achievements) {
    // Use std::visit to safely handle the different types
    std::visit(
        [](const auto &value) {
          // 'if constexpr' checks the type at compile time
          if constexpr (std::is_same_v<decltype(value), const Skills &>) {
            std::cout << value.name << ":\n";
            for (const auto &skill : value.skills) {
              std::cout << "  - " << skill << "\n";
            }
          } else if constexpr (std::is_same_v<decltype(value),
                                              const Certificate &>) {
            std::cout << "Certificate:\n"; // Or you could use a different title
            std::cout << "  - " << value.name << " (" << value.year << ")\n";
          }
        },
        item);
  }
}
