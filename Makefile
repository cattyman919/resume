# Tools
ENGINE = pdflatex
RM = rm -fr	

all: main_cv bw_cv

main_cv: folders
	@echo "Compiling Main CV"
	$(ENGINE) -output-directory=out -output-format=pdf -jobname "Seno Pamungkas Rahman - CV" main_cv/main.tex

bw_cv: folders
	@echo "Compiling Black and White CV"
	$(ENGINE) -output-directory=out -output-format=pdf -jobname "Seno Pamungkas Rahman - CV (BW)" bw_cv/main.tex

folders:
ifneq ($(wildcard ./out/.*),)
	@echo "Found out directory"
else
	@echo "out directory not found"
	mkdir out
endif

clean: 
	$(RM) out
