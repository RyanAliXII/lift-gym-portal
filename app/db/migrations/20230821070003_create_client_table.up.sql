CREATE TABLE IF NOT EXISTS client (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    given_name VARCHAR(255) NOT NULL,
    middle_name VARCHAR(255) NOT NULL,
    surname VARCHAR(255) NOT NULL,
    date_of_birth DATE NOT NULL,
    mobile_number VARCHAR(30) NOT NULL,
    emergency_contact VARCHAR(30) NOT NULL,
    address TEXT NOT NULL,
    account_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (account_id) REFERENCES account(id)
)