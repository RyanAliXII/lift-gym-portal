ALTER TABLE package_request
DROP FOREIGN KEY package_request_ibfk_4,
DROP COLUMN package_snapshot_id;


ALTER TABLE subscription
DROP FOREIGN KEY subscription_ibfk_3,
DROP COLUMN membership_plan_snapshot_id;