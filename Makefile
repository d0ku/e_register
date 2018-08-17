DB_NAME = e_register
DB_USERNAME = test_user
DB_PASSWORD = test_password
DB_ADD_TEST_DATA = 1
BINARY_FILE_NAME = e_register_bin
SASS_DEV_PATH = page/SASS/
JS_DEV_PATH = page/js/

all: deps sql css js app tests

app:
	go build -o $(BINARY_FILE_NAME)

sql:
	cd sql && make DB_NAME=$(DB_NAME) DB_USERNAME=$(DB_USERNAME) DB_ADD_TEST_DATA=$(DB_ADD_TEST_DATA) DB_PASSWORD=$(DB_PASSWORD)

css:
	cd $(SASS_DEV_PATH) && make

js:
	cd $(JS_DEV_PATH) && make

deps:
	go get github.com/tdewolff/minify/cmd/minify
	go get github.com/lib/pq
	npm install -g sass

clean:
	-rm *~
	cd $(SASS_DEV_PATH) && make clean

tests:
	go test ./...

.PHONY:	sql
