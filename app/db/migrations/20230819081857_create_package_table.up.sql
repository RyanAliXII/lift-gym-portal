CREATE TABLE IF NOT EXISTS `package` (
    id INT NOT NULL AUTO_INCREMENT  PRIMARY KEY,
    description TEXT NOT NULL, 
    price DECIMAL(13, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
)