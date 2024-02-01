CREATE TABLE general_inventory(
 id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
 name VARCHAR(100) DEFAULT '',
 brand VARCHAR(100) DEFAULT '',
 quantity INT DEFAULT 0,
 unit_of_measure ENUM('grams', 'kilograms','pack','pounds','pieces', 'millilitres', 'litres') NOT NULL,
 cost_price DECIMAL(13, 2) NOT NULL,
 date_received DATE NOT NULL,
 quantity_threshold INT DEFAULT 0,
 expiration_date DATE,
 image TEXT,
 deleted_at TIMESTAMP NULL,
 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
 updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
)