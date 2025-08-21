cmake_minimum_required(VERSION 3.15...4.0)

include(cmake/compiler_options.cmake)

set(CMAKE_EXPORT_COMPILE_COMMANDS ON)
set(CMAKE_RUNTIME_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/bin)

project(
    Resume
    VERSION 0.1
    DESCRIPTION "Resume Project Builder"
    LANGUAGES CXX
)

include(cmake/dependencies.cmake)

add_subdirectory(src)

target_include_directories(${PROJECT_NAME}
  PRIVATE
  ${Boost_INCLUDE_DIRS}
)

target_link_libraries(
  ${PROJECT_NAME} PRIVATE
  yaml-cpp::yaml-cpp
  Boost::fiber
  Boost::thread
)

target_link_libraries(
  test PRIVATE
  Boost::fiber
  Boost::thread
)
