ALTER TABLE user 
ADD column avatar varchar(100) NOT NULL DEFAULT '';

ALTER TABLE coach
ADD column avatar varchar(100) NOT NULL DEFAULT '';

ALTER TABLE client
ADD column avatar varchar(100) NOT NULL DEFAULT '';