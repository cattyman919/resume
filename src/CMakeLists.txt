set(SOURCES "main.cpp")
add_executable(app "${SOURCES}")

target_include_directories(app
  PUBLIC
  "${CMAKE_SOURCE_DIR}/include"
)
