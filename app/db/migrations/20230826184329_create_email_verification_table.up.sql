CREATE TABLE IF NOT EXISTS email_verification (
    id INT NOT NULL AUTO_INCREMENT  PRIMARY KEY,
    public_id varchar(50) NOT NULL,
    client_id INT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)