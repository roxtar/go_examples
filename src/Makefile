SRC=$(wildcard *.go)
OUT=$(patsubst %.go, %.out, $(SRC))
LINK=$(OUT:.out=.6)
all: $(OUT)

%.out:%.6
	6l -o $@ $<

%.6:%.go
	6g $<

clean:
	rm -rf $(OUT) $(LINK)

.PHONY:clean