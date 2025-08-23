#include "model/experience.h"
#include "model/general.h"
#include "model/project.h"

class YAMLProcessor final {
 public:
  static void parseGeneral(General &general);
  static void parseProject(std::vector<Project> &projects);
  static void parseExperience(std::vector<Experience> &experiences);
};
