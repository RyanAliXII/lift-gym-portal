ALTER TABLE package_request
ADD COLUMN package_snapshot_id INT;

ALTER TABLE package_request
ADD FOREIGN KEY (package_snapshot_id) REFERENCES package_snapshot(id);

ALTER TABLE subscription
ADD COLUMN membership_plan_snapshot_id INT;

ALTER TABLE subscription
ADD FOREIGN KEY (membership_plan_snapshot_id) REFERENCES membership_plan_snapshot(id);


ALTER TABLE membership_request
ADD COLUMN membership_plan_snapshot_id INT;

ALTER TABLE membership_request
ADD FOREIGN KEY (membership_plan_snapshot_id) REFERENCES membership_plan_snapshot(id);