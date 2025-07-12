# Tools
ENGINE = pdflatex
RM = rm -fr

# List of CV types. This is the only line you need to edit to add a new CV!
TYPES := main fullstack devops

# Generate the target names from the TYPES list (e.g., main-cv, fullstack-cv)
CV_TARGETS := $(TYPES:%=%-cv)

# A helper function to capitalize the first letter of a word
capitalize = $(shell echo $(1) | sed 's/./\u&/')

# Compile Command
define COMPILE_LATEX
	$(ENGINE) -output-directory=out -output-format=pdf -jobname "Seno Pamungkas Rahman - CV ($(call capitalize, $(1)))" main_cv/main.tex
	$(ENGINE) -output-directory=out -output-format=pdf -jobname "Seno Pamungkas Rahman - CV ($(call capitalize, $(1))) (BW)" bw_cv/main.tex
endef

.PHONY: all clean $(CV_TARGETS)

all: $(CV_TARGETS)

$(CV_TARGETS): %-cv:
	@echo "Generating and compiling $(call capitalize, $*) CV..."
	go run generate_sections.go --type=$*
	$(call COMPILE_LATEX, $*)

clean:
	$(RM) -rf out/*
