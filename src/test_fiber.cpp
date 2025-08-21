#include <iostream>
#include <vector>
#include <thread>
#include <chrono>
#include <sstream>

#include <boost/fiber/all.hpp>
// Make sure to include the barrier header

void fiber_function(int fiber_id) {
    std::stringstream ss;
    ss << "Fiber " << fiber_id << " is STARTING on OS thread " << std::this_thread::get_id() << std::endl;
    std::cout << ss.str();

    // Use fiber-aware sleep to simulate work
    boost::this_fiber::sleep_for(std::chrono::milliseconds(20));

    ss.str(""); // Clear the stream
    ss << "Fiber " << fiber_id << " is FINISHING on OS thread " << std::this_thread::get_id() << std::endl;
    std::cout << ss.str();
}

int main() {
    unsigned int thread_count = std::thread::hardware_concurrency();
    std::cout << "Main thread " << std::this_thread::get_id() << " will use a total of " << thread_count << " threads." << std::endl;

    // 1. Create a barrier that will wait for all threads (main + workers)
    boost::fibers::barrier b(thread_count);

    std::vector<std::thread> threads;
    std::vector<boost::fibers::fiber> fibers;

    // 2. Launch worker threads
    // Start from 1 because the main thread is thread 0
    for (unsigned int i = 1; i < thread_count; ++i) {
        threads.emplace_back([&b, thread_count]() {
            // Each thread must register with the scheduler
            boost::fibers::use_scheduling_algorithm<boost::fibers::algo::work_stealing>(thread_count);

            // *** THE KEY FIX ***
            // Wait on the barrier. This blocks the thread's main fiber,
            // making the thread available to the scheduler for running other fibers (work-stealing).
            // It will only unblock when the main thread also calls wait().
            b.wait();
        });
    }

    // 3. The main thread also registers with the scheduler
    boost::fibers::use_scheduling_algorithm<boost::fibers::algo::work_stealing>(thread_count);

    // 4. Launch all the fibers from the main thread
    for (int i = 0; i < 40; ++i) {
        fibers.emplace_back(fiber_function, i);
    }

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

    return 0;
}
