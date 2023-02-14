package dto

import "time"

type AuthDto struct {
	Id            string    `json:"id"`
	UserName      string    `json:"userName"`
	Role          string    `json:"role"`
	LastLoginTime time.Time `json:"lastLoginTime"`
}
