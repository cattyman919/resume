#pragma once

#include <concepts>
#include <cstddef>
#include <thread>
#include <utility>

#include "boost/fiber/all.hpp"

class BoostFiberScheduler {
 public:
  explicit BoostFiberScheduler(size_t thread_count);

  // Remove Copy Constructor and Assignment
  BoostFiberScheduler(const BoostFiberScheduler& other) = delete;
  BoostFiberScheduler& operator=(const BoostFiberScheduler& other) = delete;

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

  size_t thread_count_{};
  boost::fibers::barrier b_;
  std::vector<std::thread> threads{};
  std::vector<boost::fibers::fiber> fibers{};
  bool is_shutdown_{false};
};
