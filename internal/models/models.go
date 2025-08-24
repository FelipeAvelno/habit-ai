package models

import (
	"log"
	"time"

	"habit-ai/pkg"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	FullName  string `gorm:"not null" json:"full_name"`
	Email     string `gorm:"unique;not null" json:"email"`
	Password  string `gorm:"column:senha;not null" json:"password"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	Habits []Habit `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // Relacionamento
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

type Habit struct {
	ID              string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID          string    `gorm:"type:uuid;not null" json:"user_id"`
	User           User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Nome            string    `gorm:"not null" json:"name"`
	Categoria       string    `json:"category"`
	HorarioPreferido string   `json:"preferred_hour"`
	Frequencia      int       `json:"frequency"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func Migrate() {
	err := pkg.DB.AutoMigrate(&User{}, &Habit{})
	if err != nil {
		log.Fatal("Erro ao migrar tabelas:", err)
	}
	log.Println("Migração das tabelas concluída")
}