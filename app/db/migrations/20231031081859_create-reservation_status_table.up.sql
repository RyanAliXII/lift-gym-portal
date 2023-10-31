CREATE TABLE IF NOT EXISTS reservation_status (
    id INT NOT NULL PRIMARY KEY,
    description varchar(20), 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
INSERT INTO reservation_status(id, description) VALUES
(1, 'Default'),
(2, "Attended"),
(3, "No-Show"),
(4, "Cancelled");