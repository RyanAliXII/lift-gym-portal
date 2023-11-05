ALTER TABLE client
ADD COLUMN public_id varchar(40) GENERATED ALWAYS AS  
(CONCAT('C',
year(client.created_at),
month(client.created_at),
day(client.created_at),
hour(client.created_at),
minute(client.created_at),
second(client.created_at)
)) STORED;

ALTER TABLE client
ADD CONSTRAINT unique_client_public_id UNIQUE (public_id);