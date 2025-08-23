#include <boost/fiber/all.hpp>
#include <chrono>
#include <iostream>
#include <sstream>
#include <thread>
#include <vector>
// Make sure to include the barrier header

void fiber_function() {
  std::stringstream ss;
  ss << "Another thread " << std::this_thread::get_id() << std::endl;
  std::cout << ss.str();

  // Use fiber-aware sleep to simulate work
  boost::this_fiber::sleep_for(std::chrono::milliseconds(1000));

  ss.str("");  // Clear the stream
  ss << "Another thread " << std::this_thread::get_id() << std::endl;
  std::cout << ss.str();
}

boost::fibers::fiber f1;  // not-a-fiber
int main() {
  // std::cout << "Hello world\n";
  std::cout << "BEFORE Main Thread: " << std::this_thread::get_id() << '\n';

  std::thread oof{[&] {
    boost::fibers::use_scheduling_algorithm<boost::fibers::algo::work_stealing>(
        2);

    boost::fibers::fiber bruh(fiber_function);

    bruh.join();
  }};

  oof.join();

  std::cout << "AFTER Main Thread: " << std::this_thread::get_id() << '\n';
  //
  // f1 = std::move(f2);  // f2 moved to f1
  return 0;
}
