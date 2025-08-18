package org.senop;

import org.senop.parse.YAMLProcessor;

public class Main {
    public static void main(String[] args) {
        YAMLProcessor yamlProcessor = new YAMLProcessor();
        yamlProcessor.Parse();
        yamlProcessor.Print();
    }
}
