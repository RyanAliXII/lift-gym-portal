package repository

import (
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/model"

	"github.com/jmoiron/sqlx"
)

type WorkoutCategoryRepository struct {
	db *sqlx.DB
}
func NewWorkoutCategoryRepository( ) WorkoutCategoryRepository {
	return WorkoutCategoryRepository{
		db: db.GetConnection(),
	}
}
func(repo * WorkoutCategoryRepository) NewCategory(category model.WorkoutCategory) error {
	_, err := repo.db.Exec("INSERT INTO workout_category(name) VALUES(?)", category.Name)
	return err

}
func(repo * WorkoutCategoryRepository) GetCategories() ( []model.WorkoutCategory, error ) {
	categories := make([]model.WorkoutCategory, 0)
	err := repo.db.Select(&categories, "SELECT id, name from workout_category where deleted_at is null")
	return categories, err

}