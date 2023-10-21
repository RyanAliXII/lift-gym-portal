CREATE TABLE IF NOT EXISTS coaching_rate_snapshot(
   id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
   coach_id int not null,
   description VARCHAR(255) DEFAULT '',
   price DECIMAL(13, 2) default 0,
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   FOREIGN KEY (coach_id) REFERENCES coach(id)
)