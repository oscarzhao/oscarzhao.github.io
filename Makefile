all: build

deps:
	gem update github-pages

build: clean deps
	jekyll build

clean:
	rm -rf _site

serve: clean
	jekyll serve

.PHONY: all deps build clean