#import "../config/personal_info.typ"
#import "../config/experiences.typ"
#import "../config/projects.typ"
#import "../../template/template.typ"

// Generic two by two component for resume
#let generic-two-by-two(
  top-left: "",
  top-right: "",
  bottom-left: "",
  bottom-right: "",
) = {
  [
    #top-left #h(1fr) #top-right \
    #bottom-left #h(1fr) #bottom-right
  ]
}



#let headers = (
 personal_info.name,
 personal_info.location,
 personal_info.email,
 personal_info.phone,
 personal_info.github,
 personal_info.linkedin,
 personal_info.personal-site,
)

#let projects(
  projects: ()
) = [
  == Projects
  #for proj in projects [
    #template.project(
      title: proj.title,
      description: proj.description,
      url: proj.url,
      url-handle: proj.url_handle,
      points: proj.points
    )
  ]
]

#let experiences(
  experiences: ()
) = [
  == Experience
  #for exp in experiences [
    #template.experience(
      company: exp.company,
      dates: exp.dates,
      role: exp.role,
      job-type: exp.job-type,
      location: exp.location,
      points: exp.points
    )
  ]
]

#let certificates = [
== Certificates
  #for cert in personal_info.certificates{
    [- (*#cert.year*) #cert.name ]
  }
]

#let awards = [
== Awards
#for award in personal_info.awards [
  - #template.award(
    title: award.title,
    organization: award.organization,
    date: award.date,
    points: award.points
  )
]
]

#let research-interests = [
  == Research Interests
  #for interest in personal_info.research_interests [
    - #interest
  ]
]

#let skills = [
  == Skills
  #for skill in personal_info.skills_achievements [
    - #template.skill(
      title: skill.title,
      items: skill.items
    )
  ]
]

// #let educations = personal_info.educations
#let educations = [
  == Educations
  #for edu in personal_info.educations [
    #template.education(
      institution: edu.institution,
      degree: edu.degree,
      details: edu.details,
      dates:edu.dates,
    )
  ]
]

#let projects-mark = "SECTION_PROJECTS"
#let experiences-mark = "SECTION_EXPERIENCES"
