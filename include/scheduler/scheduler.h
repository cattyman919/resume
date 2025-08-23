#pragma once

#include <concepts>
#include <thread>
#include <utility>

#include "boost/fiber/all.hpp"

class BoostFiberScheduler {
 public:
  BoostFiberScheduler(unsigned int thread_count);
  ~BoostFiberScheduler();

  // A user-friendly way to add a task (a fiber).
  // Takes any function and its arguments.
  template <typename F, typename... Args>
    requires std::invocable<F, Args...>
  void addTask(F&& func, Args&&... args) {
    // Debug Version
    fibers.emplace_back([&] {
      std::stringstream ss;
      ss << "Running Task in thread (" << std::this_thread::get_id() << ")\n";
      std::cout << ss.str();
      std::invoke(std::forward<F>(func), std::forward<Args>(args)...);
    });
    // Release Version
    // fibers.emplace_back(std::forward<F>(func), std::forward<Args>(args)...);
  };

  // Waits for all fibers to complete
  void join();

  size_t getTotalTask() { return fibers.size(); }
  size_t getFiberCapacity() { return fibers.capacity(); }

 private:
  void setupWorkerThreads();

  // Waits for all fibers to complete and then cleanly shuts down the worker
  void shutdown();

  unsigned int thread_count_;
  boost::fibers::barrier b_;
  std::vector<std::thread> threads;
  std::vector<boost::fibers::fiber> fibers;
  bool is_shutdown_ = false;
};
