CREATE TABLE models (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  volume REAL,
  surface_area REAL,
  updated_at TIMESTAMP DEFAULT NULL,
  created_at TIMESTAMP DEFAULT NULL
);
