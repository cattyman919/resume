#pragma once

#include <string>
#include <vector>
class Experience
{
  public:
    std::string company{};
    std::string location{};
    std::string role{};
    std::string dates{};
    std::string job_type{};
    std::vector<std::string> cv_type{};
    std::vector<std::string> points{};

    Experience() = default;

    explicit Experience(std::string company,
        std::string location,
        std::string role,
        std::string dates,
        std::string job_type,
        std::vector<std::string> cv_type,
        std::vector<std::string> points
        ) :
      company(company),
      location(location),
      role(role),
      dates(dates),
      job_type(role),
      cv_type(cv_type),
      points(points) {}

    friend std::ostream& operator<<(std::ostream& os, const Experience& e);
};
