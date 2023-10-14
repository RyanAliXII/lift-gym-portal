-- Add the column with a default value of NULL
ALTER TABLE user
ADD COLUMN role_id INT DEFAULT NULL;

-- Add the foreign key constraint
ALTER TABLE user
ADD CONSTRAINT fk_user_role
FOREIGN KEY (role_id)
REFERENCES role(id);
