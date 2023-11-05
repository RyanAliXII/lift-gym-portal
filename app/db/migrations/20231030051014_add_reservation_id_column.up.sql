ALTER TABLE reservation
ADD COLUMN reservation_id varchar(40) GENERATED ALWAYS AS  
(CONCAT('R',
year(reservation.created_at),
month(reservation.created_at),
day(reservation.created_at),
hour(reservation.created_at),
minute(reservation.created_at),
second(reservation.created_at)
)) STORED; 

ALTER TABLE reservation
ADD CONSTRAINT unique_reservation_id UNIQUE (reservation_id);