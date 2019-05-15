--TODO: display in app that you should changed your login data after first log on.

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
        INSERT INTO Users("username", "hashed_password", "user_type", "final_id", "change_password") VALUES(user_name, crypt(password, gen_salt('bf',8)), usertype, final_id, true);
        RETURN TRUE;
--        EXCEPTION WHEN integrity_constraint_violation THEN
--           RETURN FALSE;
    END IF;
END
$$ language plpgsql;

--That type is returned by function defined below.
CREATE TYPE login_data AS (username text, password text);

--This function is called whenever new parent/teacher/schoolAdmin/student is created and new username and password has to be generated.
CREATE OR REPLACE FUNCTION add_new_user(IN forename text, surname text, usertype user_type, user_id integer) RETURNS login_data AS $$ 
DECLARE 
    new_username text;
    password text;
    counter integer;
    login_try boolean;
    return_data login_data;
BEGIN
    --TODO: the complexity is comparable to going through a linked list, can it be lowered down?
    new_username = forename || '.' || surname ;
    counter = 0;
    password = md5(random()::text);

    return_data.username = new_username;
    return_data.password = password;

    --Spin around until username is found.
    LOOP
    SELECT INTO login_try add_user(new_username, password , usertype, user_id);
    IF login_try THEN
        RETURN return_data;
    ELSE
        counter = counter + 1;
        new_username = forename || '.' || surname || counter;
        return_data.username = new_username;
    END IF;
    END LOOP;
END
$$ language plpgsql;

CREATE TYPE user_data AS (user_exists BOOLEAN, user_type user_type, id integer, change_password_out BOOLEAN);

CREATE OR REPLACE FUNCTION check_login_data(IN provided_username text, provided_password text, provided_user_type user_type) RETURNS user_data AS $$
DECLARE
return_data user_data;
BEGIN
    SELECT user_type, final_id, change_password
    INTO return_data.user_type, return_data.id, return_data.change_password_out
    FROM Users
    WHERE username = provided_username
    AND hashed_password = crypt(provided_password, hashed_password)
    AND user_type = provided_user_type;

    return_data.user_exists=true;

    IF return_data.id IS NULL OR return_data.user_type IS NULL
        THEN
        return_data.user_exists=false;
    END IF;

    RETURN return_data;
END
$$ language plpgsql;

CREATE OR REPLACE FUNCTION add_semester(IN sem_type semester_type, year integer) RETURNS integer AS $$
DECLARE
--return id of semester, -1 if can't be added.
    to_return integer;
BEGIN
    INSERT INTO Semesters("semester", "year") VALUES(sem_type, year) RETURNING id_semester INTO to_return;
    RETURN to_return;
    EXCEPTION WHEN integrity_constraint_violation THEN
        RETURN -1;
END
$$ language plpgsql;



--Call this function when you want to create new student.
--TODO: automatically add this student to Users so he can log in.
CREATE OR REPLACE FUNCTION add_student(IN school_id integer, class_id integer, address_id integer, forename_input text, surname_input text, sex_input sex_type, born_date DATE) RETURNS integer AS $$
--If student was successfully added -> return his id, else return -1 (-1 as response from that function means there was some error in the way);
DECLARE
    to_return integer;
BEGIN
    -- TODO: test it thoroughly, when done with add_address, add_class, add_school,.
-- SELECT add_student(1,1,1,'testFore','testSure','male','Jan-08-1999');
    INSERT INTO Students("id_school","id_class","id_address","forename","surname","sex","date_of_birth") VALUES(school_id, class_id, address_id, forename_input, surname_input, sex_input, born_date) RETURNING id_student INTO to_return;
    RETURN to_return;
    EXCEPTION WHEN integrity_constraint_violation THEN
        RETURN -1;
END
$$ language plpgsql;


CREATE OR REPLACE FUNCTION get_address_id(IN city_in text, street_in text, house_number_in text, flat_number_in text, postal_code_in text) RETURNS integer AS $$
--It creates address if it does not already exist, if it exists it returns id.
--If address was successfully added -> return his id, else return -1 (-1 as response from that function means there was some error in the way);
DECLARE
    to_return integer;
BEGIN
    LOCK Addresses IN EXCLUSIVE MODE;
    IF EXISTS( SELECT id_address 
        FROM Addresses
        WHERE city=city_in AND
            street=street_in AND
            house_number=house_number_in AND
            flat_number=flat_number_in AND
            postal_code=postal_code_in
        )
        THEN
        --return id of address that exists.
        SELECT id_address
        INTO to_return
        FROM Addresses
        WHERE city=city_in AND
            street=street_in AND
            house_number=house_number_in AND
            flat_number=flat_number_in AND
            postal_code=postal_code_in;
        RETURN to_return;

    ELSE
        --address has to be added.
        INSERT INTO Addresses("city","street","house_number","flat_number","postal_code") VALUES(city_in, street_in, house_number_in, flat_number_in, postal_code_in) RETURNING id_address INTO to_return;
        RETURN to_return;
    END IF;
END
$$ language plpgsql;

CREATE OR REPLACE FUNCTION add_subject(IN name text) RETURNS integer AS $$
DECLARE
    to_return integer;
BEGIN
    --try to add subject.
    INSERT INTO Subjects("name") VALUES(name) RETURNING id_subject INTO to_return;
    RETURN to_return;
    --such subject can't be added, return -1.
    EXCEPTION WHEN integrity_constraint_violation THEN
        RETURN -1;
END
$$ language plpgsql;

CREATE OR REPLACE FUNCTION add_school(IN name_in text, address_id integer, school_type school_type, bought_package offer_type ) RETURNS integer AS $$
-- return id of created school if can be created, -1 otherwise.
DECLARE
    to_return integer;
BEGIN
    --try to add school.
    INSERT INTO Schools("full_name","id_address","typ","bought_offer") VALUES(name_in, address_id, school_type, bought_package) RETURNING id_school INTO to_return;
    RETURN to_return;
    --such subject can't be added, return -1.
    EXCEPTION WHEN integrity_constraint_violation THEN
        RETURN -1;
END
$$ language plpgsql;

--Call this function when you want to give teacher admin privileges.
--TODO: automatically add to Users.
CREATE OR REPLACE FUNCTION add_school_admin(IN teacher_id integer, school_id integer) RETURNS BOOLEAN AS $$
--return true if successfully added, false otherwise.
DECLARE
BEGIN
    --try to add school admin.
    INSERT INTO SchoolsAdministrators("id_teacher","id_school") VALUES(teacher_id, school_id);
    RETURN TRUE;
    --if there is an error, return false.
    EXCEPTION WHEN integrity_constraint_violation THEN
        RETURN FALSE;
END
$$ language plpgsql;

--Call this function when you want to add completely new teacher instance.
--TODO: automatically add to Users.
CREATE OR REPLACE FUNCTION add_teacher(IN address_id integer, forename text, surname text, degree teacher_rank, sex_in sex_type, born_date DATE) RETURNS INTEGER AS $$
--Return id if succesfully added, -1 otherwise.
DECLARE
    to_return integer;
BEGIN
    INSERT INTO Teachers("id_address", "forename", "surname", "teacher_rank", "sex", "date_of_birth") VALUES(address_id, forename, surname, degree, sex_in, born_date) RETURNING id_teacher INTO to_return;
    RETURN to_return;
    --if there is an error, return -1.
    EXCEPTION WHEN integrity_constraint_violation THEN
        RETURN -1;
END
$$ language plpgsql;

CREATE OR REPLACE FUNCTION add_teacher_to_school(IN teacher_id integer, school_id integer) RETURNS BOOLEAN AS $$
BEGIN
    INSERT INTO TeachersSchools("id_teacher","id_school") VALUES(teacher_id, school_id);
    RETURN TRUE;
    EXCEPTION WHEN integrity_constraint_violation THEN
        RETURN FALSE;
END
$$ language plpgsql;

CREATE OR REPLACE FUNCTION add_grade(IN student_id integer, teacher_id integer, subject_id integer, grade integer) RETURNS BOOLEAN AS $$
BEGIN
    INSERT INTO Grades("id_student","id_teacher","id_subject","grade") VALUES(student_id, teacher_id, subject_id, grade);
    RETURN TRUE;
    EXCEPTION WHEN integrity_constraint_violation THEN
        RETURN FALSE;
END
$$ language plpgsql;

CREATE OR REPLACE FUNCTION add_class(IN school_id integer, teacher_id integer, start_date DATE, end_date DATE, char_in text, num_in integer) RETURNS INTEGER AS $$
DECLARE
    to_return integer;
BEGIN
    INSERT INTO Classes("id_school","id_educator","start_date","end_date","letter","class_level") VALUES(school_id, teacher_id, start_date, end_date, char_in, num_in) RETURNING id_class INTO to_return;
    RETURN to_return;
    EXCEPTION WHEN integrity_constraint_violation THEN
        RETURN -1;
END
$$ language plpgsql;

CREATE OR REPLACE FUNCTION add_lesson(IN id_subject_in integer, id_school_in integer, start_hour time, end_hour time, day_name day_type) RETURNS INTEGER AS $$
--If successfully added return lesson_id, else return -1.
DECLARE
    to_return integer;
BEGIN
    INSERT INTO Lessons("id_subject", "id_school","start_hour","end_hour", "day") VALUES(id_subject_in, id_school_in, start_hour, end_hour, day_name) RETURNING id_lesson INTO to_return;
        RETURN to_return;
    EXCEPTION WHEN integrity_constraint_violation THEN
        RETURN -1;
END
$$ language plpgsql;

CREATE OR REPLACE FUNCTION add_teacher_to_lesson(IN lesson_id integer, teacher_id integer) RETURNS BOOLEAN AS $$
BEGIN
    --TODO: write trigger for inserting into that table, that does not allow one teacher to have two lessons at same time. It is possible that exception catch should be changed then. Teachers should be allowed to have lessons at different schools, part-time job etc.
    INSERT INTO LessonsTeachers("id_lesson", "id_teacher") VALUES(lesson_id, teacher_id);
    RETURN TRUE;
    EXCEPTION WHEN integrity_constraint_violation THEN
        RETURN FALSE;
END
$$ language plpgsql;

CREATE OR REPLACE FUNCTION check_if_teacher_is_school_admin(IN teacher_id integer) RETURNS INTEGER AS $$
DECLARE
    to_return integer;
BEGIN
    --TODO: consider that one teacher can be admin in many schools

    SELECT id_school FROM SchoolsAdministrators WHERE id_teacher=teacher_id INTO to_return;
    IF to_return IS NULL THEN
        RETURN -1;
    END IF;
    RETURN to_return;
END
$$ language plpgsql;

CREATE OR REPLACE FUNCTION get_schools_teacher_teaches_in(IN teacher_id integer ) RETURNS TABLE (
   id_teacher_out integer) AS $$
BEGIN
    RETURN QUERY
    SELECT id_school
    FROM TeachersSchools
    WHERE id_teacher=teacher_id;
END
$$ language plpgsql;

CREATE OR REPLACE FUNCTION get_schools_details_where_teacher_teaches(IN teacher_id integer) RETURNS TABLE (
    id_school_out integer,
    full_name_out text,
    city_out text,
    street_out text,
    school_type_out school_type)
    AS $$ 
BEGIN
    RETURN QUERY 
    SELECT Schools.id_school, Schools.full_name, Addresses.city, Addresses.street, Schools.typ 
    FROM Schools INNER JOIN TeachersSchools ON (Schools.id_school=TeachersSchools.id_school)
    INNER JOIN Addresses ON (Schools.id_address=Addresses.id_address)
    WHERE id_teacher=teacher_id;
END
$$ language plpgsql;

CREATE OR REPLACE FUNCTION add_student_to_lesson(IN lesson_id integer, student_id integer) RETURNS BOOLEAN AS $$
BEGIN
    -- TODO: does it make sense to check whether lessons connected with user overlap each other?
    INSERT INTO StudentsLessons("id_student", "id_lesson") VALUES (student_id, lesson_id);
    RETURN TRUE;
    EXCEPTION WHEN integrity_constraint_violation THEN
        RETURN FALSE;
END
$$ language plpgsql;

CREATE OR REPLACE FUNCTION add_class_to_lesson(IN lesson_id integer, class_id integer) RETURNS BOOLEAN AS $$
BEGIN
    -- TODO: does it make sense to check whether lessons connected with user overlap each other?
    INSERT INTO ClassesLessons("id_class", "id_lesson") VALUES (class_id, lesson_id);
    RETURN TRUE;
    EXCEPTION WHEN integrity_constraint_violation THEN
        RETURN FALSE;
END
$$ language plpgsql;

CREATE OR REPLACE FUNCTION add_warning(IN user_id integer, teacher_id integer, content text) RETURNS BOOLEAN AS $$
BEGIN
    INSERT INTO Warnings("id_student", "id_teacher", "content", "time") VALUES(user_id, teacher_id, content, now());
    RETURN TRUE;
    EXCEPTION WHEN integrity_constraint_violation THEN
        RETURN FALSE;
END
$$ language plpgsql;

CREATE OR REPLACE FUNCTION add_final_grade(IN student_id integer, teacher_id integer, lesson_id integer, value integer) RETURNS BOOLEAN AS $$
DECLARE 
    id_subject_get integer;
BEGIN
    --TODO: write some triggers, so only teacher who is connected to lesson can insert grades.
    SELECT id_subject
    INTO id_subject_get
    FROM Lessons
    WHERE Lessons.id_lesson=lesson_id;

    INSERT INTO FinalGrades("id_student", "id_teacher", "id_subject", "grade") VALUES(student_id, teacher_id, id_subject_get, value);
    RETURN TRUE;
    EXCEPTION WHEN integrity_constraint_violation THEN
        RETURN FALSE;
END
$$ language plpgsql;

CREATE OR REPLACE FUNCTION add_final_grade_behaviour(IN student_id integer, teacher_id integer, value integer) RETURNS BOOLEAN AS $$
BEGIN
    --TODO: implement trigger, so only class main teacher can insert final behaviour grades
    INSERT INTO StudentsBehaviourFinalGrades("id_student", "id_teacher", "value") VALUES(student_id, teacher_id, value);
    RETURN TRUE;
    EXCEPTION WHEN integrity_constraint_violation THEN
        RETURN FALSE;
END
$$ language plpgsql;

--Call this function when you want to say that a teacher is capable of teaching new subject.
CREATE OR REPLACE FUNCTION add_teacher_to_subject(IN teacher_id integer, subject_id integer) RETURNS BOOLEAN AS $$
BEGIN
    INSERT INTO TeachersSubjects("id_teacher", "id_subject") VALUES (teacher_id, subject_id);
    RETURN TRUE;
    EXCEPTION WHEN integrity_constraint_violation THEN
        RETURN FALSE;
END
$$ language plpgsql;

--Call this function when you want to create parent for a child who is already registered.
--TODO: add this new parent to Users automatically.
--Returns id of newly created parent, or -1 if could not be added.
CREATE OR REPLACE FUNCTION add_parent(IN address_id integer, surname text, forename text, sex sex_type) RETURNS INTEGER AS $$
DECLARE
    to_return integer;
BEGIN
    INSERT INTO Parents("id_address", "forename", "surname", "sex") VALUES(address_id, forename, surname, sex) RETURNING id_parent INTO to_return;
    RETURN to_return;
    EXCEPTION WHEN integrity_constraint_violation THEN
        RETURN -1;
END
$$ language plpgsql;
