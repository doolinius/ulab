all: build install

build:
	go build -v ./cmd/lab
	go build -v ./cmd/gradelab
	
install:

	install /opt/ulab
