package repository

import (
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/model"

	"github.com/jmoiron/sqlx"
)

type WorkoutRepository struct {
	db *sqlx.DB
}

func (repo * WorkoutRepository)NewWorkout(workout model.Workout) error {
	_, err := repo.db.Exec("INSERT INTO workout(name, description, image_path)VALUES(?,?,?)", workout.Name, workout.Description, workout.ImagePath)
	return err
}
func NewWorkoutRepository() WorkoutRepository {
	return WorkoutRepository{
		db: db.GetConnection(),
	}
}