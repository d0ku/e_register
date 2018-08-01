all: deps sql css js app tests

app:
	go build -o database_project_go

sql:
	ls
	#TODO

css:
	cd page/SASS/ && make

js:
	cd page/js/ && make

deps:
	go get github.com/tdewolff/minify/cmd/minify
	go get github.com/lib/pq

clean:
	-rm *~
	cd page/SASS/ && make clean

tests:
	go test ./...
