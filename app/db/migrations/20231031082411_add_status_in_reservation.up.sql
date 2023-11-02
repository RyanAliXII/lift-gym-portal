ALTER TABLE reservation
ADD COLUMN status_id INT DEFAULT 1;

ALTER TABLE reservation
ADD CONSTRAINT fk_reservation_status
FOREIGN KEY (status_id)
REFERENCES reservation_status (id);
