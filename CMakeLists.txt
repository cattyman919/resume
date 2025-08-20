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

add_subdirectory(src)
