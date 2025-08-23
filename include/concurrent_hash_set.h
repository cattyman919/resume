#pragma once

#include <boost/fiber/mutex.hpp>
#include <unordered_set>

/**
 * @class ConcurrentHashSet
 * @brief A thread-safe (or fiber-safe) hash set using a mutex.
 *
 * This class wraps a std::unordered_set and a boost::fibers::mutex to ensure
 * that all write operations are synchronized, preventing race conditions when
 * accessed by multiple fibers concurrently.
 */
template <typename T>

class ConcurrentHashSet final {
 public:
  ConcurrentHashSet() = default;

  /**
   * @brief Inserts an element into the hash set in a fiber-safe manner.
   * @param value The value to insert.
   */
  void insert(const T& value) {
    // std::lock_guard is a RAII-style mutex wrapper.
    // It automatically locks the mutex when created and unlocks it
    // when it goes out of scope (even if an exception is thrown).
    std::lock_guard<boost::fibers::mutex> lock(mtx_);

    // This is the critical section. Only one fiber can execute this code at a
    // time.
    set_.insert(value);
  }

  /**
   * @brief Returns the number of elements in the set.
   *
   * Note: For truly safe concurrent access, this read operation should
   * also be protected by the mutex if you have concurrent readers and writers.
   * Since the prompt only specified concurrent writers, we'll protect it anyway
   * as a best practice to get an accurate size after all writes are done.
   * @return The number of elements.
   */
  size_t size() {
    std::lock_guard<boost::fibers::mutex> lock(mtx_);
    return set_.size();
  }

  const std::unordered_set<T>& getSet() const { return set_; }

 private:
  std::unordered_set<T> set_;
  boost::fibers::mutex mtx_;
};
