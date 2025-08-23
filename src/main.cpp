#include <functional>
#include <iostream>
#include <memory>
#include <thread>
#include <vector>

#include "concurrent_hash_set.h"
#include "parse/YAMLProcessor.h"
#include "scheduler/scheduler.h"
#include "yaml-cpp/yaml.h"

void getAllCVTypes(std::shared_ptr<ConcurrentHashSet<std::string>> shared_set,
                   const std::vector<Project>& projects) {
  for (const auto& proj : projects) {
    const auto& cv_types = proj.cv_type;
    for (const auto& cv_type : cv_types) {
      shared_set->insert(cv_type);
    }
  }
}

void getAllCVTypes(std::shared_ptr<ConcurrentHashSet<std::string>> shared_set,
                   const std::vector<Experience>& experiences) {
  for (const auto& exp : experiences) {
    const auto& cv_types = exp.cv_type;
    for (const auto& cv_type : cv_types) {
      shared_set->insert(cv_type);
    }
  }
}

int main() {
  unsigned int thread_count = std::thread::hardware_concurrency();

  // If hardware_concurrency() returns 0, we can assume a single thread.
  if (thread_count == 0) {
    thread_count = 1;
  }

  std::cout << "Main thread (" << std::this_thread::get_id()
            << ") will use a total of " << thread_count << " Logical Cores."
            << std::endl;

  BoostFiberScheduler executor(thread_count);

  General general{};
  std::vector<Experience> experiences{};
  std::vector<Project> projects{};

  try {
    executor.addTask(YAMLProcessor::parseGeneral, std::ref(general));
    executor.addTask(YAMLProcessor::parseExperience, std::ref(experiences));
    executor.addTask(YAMLProcessor::parseProject, std::ref(projects));

    executor.join();
    std::cout << "\n--- All YAML files parsed successfully ---\n\n";

  } catch (const YAML::Exception& e) {
    std::cerr << "\n--- A YAML parsing task failed --- \n";
    std::cerr << "Error: " << e.what() << std::endl;
    return 1;

  } catch (const std::exception& e) {
    std::cerr << "\n--- An unexpected error occurred --- \n";
    std::cerr << "Error: " << e.what() << std::endl;
    return 1;
  }

  // std::cout << "\n--- General Information ---\n";
  // std::cout << general << "\n";
  // for (const auto& exp : experiences) {
  //   std::cout << exp << "\n";
  // }
  //
  // for (const auto& proj : projects) {
  //   std::cout << proj << "\n";
  // }

  auto shared_set = std::make_shared<ConcurrentHashSet<std::string>>();

  try {
    executor.addTask([&] { getAllCVTypes(shared_set, projects); });
    executor.addTask([&] { getAllCVTypes(shared_set, experiences); });

    executor.join();

  } catch (const std::exception& e) {
    std::cerr << "\n--- An unexpected error occurred --- \n";
    std::cerr << "Error: " << e.what() << std::endl;
    return 1;
  }

  std::cout << "CV Types\n";

  for (const auto& item : shared_set->getSet()) {
    std::cout << "- " << item << '\n';
  }

  return 0;
}
