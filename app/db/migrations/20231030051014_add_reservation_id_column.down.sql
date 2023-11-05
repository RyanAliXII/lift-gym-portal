ALTER TABLE reservation
DROP COLUMN reservation_id,
DROP CONSTRAINT unique_reservation_id;
