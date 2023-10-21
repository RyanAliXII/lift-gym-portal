-- Add the column with a default value of NULL
ALTER TABLE coaching_rate;
ADD COLUMN deleted_at DEFAULT NULL;


