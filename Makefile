.PHONY: example demo

example:
	npm --prefix alpine run build
	templ generate
	go run ./example

demo: example
