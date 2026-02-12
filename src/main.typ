#import "../template/template.typ": *
#import "./config/personal_info.typ": *
#import "./cv_type/section.typ"
#import "./cv_type/cv_type.typ"

// 1. Get the target CV type from CLI inputs (defaults to "default")
#let target-type = sys.inputs.at("type", default: "default")

// 2. Fetch the specific configuration from the map we created in Step 1
#let selected-cv = cv_type.cv-map.at(
  target-type,
  default: cv_type.default
)


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

// 3. Render the layout defined in the selected CV
#for sec in selected-cv.layout {
  [#sec]
}
