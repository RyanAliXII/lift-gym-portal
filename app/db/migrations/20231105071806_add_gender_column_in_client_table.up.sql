ALTER TABLE client
ADD COLUMN gender ENUM('male', 'female', 'other','prefer not to answer', '') DEFAULT '';