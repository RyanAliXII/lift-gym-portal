CREATE TABLE IF NOT EXISTS coach_image (
    id INT NOT NULL AUTO_INCREMENT  PRIMARY KEY,
    coach_id INT NOT NULL,
    path text,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (coach_id) REFERENCES coach(id)
)