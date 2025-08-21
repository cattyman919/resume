set(SOURCES
  "main.cpp"
  "model/project.cpp"
  "model/experience.cpp"
  "model/general.cpp"
  "parse/YAMLProcessor.cpp"
)
add_executable(${PROJECT_NAME} "${SOURCES}")

add_executable(test "test_fiber.cpp")

# target_compile_definitions(${PROJECT_NAME}
#   PUBLIC
#   DEBUG
# )

target_include_directories(${PROJECT_NAME}
  PUBLIC
  "${CMAKE_SOURCE_DIR}/include"
)

target_include_directories(test
  PUBLIC
  "${CMAKE_SOURCE_DIR}/include"
)
