#pragma once

#include <string>
#include <vector>

struct Project
{
    std::string name{};
    std::string github{};
    std::string github_handle{};
    std::vector<std::string> cv_type{};
    std::string description{};
    std::vector<std::string> points{};

    Project() = default;

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
