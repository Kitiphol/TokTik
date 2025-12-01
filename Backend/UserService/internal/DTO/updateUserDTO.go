package DTO

import ()


type UpdateUserRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
    ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}
