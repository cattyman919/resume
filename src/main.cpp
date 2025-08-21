#include "parse/YAMLProcessor.h"
#include "boost/fiber/all.hpp"
#include <thread>
#include <iostream>

int main(){
   unsigned int thread_count = 3;

   std::cout << "Main thread (" << std::this_thread::get_id() << ") will use a total of " << thread_count << " threads." << std::endl;

  boost::fibers::barrier b(thread_count);
  std::vector<std::thread> threads;

  // Launch Worker threads
  for (unsigned int i = 1; i < thread_count; ++i) {
        threads.emplace_back([&b, thread_count]() {
            boost::fibers::use_scheduling_algorithm<boost::fibers::algo::work_stealing>(thread_count);
            b.wait();
        });
   }

    std::vector<boost::fibers::fiber> fibers;

    // The main thread also registers with the scheduler
    boost::fibers::use_scheduling_algorithm<boost::fibers::algo::work_stealing>(thread_count);

    YAMLProcessor yaml_processor{};

    // Launch all the fibers from the main thread
    fibers.emplace_back(&YAMLProcessor::parseGeneral, &yaml_processor);
    fibers.emplace_back(&YAMLProcessor::parseExperience, &yaml_processor);
    fibers.emplace_back(&YAMLProcessor::parseProject, &yaml_processor);

        // 5. Join all the fibers
    // While waiting, the main thread will also participate in running fibers.
    for (auto& f : fibers) {
        f.join();
    }

    // 6. Signal the barrier
    // Now that all fibers are done, this will unblock all the waiting worker threads.
    std::cout << "\nAll fibers have completed. Signaling threads to exit.\n";
    b.wait();

    // 7. Join all the worker threads to ensure they exit cleanly
    for (auto& t : threads) {
        t.join();
    }

    std::cout << "All threads have finished." << std::endl;


   std::cout << "\n--- All YAML files parsed successfully ---\n\n";
  // yaml_processor.general.print_skills_achivements();
  return 0;
}
