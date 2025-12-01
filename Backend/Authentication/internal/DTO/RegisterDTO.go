package DTO

import ()


type RegisterRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
    ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
    Email string `json:"email" binding:"required,email"`
}
