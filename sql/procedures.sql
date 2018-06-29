CREATE OR REPLACE FUNCTION add_user(IN username text, password text) RETURNS VOID AS $$
DECLARE
BEGIN
    INSERT INTO Users("username","password") VALUES(username, crypt(password, gen_salt('bf',8)));
END
$$ language plpgsql;

CREATE OR REPLACE FUNCTION check_login_data(IN provided_username text, provided_password text) RETURNS text AS $$
DECLARE
BEGIN
   return (SELECT username FROM Users WHERE username = provided_username AND password = crypt(provided_password, password));
END
$$ language plpgsql;
