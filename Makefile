# Tools
ENGINE = pdflatex
RM = rm -fr	

DOCNAME = resume

all: folders
	$(ENGINE) -output-directory=out  main.tex

folders:
ifneq ($(wildcard ./out/.*),)
	@echo "Found out directory"
else
	@echo "out directory not found"
	mkdir out
endif

clean: 
	$(RM) out
