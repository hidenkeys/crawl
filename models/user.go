package models

import "time"

type User struct {
	BaseModel
	FirstName      string  `gorm:"size:100;not null" json:"first_name"`
	LastName       string  `gorm:"size:100;not null" json:"last_name"`
	Username       string  `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Email          string  `gorm:"size:255;uniqueIndex;not null" json:"email"`
	PhoneNumber    string  `gorm:"size:20" json:"phone_number"`
	HashedPassword string  `gorm:"size:255;not null" json:"-"`
	ProfileImage   string  `gorm:"size:255" json:"profile_image_url"`
	Bio            string  `gorm:"type:text" json:"bio"`
	IsArtist       bool    `gorm:"default:false" json:"is_artist"`
	Roles          []Role  `gorm:"many2many:user_roles;joinForeignKey:UserID;joinReferences:RoleID" json:"roles,omitempty"`
	ArtistProfile  *Artist `gorm:"foreignKey:UserID" json:"artist_profile,omitempty"`
}

type Role struct {
	BaseModel
	Name        string `gorm:"size:50;uniqueIndex" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	Users       []User `gorm:"many2many:user_roles;joinForeignKey:RoleID;joinReferences:UserID" json:"users,omitempty"`
}

type UserRole struct {
	UserID    uint      `gorm:"primaryKey" json:"user_id"`
	RoleID    uint      `gorm:"primaryKey" json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
}
