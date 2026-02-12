#import "../../template/template.typ"

#import "../config/experiences.typ"
#import "../config/projects.typ"
#import "section.typ"

// Helper: Select specific points from an experience
#let pick(exp, indices) = {
  let selected-points = ()

  // Handle cases where points might be a dictionary (legacy support)
  let source-points = if type(exp.points) == dictionary {
    exp.points.at("default", default: ())
  } else {
    exp.points
  }

  // Loop through the requested indices and grab them
  for i in indices {
    // Safety check to avoid crashing if index is out of bounds
    if i < source-points.len() {
      selected-points.push(source-points.at(i))
    }
  }

  // Return a new dictionary merging the old experience with the NEW points
  return exp + (points: selected-points)
}

#let cv(
  // Default Layout
  layout: (
    section.experiences,
    section.educations,
    section.projects,
    section.skills,
    section.certificates,
    section.awards,
  ),
  experiences: (),
  projects: (),
) = {

  let set-value(sec) = {
    if sec == section.experiences{
      return section.experiences(experiences: experiences)
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


#let default = cv(
  experiences: (
      pick(experiences.xlsmart, range(0,5) + (6, 8, 10, 12, 15)),
      experiences.superbank,
      experiences.xlaxiata,
      experiences.bank-victoria
  ),
  projects: (
    projects.vxlang,
    projects.yazi
  )
)

#let frontend = cv(
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

#let cv-map = (
  "default": default,
  "frontend": frontend,
  // Add more types here as you define them, e.g.:
  // "devops": cv-devops,
)
