DB_NAME = test_database
DB_USERNAME = postgres
DB_ADD_TEST_DATA = 0
BINARY_FILE_NAME = e_register_bin
SASS_DEV_PATH = page/SASS/
JS_DEV_PATH = page/js/

all: deps sql css js app tests

app:
	go build -o $(BINARY_FILE_NAME)

sql:
	cd sql && make DB_NAME=$(DB_NAME) DB_USERNAME=$(DB_USERNAME) DB_ADD_TEST_DATA=$(DB_ADD_TEST_DATA)

css:
	cd $(SASS_DEV_PATH) && make

js:
	cd $(JS_DEV_PATH) && make

deps:
	go get github.com/tdewolff/minify/cmd/minify
	go get github.com/lib/pq

clean:
	-rm *~
	cd $(SASS_DEV_PATH) && make clean

tests:
	go test ./...

.PHONY:	sql
