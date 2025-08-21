#include "parse/YAMLProcessor.h"

int main(){
  YAMLProcessor yaml_processor{};
  yaml_processor.parseGeneral();
  yaml_processor.parseProject();
  yaml_processor.parseExperience();
  yaml_processor.general.print_skills_achivements();
  return 0;
}
