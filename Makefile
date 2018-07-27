all: sql css app

app:
	go build -o database_project_go

sql:
	ls
	#TODO

css:
	cd page/CSS && make

clean:
	cd page/CSS && make clean
