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
}

// CreateWorkout inserts a Workout and its associated WorkoutEntries into the database.
// - Uses a transaction to ensure all-or-nothing behavior.
// - Uses $1, $2... placeholders to prevent SQL injection.
// - Populates the IDs of the workout and its entries after insertion.
// Returns a pointer to the inserted Workout with all IDs populated, or an error if any occurs.
func (pg *PostgresWorkoutStore) CreateWorkout(workout *Workout) (*Workout, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	//The doller sign is used to prevent sql injection, we inject this query into the queryrow and these then prevents them
	query := `
		INSERT INTO workouts (user_id, title, description, duration_minutes, calories_burned)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id 
		`

	err = tx.QueryRow(query, workout.UserID, workout.Title, workout.Description, workout.DurationMinutes, workout.CaloriesBurned).Scan(&workout.ID)
	//We scan to get the id
	if err != nil {
		return nil, err
	}

	//We also need to insert the entries
	for _, entry := range workout.Entries {
		query := `
			INSERT INTO workout_entries (workout_id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id
    			`
		err = tx.QueryRow(query, workout.ID, entry.ExerciseName, entry.Sets, entry.Reps, entry.DurationSeconds, entry.Weight, entry.Notes, entry.OrderIndex).Scan(&entry.ID)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return workout, nil
}

func (pg *PostgresWorkoutStore) GetWorkoutById(id int64) (*Workout, error){
	workout := &Workout{}
	return workout, nil
}
