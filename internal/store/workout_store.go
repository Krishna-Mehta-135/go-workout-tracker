package store

import "database/sql"

type Workout struct {
	ID              int            `json:"id"`
	UserID          int            `json:"user_id"`
	Title           string         `json:"title"`
	Description     string         `json:"description"`
	DurationMinutes int            `json:"duration_minutes"`
	CaloriesBurned  int            `json:"calories_burned"`
	Entries         []WorkoutEntry `json:"entries"`
}

// We used pointer because we explicitly wanted to check if the value is nil or not
type WorkoutEntry struct {
	ID              int      `json:"id"`
	ExerciseName    string   `json:"exercise_name"`
	Sets            int      `json:"sets"`
	Reps            *int     `json:"reps"`
	DurationSeconds *int     `json:"duration_seconds"`
	Weight          *float64 `json:"weight"`
	Notes           string   `json:"notes"`
	OrderIndex      int      `json:"order_index"`
}

// PostgresWorkoutStore is a store struct that encapsulates a Postgres database connection.
// All workout-related database operations will be attached to this struct.
type PostgresWorkoutStore struct {
	db *sql.DB // Database connection
}

// NewPostgresWorkoutStore initializes a new PostgresWorkoutStore with a given *sql.DB connection.
// This allows injecting the database connection and keeps database logic encapsulated.
func NewPostgresWorkoutStore(db *sql.DB) *PostgresWorkoutStore {
	return &PostgresWorkoutStore{db: db}
}

type WorkoutStore interface {
	CreateWorkout(*Workout) (*Workout, error)
	GetWorkoutById(id int64) (*Workout, error)
	UpdateWorkout(*Workout) error
}


// CreateWorkout inserts a new workout along with its entries into the database.
func (pg *PostgresWorkoutStore) CreateWorkout(workout *Workout) (*Workout, error) {
	// Start a new transaction. This ensures that either both the workout and all its entries are saved,
	// or none of them are saved if an error occurs.
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() // If the function exits before tx.Commit(), the transaction will roll back.

	// Insert the workout into the 'workouts' table.
	// $1, $2... are placeholders to safely inject parameters and prevent SQL injection.
	query := `
		INSERT INTO workouts (user_id, title, description, duration_minutes, calories_burned)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id 
	`
	// Execute the query and scan the generated ID back into workout.ID
	err = tx.QueryRow(query, workout.UserID, workout.Title, workout.Description, workout.DurationMinutes, workout.CaloriesBurned).Scan(&workout.ID)
	if err != nil {
		return nil, err
	}

	// Insert each WorkoutEntry into the 'workout_entries' table.
	for _, entry := range workout.Entries {
		entryQuery := `
			INSERT INTO workout_entries (workout_id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id
    	`
		// Scan the generated entry ID into entry.ID
		err = tx.QueryRow(entryQuery, workout.ID, entry.ExerciseName, entry.Sets, entry.Reps, entry.DurationSeconds, entry.Weight, entry.Notes, entry.OrderIndex).Scan(&entry.ID)
		if err != nil {
			return nil, err
		}
	}

	// Commit the transaction. If this succeeds, all inserts are permanently saved.
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// Return the workout struct including its ID and entries with IDs populated.
	return workout, nil
}


// GetWorkoutById retrieves a workout and its associated entries from the database by workout ID.
func (pg *PostgresWorkoutStore) GetWorkoutById(id int64) (*Workout, error) {
	workout := &Workout{}

	// Query the workouts table for the basic workout information
	query := `
	SELECT id, title, description, duration_minutes, calories_burned
	FROM workouts
	WHERE id = $1
	`
	err := pg.db.QueryRow(query, id).Scan(&workout.ID, &workout.Title, &workout.Description, &workout.DurationMinutes, &workout.CaloriesBurned)

	if err == sql.ErrNoRows {
		// Return nil if no workout is found
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Query all associated entries for the workout
	entryQuery := `
	SELECT id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index
	FROM workout_entries
	WHERE workout_id = $1
	ORDER BY order_index
	`

	rows, err := pg.db.Query(entryQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Close the rows when done

	// Loop through each row and append to workout.Entries
	for rows.Next() {
		var entry WorkoutEntry
		err = rows.Scan(
			&entry.ID,
			&entry.ExerciseName,
			&entry.Sets,
			&entry.Reps,
			&entry.DurationSeconds,
			&entry.Weight,
			&entry.Notes,
			&entry.OrderIndex,
		)
		if err != nil {
			return nil, err
		}
		workout.Entries = append(workout.Entries, entry)
	}

	// Return the complete workout struct with its entries
	return workout, nil
}

