package models

type User struct {
    ID       string
    Username string
    Password string
    Email    string
}

type UserLogin struct {
    Username string `json:"username" validate:"required"`
    Password string `json:"password" validate:"required"`
}

type UserRegister struct {
    Username string `json:"username" validate:"required"`
    Password string `json:"password" validate:"required"`
    Email    string `json:"email" validate:"required,email"`
}
