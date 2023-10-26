CREATE TABLE IF NOT EXISTS password_reset(
    id INT NOT NULL AUTO_INCREMENT  PRIMARY KEY,
    public_key VARCHAR(50) NOT NULL,
    account_id INT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP NULL,
    FOREIGN KEY (account_id) REFERENCES account(id)
)