-- Add a foreign key to the models table to link each model to the user that uploaded it
ALTER TABLE models ADD COLUMN IF NOT EXISTS user_id INT NOT NULL;
ALTER TABLE models ADD FOREIGN KEY (user_id) REFERENCES users (id);
