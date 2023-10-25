ALTER TABLE hired_coach
ADD COLUMN status_id INT NOT NULL DEFAULT 1,
ADD COLUMN remarks TEXT NOT NULL;

ALTER TABLE hired_coach
ADD FOREIGN KEY (status_id)
REFERENCES hired_coaches_status(id);