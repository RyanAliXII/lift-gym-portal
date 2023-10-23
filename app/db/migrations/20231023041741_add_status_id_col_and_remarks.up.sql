ALTER TABLE hired_coach
ADD COLUMN status_id INT NOT NULL DEFAULT 1,
ADD COLUMN remarks TEXT DEFAULT '';

ALTER TABLE hired_coach
ADD FOREIGN KEY (status_id)
REFERENCES hired_coaches_status(id);