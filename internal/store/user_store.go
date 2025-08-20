package store

import (
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// password represents a user’s password.
// It keeps both the hash (always stored) and optionally
// the plain-text (for temporary use during creation/validation).
// - lowercase `password` makes it unexported (cannot be accessed outside this package).
type password struct {
	plainText *string // pointer so we can distinguish between "empty string" vs "not set"
	hash      []byte  // securely hashed password
}

// Set hashes a plain-text password and stores both the hash and (optionally) the plain-text.
// The hash is generated using bcrypt with cost=12 (secure default).
func (p *password) Set(plainTextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), 12)
	if err != nil {
		return err
	}

	// Keep plain-text temporarily (might be used in validation before discard).
	p.plainText = &plainTextPassword
	p.hash = hash
	return nil
}

// Matches checks if a given plain-text password matches the stored hash.
// Returns true if correct, false if incorrect, or error if something else went wrong.
func (p *password) Matches(plainTextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainTextPassword))
	if err != nil {
		switch {
		// User entered the wrong password.
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		// Some other error occurred.
		default:
			return false, err
		}
	}
	return true, nil
}

// User represents the user model in the system.
// JSON tags control API responses, while db operations are handled in queries.
// Note: PasswordHash is excluded from JSON with `json:"-"`.
type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash password  `json:"-"`
	Bio          string    `json:"bio"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// PostgresUserStore implements UserStore using PostgreSQL as the backend.
// It wraps a sql.DB connection.
type PostgresUserStore struct {
	db *sql.DB
}

// NewPostgresUserStore is a constructor for PostgresUserStore.
func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{
		db: db,
	}
}

// UserStore defines an abstraction (interface) for user persistence.
// This allows swapping out Postgres for another backend (MySQL, mock for tests, etc.)
type UserStore interface {
	CreateUser(*User) error
	GetUserByUsername(username string) (*User, error)
	UpdateUser(*User) error
}

// CreateUser inserts a new user into the database.
// Password is stored as a bcrypt hash, NOT plain-text.
// Returns error if insertion fails.
func (s *PostgresUserStore) CreateUser(user *User) error {
	query := `
	INSERT INTO users (username, email, password_hash, bio)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, updated_at
	`

	// Use QueryRow + Scan to capture the generated fields.
	err := s.db.QueryRow(query,
		user.Username,
		user.Email,
		user.PasswordHash.hash, // store only the hash, never the plain text
		user.Bio,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return err
	}
	return nil
}

// GetUserByUsername fetches a user by username.
// Returns (*User, nil) if found, (nil, nil) if not found, or (nil, error) if query fails.
func (s *PostgresUserStore) GetUserByUsername(username string) (*User, error) {
	user := &User{
		PasswordHash: password{},
	}
	query := `
	SELECT id, username, email, password_hash, bio, created_at, updated_at
	FROM users
	WHERE username = $1
	`

	err := s.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash.hash, // hydrate password hash for login checks
		&user.Bio,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	// Handle "no rows" as a valid non-error result
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser updates basic user fields in the database.
// Updates: username, email, bio.
// updated_at is set to CURRENT_TIMESTAMP automatically.
// Returns sql.ErrNoRows if user ID does not exist.
func (s *PostgresUserStore) UpdateUser(user *User) error {
	query := `
	UPDATE users
	SET username = $1, email = $2, bio = $3, updated_at = CURRENT_TIMESTAMP
	WHERE id = $4
	RETURNING updated_at
	`

	result, err := s.db.Exec(query, user.Username, user.Email, user.Bio, user.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// If no rows updated → user not found
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
