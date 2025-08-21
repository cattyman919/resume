#include "scheduler/scheduler.h"
#include "parse/YAMLProcessor.h"

#include <thread>
#include <iostream>

int main(){
   unsigned int thread_count = std::thread::hardware_concurrency();

    // If hardware_concurrency() returns 0, we can assume a single thread.
    if (thread_count == 0) {
        thread_count = 1;
    }

   std::cout << "Main thread (" << std::this_thread::get_id() << ") will use a total of " << thread_count << " threads." << std::endl;

    BoostFiberScheduler scheduler(thread_count);

    YAMLProcessor yaml_processor{};

    scheduler.addTask(&YAMLProcessor::parseGeneral, &yaml_processor);
    scheduler.addTask(&YAMLProcessor::parseExperience, &yaml_processor);
    scheduler.addTask(&YAMLProcessor::parseProject, &yaml_processor);

    scheduler.wait();
    std::cout << "\n--- All YAML files parsed successfully ---\n\n";

  // yaml_processor.general.print_skills_achivements();
  return 0;
}
