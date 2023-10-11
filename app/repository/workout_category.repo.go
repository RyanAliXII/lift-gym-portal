package repository

import (
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/model"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
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
	transaction, err := repo.db.Begin()
	if err != nil {
		transaction.Rollback()
		return err
	}
	result, err := transaction.Exec("INSERT INTO workout_category(name) VALUES(?)", category.Name)

	if err != nil {
		transaction.Rollback()
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		transaction.Rollback()
		return err
	}
	records := make([]goqu.Record, 0)
	for _, workout := range category.Workouts {
		records = append(records, goqu.Record{
			"category_id":  id,
			"workout_id": workout.Id,
		})
	}
	dialect := goqu.Dialect("mysql")
	ds := dialect.Insert(goqu.T("category_workout")).Prepared(true).Rows(records)
	query, args, _:= ds.ToSQL()
	_, err = transaction.Exec(query, args...)
	if err != nil {
		transaction.Rollback()
		return err
	}
	transaction.Commit()
	return err

}
func(repo * WorkoutCategoryRepository) GetCategories() ( []model.WorkoutCategory, error ) {

	categories := make([]model.WorkoutCategory, 0)
	err := repo.db.Select(&categories, `SELECT workout_category.id, workout_category.name, CONCAT('[',GROUP_CONCAT(JSON_OBJECT('id', workout.id, 'name', workout.name, 'description', workout.description, 'imagePath', workout.image_path)), ']') as workouts from workout_category 
	INNER JOIN category_workout on workout_category.id = category_workout.category_id
	INNER JOIN workout on category_workout.workout_id = workout.id
	where workout_category.deleted_at is null GROUP BY workout_category.id order by workout_category.updated_at desc`)
	return categories, err
}
func(repo * WorkoutCategoryRepository) GetCategoryById(id int) ( model.WorkoutCategory, error ) {
	category := model.WorkoutCategory{}
	err := repo.db.Get(&category, `SELECT workout_category.id, workout_category.name, CONCAT('[',GROUP_CONCAT(JSON_OBJECT('id', workout.id, 'name', workout.name,'description',workout.description,'imagePath', workout.image_path)), ']') as workouts from workout_category 
	INNER JOIN category_workout on workout_category.id = category_workout.category_id
	INNER JOIN workout on category_workout.workout_id = workout.id
	where workout_category.deleted_at is null and workout_category.id = ? GROUP BY workout_category.id order by workout_category.updated_at desc LIMIT 1`, id)
	return category, err
}
func(repo * WorkoutCategoryRepository) UpdateCategory(category model.WorkoutCategory) (  error ) {
	transaction, err := repo.db.Begin()
	if err != nil {
		transaction.Rollback()
		return err
	}
	_, err = transaction.Exec("UPDATE workout_category SET name = ? where id = ?", category.Name, category.Id)

	_, err = transaction.Exec("DELETE FROM category_workout where category_id = ?", category.Id)
	if err != nil {
		transaction.Rollback()
		return err
	}
	records := make([]goqu.Record, 0)
	for _, workout := range category.Workouts {
		records = append(records, goqu.Record{
				"category_id":  category.Id,
				"workout_id": workout.Id,
		})
	}
	dialect := goqu.Dialect("mysql")
	ds := dialect.Insert(goqu.T("category_workout")).Prepared(true).Rows(records)
	query, args, _:= ds.ToSQL()
	_, err = transaction.Exec(query, args...)
	if err != nil {
		transaction.Rollback()
		return err
	}
	transaction.Commit()
	return err

}
func(repo * WorkoutCategoryRepository) DeleteCategory(id int) (  error ) {
	_, err := repo.db.Exec("UPDATE workout_category SET deleted_at = now() where id = ?", id)
	return err

}