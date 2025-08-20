package org.senop;

import org.senop.logic.CVTypesProcessor;
import org.senop.parse.YAMLProcessor;

public class Main {
    public static void main(String[] args) {
        YAMLProcessor yamlProcessor = new YAMLProcessor();
        yamlProcessor.Parse();
        yamlProcessor.Print();

        CVTypesProcessor cvTypesProcessor = new CVTypesProcessor();
        cvTypesProcessor.ProcessCVTypes(
            yamlProcessor.getExperiences(),
            yamlProcessor.getProjects()
        );
        cvTypesProcessor.Print();
    }
}
