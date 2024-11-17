CREATE SEQUENCE IF NOT EXISTS task_id START 1 MAXVALUE 2147483647;
CREATE TABLE IF NOT EXISTS task (
	id INTEGER NOT NULL DEFAULT nextval('task_id'),
	title VARCHAR(256),
	location_name VARCHAR(256),
	location_address VARCHAR(256),
	description TEXT,
	lat DOUBLE PRECISION,
	lng DOUBLE PRECISION,
	uploaded TIMESTAMP,
	start TIMESTAMP,
	stop TIMESTAMP,
	initial_img_id INTEGER,
	likes INTEGER,
	num_submissions INTEGER,
	UNIQUE (id)
);

CREATE SEQUENCE IF NOT EXISTS img_id START 1 MAXVALUE 2147483647;
CREATE TABLE IF NOT EXISTS img (
	id INTEGER NOT NULL DEFAULT nextval('img_id'),
	uploaded TIMESTAMP,
	caption TEXT,
	UNIQUE (id)
);