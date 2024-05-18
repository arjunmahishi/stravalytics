INSERT into activities (
  id, name, distance, moving_time, elapsed_time, total_elevation_gain, type, start_date, start_latlng, end_latlng, average_speed, max_speed, average_cadence, kilojoules, average_heartrate, max_heartrate, has_kudoed)
VALUES (%s);
