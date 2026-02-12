#import "../template/template.typ": *
#import "./config/personal_info.typ": *
#import "./cv_type/section.typ"
#import "./cv_type/cv_type.typ"

#show: resume.with(
  author: name,
  location: location,
  email: email,
  github: github,
  linkedin: linkedin,
  phone: phone,
  personal-site: personal-site,
  accent-color: "#26428b",
  font: "New Computer Modern",
  paper: "us-letter",
  author-position: center,
  personal-info-position: center,
)

#for sec in cv_type.cv-default.layout{
  [#sec]
}
