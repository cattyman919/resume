package org.senop.parse;

import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.dataformat.yaml.YAMLFactory;
import java.nio.file.Path;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.Callable;
import java.util.concurrent.Executors;
import java.util.concurrent.Future;
import java.util.function.Consumer;
import org.senop.model.Experience;
import org.senop.model.General;
import org.senop.model.Project;

public final class YAMLProcessor {

  private General general;
  private List<Experience> experiences;
  private List<Project> projects;

  public General getGeneral() {
    return general;
  }

  public void setGeneral(General general) {
    this.general = general;
  }

  public List<Experience> getExperiences() {
    return experiences;
  }

  public void setExperiences(List<Experience> experiences) {
    this.experiences = experiences;
  }

  public List<Project> getProjects() {
    return projects;
  }

  public void setProjects(List<Project> projects) {
    this.projects = projects;
  }

  public YAMLProcessor() {
    this.general = new General();
    this.experiences = new ArrayList<>();
    this.projects = new ArrayList<>();
  }

   private <T> Callable<Void> createParsingTask(ObjectMapper mapper, String filename, Class<T> type, Consumer<T> resultConsumer) {
        return () -> {
            try {
                var res = mapper.readValue(Path.of("config", filename).toFile(), type);
                resultConsumer.accept(res);
            } catch (Exception e) {
                throw new YAMLParsingException("Failed to parse file: " + filename, e);
            }
            return null; 
        };
    }

    private <T> Callable<Void> createParsingTask(ObjectMapper mapper, String filename, TypeReference<T> typeRef, Consumer<T> resultConsumer) {
        return () -> {
            try {
                var res = mapper.readValue(Path.of("config", filename).toFile(), typeRef);
                resultConsumer.accept(res);
            } catch (Exception e) {
                throw new YAMLParsingException("Failed to parse file: " + filename, e);
            }
            return null; 
        };
    }

  public void Parse() {
    System.out.println("Parsing YAML files...");

    ObjectMapper mapper = new ObjectMapper(new YAMLFactory());

        List<Callable<Void>> tasks = List.of(
            createParsingTask(mapper, "general.yaml", General.class, res -> this.general = res),
            createParsingTask(mapper, "projects.yaml", new TypeReference<List<Project>>() {}, res -> this.projects = res),
            createParsingTask(mapper, "experiences.yaml", new TypeReference<List<Experience>>() {}, res -> this.experiences = res)
        );

        try (var executor = Executors.newVirtualThreadPerTaskExecutor()) {
            List<Future<Void>> futures = executor.invokeAll(tasks);

            for (var future : futures) {
                try {
                    future.get(); 
                } catch (Exception e) {
                    System.err.println("[ERROR] A parsing task failed catastrophically.");
                    if (e.getCause() instanceof YAMLParsingException) {
                        System.err.println(e.getCause().getMessage());
                        System.err.println(e.getCause().getCause().getClass().getSimpleName() + " - " + e.getCause().getCause().getMessage());
                    } else {
                        e.printStackTrace(); 
                    }
                    System.exit(1);
                }
            }
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new RuntimeException("Main thread was interrupted during parsing.", e);
        }
    // Automaically join all tasks when the try-with-resources block exits
    // because it implements AutoCloseable

    System.out.println("All YAML files parsed successfully.");
  }

  public void Print() {
    System.out.println("\nGeneral Information:");
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
