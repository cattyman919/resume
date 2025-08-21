
find_package(Boost REQUIRED COMPONENTS fiber context thread)

if(Boost_FOUND)
    message(STATUS "Found Boost ${Boost_VERSION}")
endif()

find_package(yaml-cpp 0.8 REQUIRED)
