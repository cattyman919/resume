#pragma once

#include <iostream>
#include <string>
#include <variant>
#include <vector>

class PersonalInfo
{
  public:
    std::string name{};
    std::string location{};
    std::string email{};
    std::string phone{};
    std::string website{};
    std::string linkedin{};
    std::string linkedin_handle{};
    std::string github{};
    std::string github_handle{};
    std::string profile_pic{};

    PersonalInfo() = default;

    explicit PersonalInfo(std::string name,
        std::string location,
        std::string email,
        std::string phone,
        std::string website,
        std::string linkedin,
        std::string linkedin_handle,
        std::string github,
        std::string github_handle,
        std::string profile_pic) :
      name(name),
      location(location),
      email(email),
      phone(phone),
      website(website),
      linkedin(linkedin),
      linkedin_handle(linkedin_handle),
      github(github),
github_handle(github_handle),
      profile_pic(profile_pic) {}

    friend std::ostream& operator<<(std::ostream& os, const PersonalInfo& p);
};

struct Education {
    const std::string institution{};
    const std::string degree{};
    const std::string dates{};
    const std::string gpa{};
    const std::vector<std::string> details{};
};

struct Award {
    const std::string title{};
    const std::string organization{};
    const std::string date{};
    const std::vector<std::string> points{};
};

struct Certificate {
    const unsigned int year{};
    const std::string name{};
};

struct Skills {
  std::string name{};
  std::vector<std::string> skills{};
};

class General {
  public:

  using AchievementItem = std::variant<Skills, Certificate>;

  PersonalInfo personal_info{};
  std::vector<AchievementItem> skills_achievements{};
  std::vector<Education> educations{};
  std::vector<Award> awards{};

  General() = default;

  explicit General(PersonalInfo personal_info,
      std::vector<AchievementItem> skills_achievements,
      std::vector<Education> educations,
      std::vector<Award> awards) :
    personal_info(personal_info),
    skills_achievements(skills_achievements),
    educations(educations),
    awards(awards) {}

    void print_skills_achivements() const;
};
