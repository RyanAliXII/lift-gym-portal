CREATE TABLE IF NOT EXISTS category_workout(
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    category_id INT,
    workout_id INT,
    FOREIGN KEY (category_id) REFERENCES workout_category(id),
    FOREIGN KEY (workout_id) REFERENCES workout(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)