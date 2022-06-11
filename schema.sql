BEGIN;

CREATE TABLE schema_version (id integer NOT NULL);
INSERT INTO schema_version VALUES (1);

-- TODO: Check if this makes sense
CREATE TABLE users (
	id serial PRIMARY KEY,
	email VARCHAR(256) NOT NULL UNIQUE,
	password VARCHAR NOT NULL,
	bio VARCHAR,
	url VARCHAR,
	suspension_notice TEXT,
	created_at timestamp WITH TIME ZONE NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
	updated_at timestamp WITH TIME ZONE NOT NULL DEFAULT (now() AT TIME ZONE 'utc')
);

CREATE TYPE visibility AS ENUM (
	'UNLISTED',
	'PRIVATE',
	'PUBLIC'
);

CREATE TABLE projects (
	id serial PRIMARY KEY,
	name VARCHAR(256) NOT NULL,
	user_id integer NOT NULL UNIQUE REFERENCES "users"(id),
	visibility visibility NOT NULL,
	created_at timestamp WITH TIME ZONE NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
	updated_at timestamp WITH TIME ZONE NOT NULL DEFAULT (now() AT TIME ZONE 'utc')
);

CREATE TABLE photos (
	id serial PRIMARY KEY,
	file_name varchar NOT NULL,
	file_size integer NOT NULL,
	content_type varchar NOT NULL,
	-- TODO: Keep this only if momento will hold multiple user accounts
	user_id integer NOT NULL UNIQUE REFERENCES "users"(id),
	project_id integer NOT NULL UNIQUE REFERENCES "projects"(id),
	created_at timestamp WITH TIME ZONE NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
	updated_at timestamp WITH TIME ZONE NOT NULL DEFAULT (now() AT TIME ZONE 'utc')
);

COMMIT;
