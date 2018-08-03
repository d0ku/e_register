#!/bin/sh

# Following script sets up PostgresSQL database with provided name as provided user.

# $1 -> database_name
# $2 -> populate_with_test_data, defaults to 0. To run it you have to provide 1 as 2nd argument.
# $3 -> username, use it only when you want to create user for database.
# $4 -> password for new user in postgres, same as above.

if [ $# -lt 1 ]
then
    echo Usage: ./setUp.sh database_name "[populate_with_test]" "[username]" "[user_password]"
    exit
fi

# Create database with provided name as superuser.
psql -c 'CREATE DATABASE '"$1"';' -U postgres

# Run all configuration files on previously created database.
psql -d $1 -f ./extensions.sql -U postgres
psql -d $1 -f ./create_tables.sql -U postgres
psql -d $1 -f ./procedures.sql -U postgres

if [ "$2" != "" ]
then
    if [ "$2" == "1" ]
    then
        psql -d $1 -f ./test_commands.sql -U postgres
    fi
fi

# Create database user and give him permissions. (only if both password and username are present)
if [ "$3" != "" ]
then
    if [ "$4" != "" ]
    then
        psql -c "CREATE USER "$3" password '"$4"'" -U postgres
        # TODO: This could and should be changed later to give only REALLY necessary permissions.
        psql -d $1 -c "GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO "$3";" -U postgres
        psql -d $1 -c "GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO "$3";" -U postgres
    fi
fi
