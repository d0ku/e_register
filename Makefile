all: sql css app tests

app:
	go build -o database_project_go

sql:
	ls
	#TODO

css:
	cd page/CSS && make

clean:
	rm *~
	cd page/CSS && make clean

tests:
	go test ./...
