CREATE TABLE IF NOT EXISTS subscription(
    id INT NOT NULL AUTO_INCREMENT  PRIMARY KEY,
    client_id INT NOT NULL,
    membership_plan_id INT NOT NULL,
    valid_until DATE NOT NULL,
    cancelled_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (client_id) REFERENCES client(id),
    FOREIGN KEY (membership_plan_id) REFERENCES membership_plan(id)
)