#pragma once

#include "boost/fiber/all.hpp"
#include <thread>

class BoostFiberScheduler {
  public:
    BoostFiberScheduler(unsigned int thread_count);
    ~BoostFiberScheduler();

    // A user-friendly way to add a task (a fiber).
    // Takes any function and its arguments.
    template<typename F, typename... Args>
    void addTask(F&& func, Args&&... args){
        fibers.emplace_back(std::forward<F>(func), std::forward<Args>(args)...);
    };

    // Waits for all fibers to complete and then cleanly shuts down the worker threads.
    void wait();

    private:
    void setupWorkerThreads();

    unsigned int thread_count_;
    boost::fibers::barrier b_;
    std::vector<std::thread> threads;
    std::vector<boost::fibers::fiber> fibers;
    bool is_shutdown_ = false;
};
