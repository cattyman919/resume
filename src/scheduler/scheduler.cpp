#include "scheduler/scheduler.h"

#include <iostream>
#include <thread>

#include "boost/fiber/all.hpp"

BoostFiberScheduler::BoostFiberScheduler(unsigned int thread_count)
    : thread_count_(thread_count),
      b_(thread_count)  // Initialize barrier to wait for all threads
{
  // Register Woker threads with the scheduler first, then main thread.
  setupWorkerThreads();

  // Register the main thread with the work-stealing scheduler.
  boost::fibers::use_scheduling_algorithm<boost::fibers::algo::work_stealing>(
      thread_count_);
}

BoostFiberScheduler::~BoostFiberScheduler() { shutdown(); }

// Waits for all fibers to complete
// threads.
// @throws YAML::Exception if any fiber throws an exception.
// @throws std::exception for any other unexpected errors.
void BoostFiberScheduler::join() {
  // Join all the fibers. The main thread will participate in work.
  for (auto &f : fibers) {
    if (f.joinable()) {
      f.join();
    }
  }

  // Signal the barrier to release the worker threads.
  std::cout << "\nAll fibers have completed.\n";
}

// Shuts down the scheduler and waits for all worker threads to finish.
void BoostFiberScheduler::shutdown() {
  // Only run this logic once.
  if (is_shutdown_) {
    return;
  }

  this->join();  // waits for the fibers to finish
  b_.wait();     // Wait for all worker threads to finish their tasks.

  // Join all the worker threads.
  for (auto &t : threads) {
    if (t.joinable()) {
      t.join();
    }
  }

  is_shutdown_ = true;
  std::cout << "All threads have finished." << std::endl;
}

void BoostFiberScheduler::setupWorkerThreads() {
  for (unsigned int i = 1; i < thread_count_; ++i) {
    threads.emplace_back([this]() {
      boost::fibers::use_scheduling_algorithm<
          boost::fibers::algo::work_stealing>(thread_count_);

      b_.wait();
    });
  }
}
