#include "model/general.h"
#include <iostream>
#include <type_traits>

void General::print_skills_achivements() const {
  std::cout << "--- Achievements ---\n";
  for (const auto &item : this->skills_achievements) {
    // Use std::visit to safely handle the different types
    std::visit(
        [](const auto &value) {
          // 'if constexpr' checks the type at compile time
          if constexpr (std::is_same_v<decltype(value), const Skills &>) {
            std::cout << value.name << ":\n";
            for (const auto &skill : value.skills) {
              std::cout << "  - " << skill << "\n";
            }
          } else if constexpr (std::is_same_v<decltype(value),
                                              const Certificate &>) {
            std::cout << "Certificate:\n"; // Or you could use a different title
            std::cout << "  - " << value.name << " (" << value.year << ")\n";
          }
        },
        item);
  }
}
