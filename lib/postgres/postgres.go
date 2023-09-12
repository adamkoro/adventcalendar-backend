package postgres

import (
	"strconv"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Connect(dbHost string, dbUser string, dbPassword string, dbName string, dbPort int, sslmode string) (*gorm.DB, error) {
	dsn := "host=" + dbHost + " user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " port=" + strconv.Itoa(dbPort) + " sslmode=" + sslmode
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func (r *Repository) Migrate() error {
	return r.db.AutoMigrate(&User{})
}

func (r *Repository) Close() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (r *Repository) CreateUser(username string, email string, password string) error {
	hashpass, err := HashPassword(password)
	if err != nil {
		return err
	}
	return r.db.Create(&User{Username: username, Email: email, Password: hashpass}).Error
}

func (r *Repository) GetUser(username string) (*User, error) {
	user := &User{}
	err := r.db.Where("username = ?", username).First(user).Error
	return user, err
}

func (r *Repository) Login(username string, password string) error {
	user := &User{}
	err := r.db.Where("username = ?", username).First(user).Error
	if err != nil {
		return err
	}
	err = CheckPasswordHash(password, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteUser(username string) error {
	return r.db.Where("username = ?", username).Delete(&User{}).Error
}

func (r *Repository) UpdateUser(username string, email string, password string) error {
	hashpass, err := HashPassword(password)
	if err != nil {
		return err
	}
	return r.db.Model(&User{}).Where("username = ?", username).Updates(User{Email: email, Password: hashpass}).Error
}

func (r *Repository) GetAllUsers() ([]User, error) {
	var users []User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *Repository) Ping() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func (r *Repository) CheckUserPassword(username string, password string) error {
	user := &User{}
	err := r.db.Where("username = ?", username).First(user).Error
	if err != nil {
		return err
	}
	err = CheckPasswordHash(password, user.Password)
	if err != nil {
		return err
	}
	return nil
}
