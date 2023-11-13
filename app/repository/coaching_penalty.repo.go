package repository

import (
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/model"

	"github.com/jmoiron/sqlx"
)

type CoachingPenalty struct{
	db * sqlx.DB
}
func NewCoachingPenalty() CoachingPenalty{
	return CoachingPenalty{
		db: db.GetConnection(),
	}
}


func (repo * CoachingPenalty) GetPenalties ()([]model.CoachAppointmentPenalty, error) {
	penalties := make([]model.CoachAppointmentPenalty,0)
	query := `
	SELECT
	coach_appointment_penalty.id, 
	client_id, 
	coach_id, 
	amount,
	(case when settled_at is null then false else true end) as is_settled,  
	JSON_OBJECT(
	   'publicId', client.public_id, 
	   'id', client.id, 
	   'givenName', client.given_name, 
	   'middleName', client.middle_name, 
	   'surname', client.surname)  as client,
	JSON_OBJECT(
	   'id', coach.id,
	   'givenName', coach.given_name,
	   'middleName', coach.middle_name,
	   'surname', coach.surname
	) as coach
	FROM coach_appointment_penalty
	INNER JOIN client on coach_appointment_penalty.client_id = client.id
	INNER JOIN coach on coach_appointment_penalty.coach_id = coach.id
	order by coach_appointment_penalty.created_at desc
	`
	err := repo.db.Select(&penalties, query)
	return penalties, err
}

func (repo * CoachingPenalty) MarkAsSettled (id int)error{
	_,err := repo.db.Exec("UPDATE coach_appointment_penalty set settled_at = now() where id = ?", id)
	return err
}
func (repo * CoachingPenalty) MarkAsUnSettled (id int)error{
	_,err := repo.db.Exec("UPDATE coach_appointment_penalty set settled_at = null where id = ?", id)
	return err
}