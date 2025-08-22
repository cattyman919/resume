#pragma once

#include <concepts>
#include <thread>

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
    fibers.emplace_back(std::forward<F>(func), std::forward<Args>(args)...);
  };

  // Waits for all fibers to complete
  void join();

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
