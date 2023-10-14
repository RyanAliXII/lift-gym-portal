CREATE TABLE IF NOT EXISTS client_log (
   id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
   client_id int not null,
   is_member BOOLEAN default false,
   amount_paid DECIMAL(13, 2) default 0,
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   deleted_at TIMESTAMP NULL
)