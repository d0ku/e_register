CREATE OR REPLACE FUNCTION add_user(IN user_name text, password text, usertype user_type, final_id integer) RETURNS BOOLEAN AS $$
DECLARE
temp integer;
BEGIN
    --Lock should prevent race condition (both connections check that user does not exist, then one can add and second can't add.).
    LOCK Users IN EXCLUSIVE MODE;
    IF EXISTS(  SELECT Users.username
        FROM Users
        WHERE Users.username=user_name
        )
        THEN
        -- could not create user, because user with such username already exists.
        RETURN FALSE;
    ELSE
        --user can be created.
        INSERT INTO Users("username", "hashed_password", "user_type", "final_id") VALUES(user_name, crypt(password, gen_salt('bf',8)), usertype, final_id);
        RETURN TRUE;
    END IF;
END
$$ language plpgsql;

CREATE TYPE user_data AS (user_exists BOOLEAN, user_type user_type, id integer);

CREATE OR REPLACE FUNCTION check_login_data(IN provided_username text, provided_password text) RETURNS user_data AS $$
DECLARE
return_data user_data;
BEGIN
    SELECT user_type, final_id
    INTO return_data.user_type, return_data.id
    FROM Users
    WHERE username=provided_username
    AND hashed_password = crypt(provided_password, hashed_password);
    return_data.user_exists=true;

    IF return_data.id IS NULL OR return_data.user_type IS NULL
        THEN
        return_data.user_exists=false;
    END IF;

    RETURN return_data;
END
$$ language plpgsql;
