package databasehandling

type UserLoginData struct {
	Exists    bool
	User_type string
	Id        int
	//Is set to true when user should change his password.
	ChangePassword bool
}

type School struct {
	Id         int
	FullName   string
	City       string
	Street     string
	Address    string
	SchoolType string
}
