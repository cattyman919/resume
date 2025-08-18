package org.senop.model;

import com.fasterxml.jackson.annotation.JsonProperty;
import java.util.List;

public final class General {
  @JsonProperty("personal_info")
  public PersonalInfo personalInfo;

  @JsonProperty("skills_achievements")
  private SkillsAchievements skillsAchievements;

  @JsonProperty("educations")
  private List<Education> educations;

  @JsonProperty("awards")
  private List<Award> awards;

  public PersonalInfo getPersonalInfo() {
    return personalInfo;
  }

  public void setPersonalInfo(PersonalInfo personalInfo) {
    this.personalInfo = personalInfo;
  }

  public SkillsAchievements getSkillsAchievements() {
    return skillsAchievements;
  }

  public void setSkillsAchievements(SkillsAchievements skillsAchievements) {
    this.skillsAchievements = skillsAchievements;
  }

  public List<Education> getEducations() {
    return educations;
  }

  public void setEducations(List<Education> educations) {
    this.educations = educations;
  }

  public List<Award> getAwards() {
    return awards;
  }

  public void setAwards(List<Award> awards) {
    this.awards = awards;
  }

  public static final class PersonalInfo {
    private String name;
    private String location;
    private String email;
    private String phone;
    private String website;
    private String linkedin;

    @JsonProperty("linkedin_handle")
    private String linkedinHandle;

    private String github;

    @JsonProperty("github_handle")
    private String githubHandle;

    @JsonProperty("profile_pic")
    private String profilePic;

    public String getName() {
      return name;
    }

    public void setName(String name) {
      this.name = name;
    }

    public String getLocation() {
      return location;
    }

    public void setLocation(String location) {
      this.location = location;
    }

    public String getEmail() {
      return email;
    }

    public void setEmail(String email) {
      this.email = email;
    }

    public String getPhone() {
      return phone;
    }

    public void setPhone(String phone) {
      this.phone = phone;
    }

    public String getWebsite() {
      return website;
    }

    public void setWebsite(String website) {
      this.website = website;
    }

    public String getLinkedin() {
      return linkedin;
    }

    public void setLinkedin(String linkedin) {
      this.linkedin = linkedin;
    }

    public String getLinkedinHandle() {
      return linkedinHandle;
    }

    public void setLinkedinHandle(String linkedinHandle) {
      this.linkedinHandle = linkedinHandle;
    }

    public String getGithub() {
      return github;
    }

    public void setGithub(String github) {
      this.github = github;
    }

    public String getGithubHandle() {
      return githubHandle;
    }

    public void setGithubHandle(String githubHandle) {
      this.githubHandle = githubHandle;
    }

    public String getProfilePic() {
      return profilePic;
    }

    public void setProfilePic(String profilePic) {
      this.profilePic = profilePic;
    }
  }

  public static final class SkillsAchievements {
    @JsonProperty("Hard Skills")
    private List<String> hardSkill;

    @JsonProperty("Soft Skills")
    private List<String> softSkill;

    @JsonProperty("Programming Languages")
    private List<String> programmingLanguages;

    @JsonProperty("Database Languages")
    private List<String> databaseLanguages;

    @JsonProperty("Misc")
    private List<String> misc;

    @JsonProperty("Certificates")
    private List<Certificate> certificates;

    public List<String> getHardSkill() {
      return hardSkill;
    }

    public void setHardSkill(List<String> hardSkill) {
      this.hardSkill = hardSkill;
    }

    public List<String> getSoftSkill() {
      return softSkill;
    }

    public void setSoftSkill(List<String> softSkill) {
      this.softSkill = softSkill;
    }

    public List<String> getProgrammingLanguages() {
      return programmingLanguages;
    }

    public void setProgrammingLanguages(List<String> programmingLanguages) {
      this.programmingLanguages = programmingLanguages;
    }

    public List<String> getDatabaseLanguages() {
      return databaseLanguages;
    }

    public void setDatabaseLanguages(List<String> databaseLanguages) {
      this.databaseLanguages = databaseLanguages;
    }

    public List<String> getMisc() {
      return misc;
    }

    public void setMisc(List<String> misc) {
      this.misc = misc;
    }

    public List<Certificate> getCertificates() {
      return certificates;
    }

    public void setCertificates(List<Certificate> certificates) {
      this.certificates = certificates;
    }
  }

  public static final class Certificate {
    private String name;
    private String year;

    public String getName() {
      return name;
    }

    public void setName(String name) {
      this.name = name;
    }

    public String getYear() {
      return year;
    }

    public void setYear(String year) {
      this.year = year;
    }
  }

  public static final class Education {
    private String institution;
    private String degree;
    private String dates;
    private String gpa;
    private List<String> details;

    public String getInstitution() {
      return institution;
    }

    public void setInstitution(String institution) {
      this.institution = institution;
    }

    public String getDates() {
      return dates;
    }

    public void setDates(String dates) {
      this.dates = dates;
    }

    public String getDegree() {
      return degree;
    }

    public void setDegree(String degree) {
      this.degree = degree;
    }

    public String getGpa() {
      return gpa;
    }

    public void setGpa(String gpa) {
      this.gpa = gpa;
    }

    public List<String> getDetails() {
      return details;
    }

    public void setDetails(List<String> details) {
      this.details = details;
    }
  }

  public static final class Award {
    private String title;
    private String organization;
    private String date;
    private List<String> points;

    public String getTitle() {
      return title;
    }

    public void setTitle(String title) {
      this.title = title;
    }

    public String getOrganization() {
      return organization;
    }

    public void setOrganization(String organization) {
      this.organization = organization;
    }

    public String getDate() {
      return date;
    }

    public void setDate(String date) {
      this.date = date;
    }

    public List<String> getPoints() {
      return points;
    }

    public void setPoints(List<String> points) {
      this.points = points;
    }
  }
}
