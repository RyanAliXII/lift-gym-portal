CREATE TABLE IF NOT EXISTS hired_coach(
   id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
   coach_id int not null,
   rate_id int not null,
   rate_snapshot_id int not null,
   client_id int not null,
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   FOREIGN KEY (coach_id) REFERENCES coach(id),
   FOREIGN KEY (rate_id) REFERENCES coaching_rate(id),
   FOREIGN KEY (client_id) REFERENCES client(id),
   FOREIGN KEY (rate_snapshot_id) REFERENCES coaching_rate_snapshot(id)
)