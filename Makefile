.PHONY: compile

compile:
	GOOS=linux GOARCH=386 go build -o bin/riptad
