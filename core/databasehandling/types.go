package databasehandling

type UserLoginData struct {
	Exists    bool
	User_type string
	Id        int
}

type School struct {
	Id         int
	FullName   string
	City       string
	Street     string
	Address    string
	SchoolType string
}
