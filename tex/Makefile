PRESENTATION=presentation.pdf
PRESENTATION_SRC=presentation.tex
MANUAL=manual.pdf
MANUAL_SRC=manual.tex
MANUAL_BIB=manual

all: $(PRESENTATION) $(MANUAL)

$(PRESENTATION): $(PRESENTATION_SRC)
	pdflatex $<
	pdflatex $<

$(MANUAL): $(MANUAL_SRC) 
	pdflatex $<
	bibtex $(MANUAL_BIB)
	pdflatex $<
	pdflatex $<

view:
	evince $(PRESENTATION) &

.PHONY: clean

clean:
	rm -rf *.aux *.log *.pdf