ALTER TABLE reservation
ADD COLUMN reservation_status_id INT DEFAULT 1;

ALTER TABLE reservation
ADD CONSTRAINT fk_reservation_status
FOREIGN KEY (reservation_status_id)
REFERENCES reservation_status (id);
