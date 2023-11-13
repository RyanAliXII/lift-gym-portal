CREATE TABLE coach_appointment_penalty (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    amount DECIMAL(13, 2) NOT NULL,
    client_id INT NOT NULL,
    coach_id INT NOT NULL,
    settled_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (client_id) REFERENCES client(id),
    FOREIGN KEY (coach_id) REFERENCES coach(id)
)