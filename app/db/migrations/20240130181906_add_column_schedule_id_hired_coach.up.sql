ALTER TABLE hired_coach
ADD COLUMN schedule_id INT,
ADD FOREIGN KEY (schedule_id) REFERENCES coach_schedule(id);