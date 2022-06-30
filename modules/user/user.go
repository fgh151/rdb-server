package user

import (
	"db-server/server/db"
	"db-server/utils"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// swagger:model
type User struct {
	// The user UUID
	// example: 6204037c-30e6-408b-8aaa-dd8219860b4b
	Id uuid.UUID `gorm:"primarykey" json:"id"`
	// User email
	Email string `gorm:"index" json:"email"`
	// Auth token
	Token string `gorm:"index" json:"token"`
	// Password hash
	PasswordHash string `json:"-"`
	// Is user admin
	Admin bool `gorm:"index;default:false;type:bool" json:"admin"`
	// Is user active
	Active bool `gorm:"index;default:true;type:bool" json:"active"`
	// Created at date time
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	// Last login date time
	LastLogin *time.Time `gorm:"last_login" json:"last_login,omitempty"`
	// user devices
	Devices []UserDevice `gorm:"foreignKey:user_id" json:"devices"`
}

// TableName Gorm table name
func (p User) TableName() string {
	return "user"
}

func (p User) List(limit int, offset int, sort string, order string, filter map[string]string) ([]interface{}, error) {
	var users []User

	db.MetaDb.ListQuery(limit, offset, sort, order, filter, users, []string{"Devices"})

	y := make([]interface{}, len(users))
	for i, v := range users {
		y[i] = v
	}

	return y, nil
}

func (p User) Total() *int64 {
	return db.MetaDb.TotalRecords(&User{})
}

func (p User) GetById(id string) (interface{}, error) {
	var user User

	conn := db.MetaDb.GetConnection()

	tx := conn.Preload("Devices").First(&user, "id = ?", id)

	if tx.RowsAffected < 1 {
		return user, errors.New("no found")
	}

	return user, nil
}

func (p User) GetByEmail(email string) (interface{}, error) {
	var user User

	conn := db.MetaDb.GetConnection()

	tx := conn.First(&user, "email = ?", email)

	if tx.RowsAffected > 0 {
		return user, nil
	}

	return user, errors.New("no found")
}

func (p User) Delete(id string) {
	conn := db.MetaDb.GetConnection()
	conn.Where("id = ?", id).Delete(&p)
}

func (p User) ValidatePassword(password string) bool {
	return p.PasswordHash == utils.HashPassword(password)
}

func (p User) UpdateLastLogin() {
	now := time.Now()
	p.LastLogin = &now
	db.MetaDb.GetConnection().Save(&p)
}

// swagger:model
type CreateUserForm struct {
	// new User email
	Email string `json:"Email"`
	// new User password
	Password string `json:"Password"`
}

func (f CreateUserForm) Save() User {
	var u = User{
		Email:        f.Email,
		PasswordHash: utils.HashPassword(f.Password),
		Token:        utils.GenerateRandomString(15),
	}

	u.Id, _ = uuid.NewUUID()

	db.MetaDb.GetConnection().Create(&u)

	return u
}

// swagger:model
type LoginForm struct {
	// User email
	Email string `json:"email"`
	// User password
	Password string `json:"password"`
}

func (f LoginForm) AdminLogin() (User, error) {
	return f.login(&User{Email: f.Email, Admin: true, Active: true})
}

func (f LoginForm) ApiLogin() (User, error) {
	return f.login(&User{Email: f.Email, Active: true})
}

func (f LoginForm) login(condition *User) (User, error) {
	var login User

	db.MetaDb.GetConnection().Where(condition).First(&login)

	if !login.ValidatePassword(f.Password) {
		return login, errors.New("invalid login or password")
	}

	login.UpdateLastLogin()

	return login, nil
}
