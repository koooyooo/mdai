PHONY: install

install:
	@ go build -o mdai mdai.go && sudo mv mdai /usr/local/bin/

run:
	@ go run mdai.go $(FILE)