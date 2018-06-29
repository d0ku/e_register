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
