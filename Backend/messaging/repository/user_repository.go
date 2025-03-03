package repository

import (
	"database/sql"
	"log"
	"messaging/models"
)

// UserRepository struct will hold the database connection
type UserRepository struct {
	DB *sql.DB
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// GetAllUsers retrieves all users from the database
func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	rows, err := r.DB.Query("SELECT id, name, email FROM users")
	if err != nil {
		log.Println("Error fetching users:", err)
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			log.Println("Error scanning user:", err)
			return nil, err
		}
		users = append(users, user)
	}

	// Check for any error that occurred during iteration
	if err := rows.Err(); err != nil {
		log.Println("Error during rows iteration:", err)
		return nil, err
	}

	return users, nil
}
