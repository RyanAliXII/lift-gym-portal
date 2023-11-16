CREATE TABLE IF NOT EXISTS staff_log (
   id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
   staff_id int not null,
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   deleted_at TIMESTAMP NULL,
   logged_out_at timestamp null,
   FOREIGN KEY (staff_id) REFERENCES user(id)
)