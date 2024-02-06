all: build install

build:
	go build -v ./cmd/lab
	go build -v ./cmd/gradelab
	
install:

	install /opt/ulab

docs:
	pandoc --standalone --to man doc/lab.md -o doc/lab.1
	gzip doc/lab.1