OUT=presentation.pdf
SRC=presentation.tex
all: $(OUT)

$(OUT): $(SRC)
	pdflatex $<
	pdflatex $<

view:
	evince $(OUT) &

.PHONY: clean

clean:
	rm -rf *.aux *.log *.pdf