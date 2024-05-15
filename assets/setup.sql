CREATE TABLE IF NOT EXISTS activities (
  id INT NOT NULL,
  name VARCHAR NOT NULL,
  distance FLOAT NOT NULL,
  moving_time INT NOT NULL,
  elapsed_time INT NOT NULL,
  total_elevation_gain FLOAT NOT NULL,
  type VARCHAR NOT NULL,
  start_date DATETIME NOT NULL,
  start_latlng VARCHAR NOT NULL,
  end_latlng VARCHAR NOT NULL,
  average_speed FLOAT NOT NULL,
  max_speed FLOAT NOT NULL,
  average_cadence FLOAT NOT NULL,
  kilojoules FLOAT NOT NULL,
  average_heartrate FLOAT NOT NULL,
  max_heartrate FLOAT NOT NULL,
  has_kudoed BOOLEAN NOT NULL DEFAULT false,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS best_efforts (
  id INT NOT NULL,
  activity_id INT NOT NULL,
  name VARCHAR NOT NULL,
  distance FLOAT NOT NULL,
  moving_time INT NOT NULL,
  elapsed_time INT NOT NULL,
  start_date DATETIME NOT NULL,
  pr_rank INT NOT NULL,
  PRIMARY KEY (id),
  FOREIGN KEY (activity_id) REFERENCES activities(id)
);

CREATE TABLE IF NOT EXISTS segments (
  id INT NOT NULL,
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

CREATE TABLE IF NOT EXISTS segment_efforts (
  id INT NOT NULL,
  segment_id INT NOT NULL,
  activity_id INT NOT NULL,
  average_cadence FLOAT NOT NULL,
  average_heartrate FLOAT NOT NULL,
  max_heartrate FLOAT NOT NULL,
  pr_rank INT NOT NULL,
  PRIMARY KEY (id),
  FOREIGN KEY (segment_id) REFERENCES segments(id),
  FOREIGN KEY (activity_id) REFERENCES activities(id)
);
