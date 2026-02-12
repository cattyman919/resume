#import "../../template/template.typ"

#import "../config/experiences.typ"
#import "../config/projects.typ"
#import "section.typ"

// TODO
// Default Personal info
// We might need to make this more flexible and configurable later for every CV type
#let cv(
  layout: () ,
  experiences: (),
  projects: (),
) = {

  let set-value(sec) = {
    if sec == section.experiences{
      return section.experiences(experiences: experiences)
      // return experiences
    } else if sec == section.projects{
      return section.projects(projects:projects)
    } else{
      return sec
    }
  }

  layout = layout.map(set-value)

  return (
    headers: section.headers,
    layout: layout,
  )
}

#let cv-default = cv(
  layout: (
    // section.research-interests,
    section.experiences,
    section.educations,
    section.projects,
    section.skills,
    section.certificates,
    section.awards
  ),
  experiences: (
      experiences.xlsmart,
      experiences.superbank,
      experiences.xlaxiata,
      experiences.bank-victoria
  ),
  projects: (
    projects.vxlang,
    projects.yazi
  )
)

#let cv-frontend = cv(
  layout: (
    section.experiences,
    section.educations,
    section.projects,
    section.skills,
    section.awards,
  ),
  experiences: (
      experiences.superbank,
      experiences.xlaxiata,
      experiences.bank-victoria,
      experiences.mileapp,
  ),
  projects: (
    projects.restomatic,
    projects.jaga,
    projects.yazi
  )
)
