CREATE TABLE IF NOT EXISTS hired_coaches_status (
    id INT NOT NULL PRIMARY KEY,
    description varchar(20), 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

INSERT INTO hired_coaches_status (id, description) VALUES
(1, 'Pending'),
(2, "Approved"),
(3, "Paid"),
(4, "Cancelled");