-- migrate:up
ALTER TABLE data 
ALTER COLUMN data TYPE VARCHAR;

-- migrate:down
ALTER TABLE data 
ALTER COLUMN data TYPE jsonb;

