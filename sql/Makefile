all: sql

sql:	setUp.sh extensions.sql create_tables.sql procedures.sql test_commands.sql
	./setUp.sh $(DB_NAME) $(DB_USERNAME) $(DB_ADD_TEST_DATA)

clean:
	-rm -rf *~
