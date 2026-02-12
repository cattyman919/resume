.PHONY: watch

compile:
	typst compile src/main.typ --root .

watch:
	typst watch src/main.typ --root .
