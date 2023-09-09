package repository

import (
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/model"

	"github.com/jmoiron/sqlx"
)

type CoachImageRepository struct {
	db *sqlx.DB
}
func(repo * CoachImageRepository)GetImagesByCoachId(coachId int)([]model.CoachImage, error){
	images := make([]model.CoachImage, 0)
	repo.db.Select(&images, "Select id, coach_id, path from coach_image where coach_id = ? order by created_at desc", coachId)
	return images, nil
}

func(repo * CoachImageRepository)NewCoachImages(images []model.CoachImage) (error){
	_, err := repo.db.NamedExec("INSERT INTO coach_image(coach_id, path)VALUES(:coach_id, :path)", images)
	return err
}

func NewCoachImageRepository() CoachImageRepository {

	return CoachImageRepository{
		db: db.GetConnection(),
	}
}

