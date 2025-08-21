#include "model/experience.h"
#include "model/general.h"
#include "model/project.h"
#include <vector>

class YAMLProcessor {
public:
  General general{};
  std::vector<Project> projects{};
  std::vector<Experience> experiences{};

  YAMLProcessor() = default;

  YAMLProcessor(const General &general, const std::vector<Project> &projects,
                const std::vector<Experience> &experiences)
      : general(general), projects(projects), experiences(experiences) {}

  void parseGeneral();
  void parseProject();
  void parseExperience();
};
