package repository

import (
	"context"
	"log"

	"first-project/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	SetupUserTable()
	CreateUser(user models.User) (int, error)
	GetUserByEmail(email string) (models.User, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) SetupUserTable() {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(100) UNIQUE NOT NULL,
		password_hash VARCHAR(100) NOT NULL
	);`
	_, err := r.db.Exec(context.Background(), createTableQuery)
	if err != nil {
		log.Fatalf("Gagal membuat tabel 'users': %v\n", err)
	}
}

func (r *userRepository) CreateUser(user models.User) (int, error) {
	query := `INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id`
	var newID int

	err := r.db.QueryRow(
		context.Background(),
		query,
		user.Email,
		user.Password,
	).Scan(&newID)

	return newID, err
}

func (r *userRepository) GetUserByEmail(email string) (models.User, error) {
	query := `SELECT id, email, password_hash FROM users WHERE email = $1`
	var user models.User

	err := r.db.QueryRow(
		context.Background(),
		query,
		email,
	).Scan(&user.ID, &user.Email, &user.Password)

	return user, err
}
