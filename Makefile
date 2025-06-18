# Tools
ENGINE = pdflatex
RM = rm -fr
PYTHON = python3

# Output file names (variables for easier changes)
MAIN_CV_OUT = "Seno Pamungkas Rahman - CV"
BW_CV_OUT = "Seno Pamungkas Rahman - CV (BW)"

all: generate main_cv bw_cv

generate:
	@echo "Generating LaTeX sections from data file..."
	$(PYTHON) generate_sections.py

main_cv: folders
	@echo "Compiling Main CV"
	$(ENGINE) -output-directory=out -output-format=pdf -jobname $(MAIN_CV_OUT) main_cv/main.tex

bw_cv: folders
	@echo "Compiling Black and White CV"
	$(ENGINE) -output-directory=out -output-format=pdf -jobname $(BW_CV_OUT) bw_cv/main.tex

folders:
	@mkdir -p out

clean:
	$(RM) -rf out

.PHONY: all generate main_cv bw_cv folders clean
