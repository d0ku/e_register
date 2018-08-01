all: sql css app tests

app:
	go build -o database_project_go

sql:
	ls
	#TODO

css:
	cd page/SASS/ && make

clean:
	-rm *~
	cd page/SASS/ && make clean

tests:
	go test ./...
