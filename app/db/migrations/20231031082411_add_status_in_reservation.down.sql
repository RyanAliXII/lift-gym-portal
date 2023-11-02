ALTER TABLE reservation
DROP CONSTRAINT fk_reservation_status,
DROP COLUMN status_id;
