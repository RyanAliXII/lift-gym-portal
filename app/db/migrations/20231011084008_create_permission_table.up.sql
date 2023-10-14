CREATE TABLE IF NOT EXISTS permission (
     id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
     value VARCHAR(100) DEFAULT '',
     role_id INT NOT NULL,
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
     FOREIGN KEY (role_id) REFERENCES role(id)
)