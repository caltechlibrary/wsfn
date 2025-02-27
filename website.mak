#
# Makefile for running pandoc on all Markdown docs ending in .md
#
PROJECT = wsfn

BASE_URL =
ifneq ($(base_url),)
  BASE_URL = $(base_url)
endif

PANDOC = $(shell which pandoc)

MD_PAGES = $(shell ls -1 *.md | grep -v 'nav.md')

HTML_PAGES = $(shell ls -1 *.md | sed -E 's/\.md/.html/g')

build: search.md $(HTML_PAGES) $(MD_PAGES) pagefind

search.md: search.md.tmpl .FORCE
	echo "" | pandoc --metadata base_url=$(base_url) -s --to markdown -o search.md --template=search.md.tmpl; git add search.md

$(HTML_PAGES): $(MD_PAGES) .FORCE
	if [ -f $(PANDOC) ]; then $(PANDOC) --metadata title=$(basename $@) -s --to html5 $(basename $@).md -o $(basename $@).html \
		--lua-filter=links-to-html.lua \
	    --template=page.tmpl; fi
	@if [ $@ = "README.html" ]; then mv README.html index.html; fi

pagefind: .FORCE
	pagefind --verbose --exclude-selectors="nav,header,footer" --site .
	git add pagefind

clean:
	@if [ -f index.html ]; then rm *.html; fi

.FORCE:
