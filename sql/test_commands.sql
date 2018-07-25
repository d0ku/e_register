SELECT add_user('testUser','testPassword','teacher',1); --should return true (correct)
SELECT add_user('testUser','testPassword','teacher',1); --should return false (duplicate)

SELECT check_login_data('testUser','testPassword', 'teacher'); --should return true,teacher,1
SELECT check_login_data('testUserFailed','testPassword', 'parent'); --should return false,-,-
SELECT check_login_data('testUser','testPassword', 'parent'); --should return false,-,- (bad user_type)

SELECT add_session('abba','makro'); --should return true
SELECT add_session('abba','makro'); --should return false

SELECT delete_session('abba'); --should return true

SELECT add_semester('summer',2018); --should return 1 (correct)
SELECT add_semester('winter',2018); --should return 2 (correct)
SELECT add_semester('winter',2018); --should return -1 (duplicate)

SELECT get_address_id('Paradise City','Imagined','12a','3','01-231'); --should return 1 (correct, create)
SELECT get_address_id('Paradise City','Imagined','12a','3','01-232'); --should return 2 (correct, create)
SELECT get_address_id('Paradise City','Imagined','12a','3','01-231'); --should return 1 (correct, find)

SELECT add_subject('Science'); --should return 1 (correct)
SELECT add_subject('Math'); --should return 2 (correct)
SELECT add_subject('Science'); --should return -1 (duplicate)

SELECT add_school('dummySchool',1,'LO','standard'); --should return 1 (correct)
SELECT add_school('dummySchool',1,'TECH','standard'); --should return 2 (correct)
SELECT add_school('dummySchool',1,'LO','standard'); --should return -1 (duplicate)

SELECT add_teacher(1,'tempFor','tempSur','mgr.','male','Jan-18-1923'); --should return 1 (correct)
SELECT add_teacher(1,'tempFor','tempSur','mgr.','male','Jan-18-1924'); --should return 2 (correct)
SELECT add_teacher(1,'tempFor','tempSur','mgr.','male','Jan-18-1923'); -- should return -1 (duplicate)

SELECT add_teacher_to_school(1,1); --should return t (correct)
SELECT add_teacher_to_school(1,1); --should return f (duplicate)
SELECT add_teacher_to_school(11,1); --should return f (non-existing teacher)
SELECT add_teacher_to_school(1,21); --should return f (non-existing school)

SELECT add_school_admin(1,1); --should return t (existing teacher to existing school)
SELECT add_school_admin(1,1); --should return f (duplicate)
SELECT add_school_admin(1,23); --should return f (existing teacher to not existing school)
SELECT add_school_admin(23, 1); --should return f (not-existing teacher to existing school)
SELECT add_school_admin(23, 23); --should return f (not-existing teacher to not-existing school)

SELECT add_class(1,1,'Oct-1-2007','Oct-1-2010','a',1); --should return 1 (correct)
SELECT add_class(1,1,'Oct-1-2009','Oct-1-2012','c',1); --should return 1 (correct)
SELECT add_class(1,1,'Oct-1-2007','Oct-1-2010','a',1); --should return -1 (duplicate)

SELECT add_student(1,1,1,'John','Snow','male','Jan-12-1998'); --should return 1 (correct)
SELECT add_student(1,1,1,'John','Snow','male','Jan-12-1998'); --should return 2 (correct, two students in same class with same names and same born_date can happen)
SELECT add_student(23,1,1,'Felicia','Snow','female','Jan-12-1998'); --should return -1 (non-existing school)
SELECT add_student(1,23,1,'Felicia','Snow','female','Jan-12-1998'); --should return -1 (non-existing class)
SELECT add_student(1,1,23,'Felicia','Snow','female','Jan-12-1998'); --should return -1 (non-existing address)

SELECT add_grade(1,1,1,'4.0'); --should return t (correct)
SELECT add_grade(1,1,8,'3.0'); --should return f (non-existing subject)
SELECT add_grade(1,8,1,'3.0'); --should return f (non-existing teacher)
SELECT add_grade(8,1,1,'3.0'); --should return f (non-existing student)

SELECT add_lesson(1,1,'08:00','08:45','Monday'); --should return 1 (correct)
SELECT add_lesson(1,1,'08:00','08:45','Monday'); --should return 2 (duplicate, but allowed)
SELECT add_lesson(1,5,'08:00','08:45','Monday'); --should return -1 (non-existing school)
SELECT add_lesson(2,5,'08:00','08:45','Monday'); --should return -1 (non-existing subject)

SELECT add_teacher_to_lesson(1,1); --should return t (correct)
SELECT add_teacher_to_lesson(1,1); --should return f (duplicate)

--TODO: there is a potential for UczniowieLekcje table where uczniowie can be assigned lessons indivdiually (e.g. language lessons where classes consist of people frmo many classes).
