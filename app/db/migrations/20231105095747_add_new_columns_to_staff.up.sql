ALTER TABLE user
ADD COLUMN date_of_birth DATE,
ADD COLUMN mobile_number VARCHAR(30) DEFAULT '',
ADD COLUMN emergency_contact VARCHAR(30) DEFAULT '',
ADD COLUMN address varchar(255) DEFAULT '',
ADD COLUMN gender ENUM('male', 'female', 'other','prefer not to answer', '') DEFAULT '',
ADD COLUMN public_id varchar(40) GENERATED ALWAYS AS  
(CONCAT('ST',
year(user.created_at),
month(user.created_at),
day(user.created_at),
hour(user.created_at),
minute(user.created_at),
second(user.created_at)
)) STORED;