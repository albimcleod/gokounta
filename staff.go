package gokounta

//Staff is the struct for a Kounta Staff
type Staff struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"primary_email_address"`
}

//Staffs is the struct for a list of Staff
type Staffs []Staff
