#include "model/project.h"
#include "model/general.h"
#include "model/experience.h"
#include <vector>

class YAMLProcessor {
  public:
    General general {};
    std::vector<Project> projects {};
    std::vector<Experience> experiences {};

    YAMLProcessor() = default;

    explicit YAMLProcessor(const General& general,
                  const std::vector<Project>& projects,
                  const std::vector<Experience>& experiences)
        : general(general), projects(projects), experiences(experiences) {}

    void parseGeneral();
    void parseProject();
    void parseExperience();
};
