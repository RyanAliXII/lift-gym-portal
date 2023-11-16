CREATE TABLE IF NOT EXISTS coach_log (
   id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
   coach_id int not null,
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   deleted_at TIMESTAMP NULL,
   logged_out_at timestamp null,
   FOREIGN KEY (coach_id) REFERENCES coach(id)
)