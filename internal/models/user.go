package models

type UserRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
	Dob  string `json:"dob" validate:"required,datetime=2006-01-02"`
}

type UserResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	Dob  string `json:"dob"`
	Age  *int   `json:"age,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
