package org.senop.model;

import java.util.List;
import com.fasterxml.jackson.annotation.JsonProperty;

public final class Project {
  private String name;
  private String github;
  @JsonProperty("github_handle")
  private String githubHandle;
  @JsonProperty("cv_type")
  private List<String> cvType;
  private String description;
  private List<String> points;

  public String getName() {
    return name;
  }
  public void setName(String name) {
    this.name = name;
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
  public List<String> getCvType() {
    return cvType;
  }
  public void setCvType(List<String> cvType) {
    this.cvType = cvType;
  }
  public String getDescription() {
    return description;
  }
  public void setDescription(String description) {
    this.description = description;
  }
  public List<String> getPoints() {
    return points;
  }
  public void setPoints(List<String> points) {
    this.points = points;
  }
}
