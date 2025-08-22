#include <iostream>
#include <thread>

#include "parse/YAMLProcessor.h"
#include "scheduler/scheduler.h"
#include "yaml-cpp/yaml.h"

int main() {
  unsigned int thread_count = std::thread::hardware_concurrency();

  // If hardware_concurrency() returns 0, we can assume a single thread.
  if (thread_count == 0) {
    thread_count = 1;
  }

  std::cout << "Main thread (" << std::this_thread::get_id()
            << ") will use a total of " << thread_count << " Logical Cores."
            << std::endl;

  BoostFiberScheduler scheduler(thread_count);

  YAMLProcessor yaml_processor;

  try {
    scheduler.addTask(&YAMLProcessor::parseGeneral, &yaml_processor);
    scheduler.addTask(&YAMLProcessor::parseExperience, &yaml_processor);
    scheduler.addTask(&YAMLProcessor::parseProject, &yaml_processor);

    scheduler.wait();
    std::cout << "\n--- All YAML files parsed successfully ---\n\n";

  } catch (const YAML::Exception &e) {
    std::cerr << "\n--- A YAML parsing task failed --- \n";
    std::cerr << "Error: " << e.what() << std::endl;
    return 1;

  } catch (const std::exception &e) {
    std::cerr << "\n--- An unexpected error occurred --- \n";
    std::cerr << "Error: " << e.what() << std::endl;
    return 1;
  }

  // std::cout << "\n--- General Information ---\n";
  // std::cout << yaml_processor.general << "\n";

  return 0;
}
