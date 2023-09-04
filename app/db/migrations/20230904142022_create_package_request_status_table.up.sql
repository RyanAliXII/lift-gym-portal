CREATE TABLE IF NOT EXISTS package_request_status (
    id INT NOT NULL PRIMARY KEY,
    description varchar(20), 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
INSERT INTO package_request_status(id, description) VALUES
(1, 'Pending'),
(2, "Approved"),
(3, "Received"),
(4, "Cancelled");