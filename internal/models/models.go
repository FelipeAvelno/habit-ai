package models

import (
	"log"
	"time"

	"habit-ai/pkg"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
    ID        string `gorm:"type:uuid;default:gen_random_uuid();primary_key"`
    FullName  string    `gorm:"not null" json:"full_name"`
    Email     string    `gorm:"unique;not null" json:"email"`
    Password  string    `gorm:"column:senha;not null" json:"password"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func Migrate() {
	err := pkg.DB.AutoMigrate(&User{})
	if err != nil {
		log.Fatal("Erro ao migrar tabelas:", err)
	}
	log.Println("Migração das tabelas concluída")
}