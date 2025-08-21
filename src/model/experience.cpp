#include "model/experience.h"
#include <iostream>

std::ostream& operator<<(std::ostream& os, const Experience& e) {
   os << "Experience {\nname: " << e.company << '\n'
       << "location: " << e.location << '\n'
       << "role: " << e.role << '\n'
       << "dates: " << e.dates << '\n'
       << "job_type: " << e.job_type << '\n'
       << "cv_type: [ " ;
        for (const auto& type : e.cv_type){
          os << type << ", ";
        }
       os << " ]\n"<< "points: [\n" ;
        for (const auto& point : e.points){
          os  << "- "<< point << "\n";
        }
       os << "]"  << " }";
    return os;
}
