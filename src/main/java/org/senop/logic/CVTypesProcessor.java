package org.senop.logic;

import java.util.List;
import java.util.Set;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.Executors;

import org.senop.model.Experience;
import org.senop.model.Project;

public class CVTypesProcessor {

    private Set<String> cv_types;

    public Set<String> getCv_types() {
      return cv_types;
    }

    public CVTypesProcessor() {
        this.cv_types = ConcurrentHashMap.newKeySet();
    }

    public void ProcessCVTypes(List<Experience> experiences, List<Project> projects ){

        try (var executor = Executors.newVirtualThreadPerTaskExecutor()) {
          
          experiences.stream().forEach(experience -> {
            executor.submit(() -> {
                cv_types.addAll(experience.getCvType());
            });
          });

          projects.stream().forEach(project -> {
            executor.submit(() -> {
                cv_types.addAll(project.getCvType());
            });
          });
        } 
    }

    public void Print() {
        System.out.println("CV Types:");
        cv_types.forEach(System.out::println);
    }
}
