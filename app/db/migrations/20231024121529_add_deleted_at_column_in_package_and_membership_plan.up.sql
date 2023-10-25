ALTER TABLE package
ADD COLUMN deleted_at timestamp null;


ALTER TABLE membership_plan
ADD COLUMN deleted_at timestamp null;