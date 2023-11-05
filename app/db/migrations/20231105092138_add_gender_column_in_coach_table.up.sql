ALTER TABLE coach
ADD COLUMN gender ENUM('male', 'female', 'other','prefer not to answer', '') DEFAULT '';