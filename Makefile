BASEDIR=$(CURDIR)

BINARY=zombie-ping
SOURCES := $(shell find $(BASEDIR) -name '*.go')
SCSS_SRC := $(shell find $(BASEDIR)/scss -name '*.scss')
SCSS_TARGET := $(shell find ./scss -name '*.scss' | sed 's/\.scss/\.css/' | sed 's/scss/static/')
TESTED=.tested

build: $(BINARY)
$(BINARY): $(SOURCES) $(TESTED) $(SCSS_TARGET)
	./tools/build

scss: $(SCSS_TARGET)
static/%.css: scss/%.scss
	pysassc $< $@

test: $(TESTED)
$(TESTED): $(SOURCES)
	./tools/test

clean:
	rm -f $(BINARY) debug debug.test cover.out $(TESTED) $(SCSS_TARGET)

.PHONY: clean test cover build run scss
