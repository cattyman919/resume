#include <iostream>
#include "yaml-cpp/yaml.h"

struct Employee
{
    int id {};
    int age {};
    double wage {};
};


std::ostream& operator<<(std::ostream& os, const Employee& e) {
    os << "Employee { id: " << e.id
       << ", age: " << e.age
       << ", wage: " << e.wage << " }";
    return os;
}

int main(){
    try {
        // Load the YAML file into a root node.
        // This line does all the parsing.
        YAML::Node config = YAML::LoadFile("config.yaml");

        // --- Accessing Map Values ---
        // You can chain accessors for nested maps.
        std::string appName = config["application"]["name"].as<std::string>();
        double appVersion = config["application"]["version"].as<double>();
        int dbPort = config["database"]["port"].as<int>();

        std::cout << "Application Name: " << appName << std::endl;
        std::cout << "Application Version: " << appVersion << std::endl;
        std::cout << "Database Port: " << dbPort << std::endl;

        // --- Checking for Optional Keys ---
        // Always check if a node exists before accessing it.
        if (config["application"]["threaded"]) {
            bool isThreaded = config["application"]["threaded"].as<bool>();
            std::cout << "Threading enabled: " << (isThreaded ? "true" : "false") << std::endl;
        }

        // --- Iterating Over a Sequence (List) ---
        const YAML::Node& users = config["database"]["users"];
        if (users && users.IsSequence()) {
            std::cout << "Database Users:" << std::endl;
            for (const YAML::Node& user : users) {
                std::cout << "- " << user.as<std::string>() << std::endl;
            }
        }

    } catch (const YAML::Exception& e) {
        // Handle file not found, parsing errors, etc.
        std::cerr << "Error parsing YAML: " << e.what() << std::endl;
        return 1;
    }
  return 0;
}
