ALTER TABLE coach
ADD COLUMN public_id varchar(40) GENERATED ALWAYS AS  
(CONCAT('CH',
year(coach.created_at),
month(coach.created_at),
day(coach.created_at),
hour(coach.created_at),
minute(coach.created_at),
second(coach.created_at)
)) STORED;

ALTER TABLE coach
ADD CONSTRAINT unique_coach_public_id UNIQUE (public_id);