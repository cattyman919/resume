package org.senop.parse;

import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.dataformat.yaml.YAMLFactory;
import java.nio.file.Path;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.Executors;
import java.util.function.Consumer;
import java.util.function.Supplier;
import org.senop.model.Experience;
import org.senop.model.General;
import org.senop.model.Project;

public final class YAMLProcessor {

  private General general;
  private List<Experience> experiences;
  private List<Project> projects;

  public YAMLProcessor() {
    this.general = new General();
    this.experiences = new ArrayList<>();
    this.projects = new ArrayList<>();
  }

  private record ParsingTask<T>(String filename, Supplier<T> parser, Consumer<T> resultConsumer) {}

  public void Parse() {
    System.out.println("Parsing YAML files...");

    ObjectMapper mapper = new ObjectMapper(new YAMLFactory());

    List<ParsingTask<?>> tasks =
        List.of(
            new ParsingTask<>(
                "general.yaml",
                () -> {
                  try {
                    return mapper.readValue(
                        Path.of("config", "general.yaml").toFile(), General.class);
                  } catch (Exception e) { // Catch the checked exception
                    throw new RuntimeException(e); // And wrap it
                  }
                },
                result -> this.general = result),
            new ParsingTask<>(
                "projects.yaml",
                () -> {
                  try {
                    return mapper.readValue(
                        Path.of("config", "projects.yaml").toFile(),
                        new TypeReference<List<Project>>() {});
                  } catch (Exception e) {
                    throw new RuntimeException(e);
                  }
                },
                result -> this.projects = result),
            new ParsingTask<>(
                "experiences.yaml",
                () -> {
                  try {
                    return mapper.readValue(
                        Path.of("config", "experiences.yaml").toFile(),
                        new TypeReference<List<Experience>>() {});
                  } catch (Exception e) {
                    throw new RuntimeException(e);
                  }
                },
                result -> this.experiences = result));

    try (var executor = Executors.newVirtualThreadPerTaskExecutor()) {

      List<CompletableFuture<Void>> futures =
          tasks.stream()
              .map(
                  task ->
                      CompletableFuture.supplyAsync(
                              () -> {
                                try {
                                  return task.parser();
                                } catch (Exception e) {
                                  throw new RuntimeException(e);
                                }
                              },
                              executor)
                          .thenAccept(result -> task.resultConsumer())
                          .exceptionally(
                              ex -> {
                                System.err.printf(
                                    "Error parsing %s: %s\n",
                                    task.filename(), ex.getCause().getMessage());
                                return null;
                              }))
              .toList();

      CompletableFuture.allOf(futures.toArray(new CompletableFuture[0])).join();
    }
  }

  public void Print() {
    System.out.println("General Information:");
    System.out.println("Name: " + general.getPersonalInfo().getName());
    System.out.println("Location: " + general.getPersonalInfo().getLocation());
    System.out.println("Email: " + general.getPersonalInfo().getEmail());
    System.out.println("Phone: " + general.getPersonalInfo().getPhone());
    System.out.println("Website: " + general.getPersonalInfo().getWebsite());
    System.out.println("LinkedIn: " + general.getPersonalInfo().getLinkedin());
    System.out.println("GitHub: " + general.getPersonalInfo().getGithub());

    System.out.println("\nExperiences:");
    experiences.forEach(
        exp -> {
          System.out.println("Company: " + exp.getCompany());
          System.out.println("Role: " + exp.getRole());
          System.out.println("Dates: " + exp.getDates());
          System.out.println("Job Type: " + exp.getJobType());
          System.out.println("Points: " + String.join(", ", exp.getPoints()));
          System.out.println();
        });

    System.out.println("\nProjects:");
    projects.forEach(
        proj -> {
          System.out.println("Project Name: " + proj.getName());
          System.out.println("GitHub: " + proj.getGithub());
          System.out.println("Description: " + proj.getDescription());
          System.out.println("Points: " + String.join(", ", proj.getPoints()));
          System.out.println();
        });
  }
}
