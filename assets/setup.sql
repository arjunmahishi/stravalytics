CREATE DATABASE IF NOT EXISTS strava;

USE strava;

DROP TABLE IF EXISTS activities;

CREATE TABLE IF NOT EXISTS activities (
  id BIGINT NOT NULL,
  name VARCHAR NOT NULL,
  distance FLOAT NOT NULL,
  moving_time INT NOT NULL,
  elapsed_time INT NOT NULL,
  total_elevation_gain FLOAT NOT NULL,
  type VARCHAR NOT NULL,
  start_date DATETIME NOT NULL,
  start_latitude FLOAT NOT NULL,
  start_longitude FLOAT NOT NULL,
  end_latitude FLOAT NOT NULL,
  end_longitude FLOAT NOT NULL,
  average_speed FLOAT NOT NULL,
  max_speed FLOAT NOT NULL,
  average_cadence FLOAT NOT NULL,
  kilojoules FLOAT NOT NULL,
  average_heartrate FLOAT NOT NULL,
  max_heartrate FLOAT NOT NULL,
  has_kudoed BOOLEAN NOT NULL DEFAULT false,
  PRIMARY KEY (id)
);

DROP TABLE IF EXISTS best_efforts;

CREATE TABLE IF NOT EXISTS best_efforts (
  id BIGINT NOT NULL,
  activity_id BIGINT NOT NULL,
  name VARCHAR NOT NULL,
  distance FLOAT NOT NULL,
  moving_time INT NOT NULL,
  elapsed_time INT NOT NULL,
  start_date DATETIME NOT NULL,
  pr_rank INT NOT NULL,
  PRIMARY KEY (id),
  -- FOREIGN KEY (activity_id) REFERENCES activities(id)
);

DROP TABLE IF EXISTS segments;

CREATE TABLE IF NOT EXISTS segments (
  id BIGINT NOT NULL,
  name VARCHAR NOT NULL,
  activity_type VARCHAR NOT NULL,
  distance FLOAT NOT NULL,
  average_grade FLOAT NOT NULL,
  maximum_grade FLOAT NOT NULL,
  elevation_high FLOAT NOT NULL,
  elevation_low FLOAT NOT NULL,
  start_latlng VARCHAR NOT NULL,
  end_latlng VARCHAR NOT NULL,
  city VARCHAR NOT NULL,
  state VARCHAR NOT NULL,
  country VARCHAR NOT NULL,
  PRIMARY KEY (id)
);

DROP TABLE IF EXISTS segment_efforts;

CREATE TABLE IF NOT EXISTS segment_efforts (
  id BIGINT NOT NULL,
  activity_id BIGINT NOT NULL,
  segment_id BIGINT NOT NULL,
  average_cadence FLOAT NOT NULL,
  average_heartrate FLOAT NOT NULL,
  max_heartrate FLOAT NOT NULL,
  pr_rank INT NOT NULL,
  PRIMARY KEY (id),
  -- FOREIGN KEY (segment_id) REFERENCES segments(id),
  -- FOREIGN KEY (activity_id) REFERENCES activities(id)
);

