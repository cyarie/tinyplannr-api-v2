/* A magical schema used for generating the tables and such for the TinyPlannr DB*/

DROP SCHEMA tinyplannr_api CASCADE;
CREATE SCHEMA tinyplannr_api;

DROP TABLE IF EXISTS tinyplannr_api.user;
CREATE TABLE tinyplannr_api.user (
  user_id SERIAL PRIMARY KEY,
  email VARCHAR(255) UNIQUE,
  first_name VARCHAR(255),
  last_name VARCHAR(255),
  zip_code INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  create_dt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  update_dt TIMESTAMP
);
CREATE INDEX api_user_idx ON tinyplannr_api.user (user_id, email);

DROP TABLE IF EXISTS tinyplannr_api.event;
CREATE TABLE tinyplannr_api.event (
  event_id SERIAL PRIMARY KEY,
  user_id INTEGER REFERENCES tinyplannr_api.user (user_id),
  title TEXT,
  description TEXT,
  location TEXT,
  all_day BOOLEAN,
  start_dt TIMESTAMP,
  end_dt TIMESTAMP,
  create_dt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  update_dt TIMESTAMP
);
CREATE INDEX api_event_idx ON tinyplannr_api.event (user_id);

DROP SCHEMA tinyplannr_auth CASCADE;
CREATE SCHEMA tinyplannr_auth;

DROP TABLE IF EXISTS tinyplannr_auth.user;
CREATE TABLE tinyplannr_auth.user (
  auth_id SERIAL PRIMARY KEY,
  user_id INTEGER REFERENCES tinyplannr_api.user (user_id) UNIQUE,
  email VARCHAR(255) REFERENCES tinyplannr_api.user (email) UNIQUE,
  hash_pw TEXT,
  create_dt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  update_dt TIMESTAMP,
  last_login_dt TIMESTAMP
);
CREATE INDEX auth_user_idx ON tinyplannr_auth.user (user_id, email);

DROP TABLE IF EXISTS tinyplannr_auth.session;
CREATE TABLE tinyplannr_auth.session (
  session_key varchar(255) NOT NULL PRIMARY KEY,
  user_id INTEGER REFERENCES tinyplannr_auth.user (user_id),
  email VARCHAR(255) REFERENCES tinyplannr_auth.user (email),
  is_valid BOOL DEFAULT TRUE,
  create_dt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  update_dt TIMESTAMP,
  expire_dt TIMESTAMP
);
CREATE INDEX auth_session_idx ON tinyplannr_auth.session (session_key);