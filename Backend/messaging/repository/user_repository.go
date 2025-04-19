package repository

import (
	"database/sql"
	"fmt"
	"log"
	"messaging/models"
)

// UserRepository struct will hold the database connection
type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) AddUser(id uint, name string, email string) error {
	// Prepare the SQL statement for inserting a user
	// Using Exec directly is fine for simple inserts like this
	// Ensure your table name is 'users' and columns are 'id', 'name', 'email'
	query := "INSERT INTO users (id, name, email) VALUES ($1, $2, $3)"
	_, err := r.DB.Exec(query, id, name, email)
	if err != nil {
		// Log the error and return a wrapped error
		// The handler will check for specific duplicate key errors (like pq.ErrCode unique_violation)
		log.Printf("Error inserting user (ID: %d, Name: %s) into messaging db: %v", id, name, err)
		return fmt.Errorf("failed to add user %d to messaging database: %w", id, err)
	}

	// Log success (optional)
	// log.Printf("Successfully added user %d (%s) to messaging user table", id, name)
	return nil // Return nil on success
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) UserExists(id uint) (bool, error) {
	var exists bool
	// Use EXISTS for efficiency, it stops scanning once a match is found
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)"
	// Use QueryRowContext for context propagation if needed, otherwise QueryRow is fine
	err := r.DB.QueryRow(query, id).Scan(&exists)
	if err != nil {
		// Don't log sql.ErrNoRows as an error here, QueryRow handles it.
		// EXISTS query always returns a row (true or false).
		log.Printf("Error checking existence for user ID %d: %v", id, err)
		return false, fmt.Errorf("failed to check user existence for ID %d: %w", id, err)
	}
	return exists, nil
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
