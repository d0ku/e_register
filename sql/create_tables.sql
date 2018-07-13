-- TODO: big spell check
CREATE TYPE user_type AS ENUM('teacher','user','parent');
CREATE TYPE school_type AS ENUM('LO','POD','TECH', 'ZAW');
CREATE TYPE offer_type AS ENUM('standard','gold','diamond');
CREATE TYPE teacher_rank AS ENUM('inż.', 'mgr.', 'mgr. inż.', 'dr.'); -- TODO: write more
CREATE TYPE sex_type AS ENUM('male', 'female');
CREATE TYPE day_type AS ENUM('Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday');
CREATE TYPE presence_type AS ENUM('present', 'absent', 'late', 'justified');
CREATE TYPE grade_type AS ENUM('1.0', '1.5', '2.0', '2.5', '3.0', '3.5', '4.0', '4.5', '5.0', '5.5', '6.0');
CREATE TYPE semester_type AS ENUM('winter', 'summer');
CREATE TYPE behaviour_grades_type AS ENUM('1.0', '2.0', '3,0', '4.0', '5.0');
CREATE TABLE Users (
    username text UNIQUE NOT NULL,
    hashed_password text NOT NULL,
    user_type user_type NOT NULL,
    final_id integer NOT NULL
);

CREATE TABLE SessionID (
    session_id text NOT NULL UNIQUE,
    username text NOT NULL references Users(username)
);

CREATE TABLE Semesters (
    id_semester SERIAL PRIMARY KEY,
    semester semester_type NOT NULL,
    year integer NOT NULL,
    UNIQUE(semester, year)
);

CREATE TABLE Addresses (
    id_address SERIAL PRIMARY KEY,
    city text NOT NULL,
    street text NOT NULL,
    house_number text NOT NULL,
    flat_number text,
    postal_code text NOT NULL,
    UNIQUE(city,street,house_number,flat_number,postal_code)
);

CREATE TABLE Schools (
    id_school SERIAL PRIMARY KEY,
    full_name text NOT NULL,
    id_address integer references Addresses(id_address),
    typ school_type NOT NULL,
    bought_offer offer_type NOT NULL,
    UNIQUE(full_name,id_address,typ,bought_offer)
);

CREATE TABLE SchoolsNumbers (
    id_school SERIAL references Schools(id_school),
    phone_number text NOT NULL,
    PRIMARY KEY(id_school, phone_number)
);

CREATE TABLE Teachers (
    id_teacher SERIAL PRIMARY KEY,
    id_school integer references Schools(id_school),
    id_address integer references Addresses(id_address),
    forename text NOT NULL,
    surename text NOT NULL,
    teacher_rank teacher_rank NOT NULL,
    sex sex_type NOT NULL,
    date_of_birth date,
    UNIQUE(id_school,id_address,forename,surename,teacher_rank,sex,date_of_birth)
);

CREATE TABLE TeachersNumbers (
    id_teacher integer references Teachers(id_teacher),
    phone_number text NOT NULL,
    PRIMARY KEY(id_teacher, phone_number)
);

CREATE TABLE Subjects (
    id_subject SERIAL PRIMARY KEY,
    name text NOT NULL UNIQUE
);

CREATE TABLE TeachersSubjects (
    id_teacher integer references Teachers(id_teacher),
    id_subject integer references Subjects(id_subject),
    PRIMARY KEY(id_teacher, id_subject)
);

CREATE TABLE SchoolsAdministrators (
    id_teacher integer references Teachers(id_teacher),
    id_school integer references Schools(id_school),
    PRIMARY KEY(id_teacher, id_school)
);

CREATE TABLE Classes (
    id_class SERIAL PRIMARY KEY,
    id_school integer references Schools(id_school),
    id_educator integer references Teachers(id_teacher),
    start_date date NOT NULL,
    end_date date NOT NULL CHECK(end_date > start_date),
    letter VARCHAR(3) NOT NULL, -- TI, a, b
    class_level integer NOT NULL, -- 1,2,3...,
    UNIQUE(id_school,id_educator,start_date,end_date,letter,class_level)
);

CREATE TABLE Lessons (
    --TODO: consider creating additional table with lesson times, and just referencing by id here?
    id_lesson SERIAL PRIMARY KEY,
    id_subject integer references Subjects(id_subject),
    start_hour time NOT NULL,
    end_hour time NOT NULL,
    day day_type NOT NULL
);

CREATE TABLE ClassesLessons (
    id_class integer references Classes(id_class),
    id_subject integer references Subjects(id_subject),
    PRIMARY KEY(id_class, id_subject)
);

CREATE TABLE LessonsTeachers (
    id_lesson integer references Lessons(id_lesson),
    id_teacher integer references Teachers(id_teacher),
    PRIMARY KEY(id_lesson, id_teacher)
);

CREATE TABLE Students (
    id_student SERIAL PRIMARY KEY,
    id_school integer references Schools(id_school),
    id_class integer references Classes(id_class),
    id_address integer references Addresses(id_address),
    surename text NOT NULL,
    forename text NOT NULL,
    sex sex_type NOT NULL,
    date_of_birth date NOT NULL
);

CREATE TABLE StudentsLessons (
    id_student integer references Students(id_student),
    id_lesson integer references Lessons(id_lesson),
    PRIMARY KEY(id_student, id_lesson)
);

CREATE TABLE StudentsNumbers (
    id_student integer references Students(id_student),
    phone_number text NOT NULL
);

CREATE TABLE StudentsInformation (
    id_student integer references Students(id_student),
    information text NOT NULL,
    PRIMARY KEY(id_student, information)
);

CREATE TABLE LessonsPresence (
    id_lesson integer references Lessons(id_lesson),
    id_student integer references Students(id_student),
    status presence_type NOT NULL,
    id_insertor integer references Teachers(id_teacher)
);

CREATE TABLE Warnings (
    -- No reference to lesson, because warnings can be added by teachers who don't have lessons with student.
    id_warning SERIAL PRIMARY KEY,
    id_student integer references Students(id_student),
    id_teacher integer references Teachers(id_teacher),
    content text NOT NULL,
    day date
);

CREATE TABLE Grades (
    id_grade SERIAL PRIMARY KEY,
    id_student integer references Students(id_student),
    id_teacher integer references Teachers(id_teacher),
    id_subject integer references Subjects(id_subject),
    grade grade_type NOT NULL
);
--Represents current semester.
CREATE TABLE FinalGrades (
    id_grade SERIAL PRIMARY KEY,
    id_student integer references Students(id_student),
    id_teacher integer references Teachers(id_teacher), --who added this grade.
    id_subject integer references Subjects(id_subject),
    grade grade_type NOT NULL
);

--Stores information about all previous semesters.
CREATE TABLE FinalGradesArchive (
id_grade SERIAL PRIMARY KEY,
    id_student integer references Students(id_student),
    id_teacher integer references Teachers(id_teacher), --who added this grade.
    id_subject integer references Subjects(id_subject),
    id_class integer references Classes(id_class),
    id_semester integer references Semesters(id_semester),
    grade grade_type NOT NULL
);

CREATE TABLE StudentsBehaviourGradesFinal (
    id_grade SERIAL PRIMARY KEY,
    id_student integer references Students(id_student),
    id_teacher integer references Teachers(id_teacher),
    value behaviour_grades_type NOT NULL
);

CREATE TABLE StudentsBehaviourGradesFinalArchive (
    id_grade SERIAL PRIMARY KEY,
    id_student integer references Students(id_student),
    id_teacher integer references Teachers(id_teacher),
    id_class integer references Classes(id_class),
    id_semester integer references Semesters(id_semester),
    value behaviour_grades_type NOT NULL
);

CREATE TABLE Parents (
    id_parent SERIAL PRIMARY KEY,
    id_address integer references Addresses(id_address),
    forename text NOT NULL,
    surename text NOT NULL,
    sex sex_type NOT NULL
);

CREATE TABLE ParentsNumbers (
    id_parent integer references Parents(id_parent),
    phone_number text NOT NULL
);

CREATE TABLE ParentsStudents (
    id_parent integer references Parents(id_parent),
    id_student integer references Students(id_student),
    PRIMARY KEY(id_parent, id_student)
);
-- 28
