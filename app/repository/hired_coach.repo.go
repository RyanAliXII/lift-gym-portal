package repository

import (
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/pkg/status"

	"github.com/jmoiron/sqlx"
)

type HiredCoachRepository struct {
	db *sqlx.DB
}
func NewHiredCoachRepository() (HiredCoachRepository){
	return HiredCoachRepository{
		db: db.GetConnection(),
	}
}
func (repo * HiredCoachRepository) Hire(hiredCoach model.HiredCoach) error {
	transaction, err := repo.db.Beginx()
	if err != nil {
		transaction.Rollback()
		return err
	}
	coachRate := model.CoachRate{}
	err = transaction.Get(&coachRate, "SELECT id, description, price, coach_id from coaching_rate where id = ? and coach_id = ? and deleted_at is null LIMIT 1", hiredCoach.RateId, hiredCoach.CoachId)
	if err != nil {
		transaction.Rollback()
		return err
	}
	result, err := transaction.Exec("INSERT INTO coaching_rate_snapshot (description, price, coach_id) VALUES(?, ?, ?)", coachRate.Description, coachRate.Price, coachRate.CoachId)
	if err != nil {
		transaction.Rollback()
		return err
	}
	snapshotId, err := result.LastInsertId()
	if err != nil {
		transaction.Rollback()
		return err
	}
	_, err = transaction.Exec("INSERT INTO hired_coach(coach_id, rate_id, rate_snapshot_id, client_id, remarks, meeting_time) VALUES(?, ?, ?, ?, ?,?)", coachRate.CoachId, coachRate.Id, snapshotId, hiredCoach.ClientId, hiredCoach.Remarks, hiredCoach.MeetingTime)
	if err != nil {
		transaction.Rollback()
		return err
	}
	transaction.Commit()
	return nil
}

func (repo * HiredCoachRepository)GetCoachReservationByClientId(clientId int )([]model.HiredCoach, error){
	query := `
	SELECT 
	hired_coach.id, 
	hired_coach.coach_id,
	hired_coach.rate_id, 
	hired_coach.rate_snapshot_id,
	hired_coach.client_id,
	COALESCE(hired_coach.meeting_time, "") as meeting_time,
	remarks,
	hired_coach.status_id,
	hcs.description as status,
	JSON_OBJECT(
		'id', hired_coach.coach_id,
		'givenName', coach.given_name,
		'middleName', coach.middle_name,
		'surname', coach.surname,
		'email', account.email,
		'mobileNumber', coach.mobile_number
	) as coach,
	JSON_OBJECT('id', coaching_rate.id, 'description', coaching_rate.description, 'price', coaching_rate.price, 'coachId', coaching_rate.coach_id) as rate,
	JSON_OBJECT('id', coaching_rate_snapshot.id, 'description', coaching_rate_snapshot.description, 'price', coaching_rate_snapshot.price) as rate_snapshot
	FROM hired_coach
	INNER JOIN coach ON hired_coach.coach_id = coach.id
	INNER JOIN account ON coach.account_id = account.id
	INNER JOIN coaching_rate ON hired_coach.rate_id = coaching_rate.id
	INNER JOIN coaching_rate_snapshot ON hired_coach.rate_snapshot_id = coaching_rate_snapshot.id
	INNER JOIN hired_coaches_status as hcs on hired_coach.status_id = hcs.id
	where client_id = ?
	ORDER BY hired_coach.created_at desc;
	`
    hiredCoaches := make([]model.HiredCoach, 0)
	err := repo.db.Select(&hiredCoaches, query, clientId)
	if err != nil {
		return hiredCoaches, err
	}
	return hiredCoaches, nil
}

func (repo * HiredCoachRepository)GetCoachAppointments(coachId int )([]model.HiredCoach, error){
	query := `
		SELECT 
		hired_coach.id, 
		hired_coach.coach_id,
		hired_coach.rate_id, 
		hired_coach.rate_snapshot_id,
		hired_coach.client_id,
		hired_coach.status_id,
		hcs.description as status,
		COALESCE(hired_coach.meeting_time, "") as meeting_time,
		remarks,
		JSON_OBJECT(
			'id', client.id,
			'givenName', client.given_name,
			'middleName', client.middle_name,
			'surname', client.surname,
			'email', account.email,
			'mobileNumber', client.mobile_number
		) as client,
		JSON_OBJECT('id', coaching_rate.id, 'description', coaching_rate.description, 'price', coaching_rate.price, 'coachId', coaching_rate.coach_id) as rate,
		JSON_OBJECT('id', coaching_rate_snapshot.id, 'description', coaching_rate_snapshot.description, 'price', coaching_rate_snapshot.price) as rate_snapshot
		FROM hired_coach
		INNER JOIN client ON hired_coach.client_id = client.id
		INNER JOIN account ON client.account_id = account.id
		INNER JOIN coaching_rate ON hired_coach.rate_id = coaching_rate.id
		INNER JOIN coaching_rate_snapshot ON hired_coach.rate_snapshot_id = coaching_rate_snapshot.id
		INNER JOIN hired_coaches_status as hcs on hired_coach.status_id = hcs.id
		where hired_coach.coach_id = ?
		ORDER BY hired_coach.created_at desc
	`
    hiredCoaches := make([]model.HiredCoach, 0)
	err := repo.db.Select(&hiredCoaches, query, coachId)
	if err != nil {
		return hiredCoaches, err
	}
	return hiredCoaches, nil
}

func(repo * HiredCoachRepository)CancelAppointmentByClient(appointment model.HiredCoach) error {
	_, err := repo.db.Exec("UPDATE hired_coach SET status_id = ?, remarks = ? where id = ? and client_id = ? and (status_id = ? OR status_id = ?)", status.CoachAppointmentStatusCancelled, appointment.Remarks, appointment.Id, appointment.ClientId, status.CoachAppointmentStatusPending, status.CoachAppointmentStatusApproved)
	return err
}

func(repo * HiredCoachRepository)MarkAppointmentAsApproved(appointment model.HiredCoach) error {
	_, err := repo.db.Exec("UPDATE hired_coach SET status_id = ? where id = ? and coach_id = ? and status_id = ?", status.CoachAppointmentStatusApproved,appointment.Id, appointment.CoachId, status.CoachAppointmentStatusPending)
	return err
}
func(repo * HiredCoachRepository)MarkAppointmentAsPaid(appointment model.HiredCoach) error {
	_, err := repo.db.Exec("UPDATE hired_coach SET status_id = ? where id = ? and coach_id = ? and status_id = ?", status.CoachAppointmentStatusPaid,appointment.Id, appointment.CoachId, status.CoachAppointmentStatusApproved)
	return err
}

func(repo * HiredCoachRepository)CancelAppointment(appointment model.HiredCoach) error {
	_, err := repo.db.Exec("UPDATE hired_coach SET status_id = ?, remarks = ? where id = ? and coach_id = ?  and (status_id = 1 or status_id = 2)", status.CoachAppointmentStatusCancelled, appointment.Remarks, appointment.Id, appointment.CoachId)
	return err
}
func(repo * HiredCoachRepository)MarkAsNoShow(appointment model.HiredCoach) error {
	dbAppointment, err := repo.GetAppointmentById(appointment.Id)
	if err != nil{
		return err
	}	
	_, err = repo.db.Exec("UPDATE hired_coach SET status_id = ? where id = ? and coach_id = ? and status_id = 2", status.CoachAppointmentStatusNoShow, appointment.Id, appointment.CoachId)

	if err != nil {
		return err
	}
	err = repo.handlePenalty(dbAppointment)
	if err  != nil {
		return err
	}
	return nil
}
func (repo * HiredCoachRepository)getLastThreeAppointmentOfClient(clientId int ) ([]model.HiredCoach, error) {
	appointments := make([]model.HiredCoach, 0)
	err := repo.db.Select(&appointments, "SELECT status_id FROM hired_coach where client_id = ? order by created_at desc LIMIT 3", clientId)
	return appointments, err
} 

func (repo *HiredCoachRepository)GetAppointmentById(id int )(model.HiredCoach, error){
	appointment := model.HiredCoach{}
	query := `
	SELECT 
		hired_coach.id, 
		hired_coach.coach_id,
		hired_coach.rate_id, 
		hired_coach.rate_snapshot_id,
		hired_coach.client_id,
		hired_coach.status_id,
		JSON_OBJECT('id', coaching_rate_snapshot.id, 'description', coaching_rate_snapshot.description, 'price', coaching_rate_snapshot.price) as rate_snapshot
		from hired_coach
		INNER JOIN coaching_rate_snapshot ON hired_coach.rate_snapshot_id = coaching_rate_snapshot.id
		where hired_coach.id = ? LIMIT 1
	`
	err := repo.db.Get(&appointment, query, id)
	return appointment, err
}

func (repo *HiredCoachRepository)handlePenalty(appointment model.HiredCoach) error{

	// appointments, err := repo.getLastThreeAppointmentOfClient(appointment.ClientId)
	// if err != nil { 
	// 	return err
	// }
	// MaxMissedAppointments := 3
	// missed := 0
	// for _, a := range appointments {
	// 	if (a.StatusId == status.CoachAppointmentStatusNoShow){
	// 		missed += 1
	// 	}
	// }
	// if missed < MaxMissedAppointments {
	// 	return nil
	// }
	// if (repo.HasPenalty(appointment.ClientId)){
	// 	return nil
	// }
	_, err := repo.db.Exec("INSERT INTO coach_appointment_penalty(amount, client_id, coach_id) VALUES (?,?,?)", appointment.RateSnapshot.Price, appointment.ClientId, appointment.CoachId)
	return err
}

func (repo *HiredCoachRepository) HasPenalty(clientId int) bool{
	recordCount := 1
	repo.db.Get(&recordCount, "SELECT count(1) as recordCount from coach_appointment_penalty where client_id = ? and settled_at is null", clientId)
	return recordCount >= 1 
}	