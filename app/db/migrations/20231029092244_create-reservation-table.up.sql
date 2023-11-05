CREATE TABLE IF NOT EXISTS reservation (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    client_id INT NOT NULL,
    date_slot_id INT NOT NULL,
    time_slot_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    cancelled_at TIMESTAMP NULL,
    attended_at TIMESTAMP NULL,
    FOREIGN KEY (date_slot_id) REFERENCES date_slot(id),
    FOREIGN KEY (time_slot_id) REFERENCES time_slot(id),
    FOREIGN KEY (client_id) REFERENCES client(id)
);