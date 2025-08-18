package org.senop.model;

import java.util.List;
import com.fasterxml.jackson.annotation.JsonProperty;

public final class Experience {
  private String company;
  private String location;
  private String role;
  private String dates;
  @JsonProperty("job_type")
  private String jobType;
  @JsonProperty("cv_type")
  private List<String> cvType;
  private List<String> points;

  public String getCompany() {
    return company;
  }
  public void setCompany(String company) {
    this.company = company;
  }

  public String getLocation() {
    return location;
  }
  public void setLocation(String location) {
    this.location = location;
  }

  public String getRole() {
    return role;
  }
  public void setRole(String role) {
    this.role = role;
  }
  public String getDates() {
    return dates;
  }
  public void setDates(String dates) {
    this.dates = dates;
  }
  public String getJobType() {
    return jobType;
  }
  public void setJobType(String jobType) {
    this.jobType = jobType;
  }
  public List<String> getCvType() {
    return cvType;
  }
  public void setCvType(List<String> cvType) {
    this.cvType = cvType;
  }
  public List<String> getPoints() {
    return points;
  }
  public void setPoints(List<String> points) {
    this.points = points;
  }
}
