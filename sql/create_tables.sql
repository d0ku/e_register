DROP TABLE PrzepisyOdBabci;
DROP TABLE Users;

CREATE TABLE Users (
    id SERIAL,
    username text UNIQUE NOT NULL,
    password text NOT NULL,
    privileges VARCHAR(10)
);

CREATE TABLE PrzepisyOdBabci (
    username text REFERENCES Users(username),
    salt CHAR(7) NOT NULL
);

CREATE TABLE SessionID (
    session_id text NOT NULL UNIQUE,
    username text NOT NULL REFERENCES Users(username)
);
