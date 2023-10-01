package postgres

import (
	"context"
	"encoding/base64"
	"strconv"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(hash, password []byte) error {
	err := bcrypt.CompareHashAndPassword(hash, password)
	return err
}

func NewRepository(db *gorm.DB, ctx *context.Context) *Repository {
	return &Repository{
		Db:  db,
		Ctx: ctx,
	}
}

func (r *Repository) Connect(dbHost string, dbUser string, dbPassword string, dbName string, dbPort int, sslmode string) (*gorm.DB, error) {
	dsn := "host=" + dbHost + " user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " port=" + strconv.Itoa(dbPort) + " sslmode=" + sslmode
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func (r *Repository) Migrate() error {
	return r.Db.WithContext(*r.Ctx).AutoMigrate(&User{})
}

func (r *Repository) Close() error {
	sqlDB, err := r.Db.WithContext(*r.Ctx).DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (r *Repository) CreateUser(user *CreateUserRequest) error {
	hashpass, err := HashPassword(user.Password)
	if err != nil {
		return err
	}
	return r.Db.Create(&User{Username: user.Username, Email: user.Email, Password: hashpass}).Error
}

func (r *Repository) GetUser(user *UserRequest) (*User, error) {
	dbUser := &User{}
	err := r.Db.WithContext(*r.Ctx).Where("username = ?", user.Username).First(dbUser).Error
	return dbUser, err
}

func (r *Repository) Login(user *LoginRequest) error {
	dbUser := &User{}
	err := r.Db.WithContext(*r.Ctx).Where("username = ?", user.Username).First(dbUser).Error
	if err != nil {
		return err
	}
	convertedPassword, err := base64.StdEncoding.DecodeString(user.Password)
	if err != nil {
		return err
	}
	err = CheckPasswordHash([]byte(dbUser.Password), convertedPassword)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteUser(user *DeleteUserRequest) error {
	return r.Db.WithContext(*r.Ctx).Where("username = ?", user.Username).Delete(&User{}).Error
}

func (r *Repository) UpdateUser(user *UpdateUserRequest) error {
	var hashpass string
	var err error
	if user.Password != "" {
		hashpass, err = HashPassword(user.Password)
		if err != nil {
			return err
		}
	}
	return r.Db.WithContext(*r.Ctx).Model(&User{}).Where("username = ?", user.Username).Updates(User{Email: user.Email, Password: hashpass}).Error
}

func (r *Repository) GetAllUsers() ([]User, error) {
	var users []User
	err := r.Db.WithContext(*r.Ctx).Find(&users).Error
	return users, err
}

func (r *Repository) Ping() error {
	sqlDB, err := r.Db.WithContext(*r.Ctx).DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func (r *Repository) CheckUserPassword(user *LoginRequest) error {
	dbUser := &User{}
	err := r.Db.WithContext(*r.Ctx).Where("username = ?", user.Username).First(dbUser).Error
	if err != nil {
		return err
	}
	err = CheckPasswordHash([]byte(user.Password), []byte(dbUser.Password))
	if err != nil {
		return err
	}
	return nil
}
