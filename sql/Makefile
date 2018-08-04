all: sql

sql:	setUp.sh extensions.sql create_tables.sql procedures.sql test_commands.sql
	./setUp.sh $(DB_NAME) $(DB_ADD_TEST_DATA) $(DB_USERNAME) $(DB_PASSWORD)

clean:
	-rm -rf *~
