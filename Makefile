# Tools
ENGINE = pdflatex
RM = rm -fr	

all: folders
	$(ENGINE) -output-directory=out -output-format=pdf  main.tex

folders:
ifneq ($(wildcard ./out/.*),)
	@echo "Found out directory"
else
	@echo "out directory not found"
	mkdir out
endif

clean: 
	$(RM) out
