.PHONY: all clean

all: sequence.pdf

clean:
	-rm sequence.aux sequence.latex sequence.log sequence.pdf

sequence.pdf: sequence.txt
	plantuml -tlatex sequence.txt
	pdflatex sequence.latex
	rm sequence.aux sequence.latex sequence.log
