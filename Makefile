PHONY: install

clean:
	@ rm -f mdai && rm -f mdai.exe && sudo rm -f /usr/local/bin/mdai

install: clean
	@ go build -o mdai mdai.go && sudo mv mdai /usr/local/bin/

run:
	@ go run mdai.go $(FILE)

less-config:
	@ less ~/.mdai/config.yml

rm-config:
	@ rm -f ~/.mdai/config.yml

init: rm-config
	@ mdai init