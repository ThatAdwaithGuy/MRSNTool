-- migrate:up
CREATE TYPE data_types AS ENUM ('text_input', 'numeric_field', 'checkbox_switch', 'text_area', 'drop_down');

ALTER TABLE data 
DROP COLUMN options;

-- Step 1: Convert design table column to ENUM with an explicit USING clause
ALTER TABLE design 
ALTER COLUMN data_type TYPE data_types USING data_type::data_types;

-- Step 2: Add the jsonb data column with a fallback default object to satisfy NOT NULL
ALTER TABLE data 
ADD COLUMN data jsonb DEFAULT '{}'::jsonb NOT NULL;

-- Step 3: Add data_type straight to the data table as the ENUM type
ALTER TABLE data 
ADD COLUMN data_type data_types DEFAULT 'text_input'::data_types NOT NULL;

-- Step 4: Establish the optional dropdown link
ALTER TABLE data 
ADD COLUMN dropdown_id INT;

ALTER TABLE data 
ADD CONSTRAINT fk_dropdown
FOREIGN KEY (dropdown_id) REFERENCES dropdown(id) 
ON DELETE CASCADE;

-- migrate:down
ALTER TABLE data 
DROP CONSTRAINT fk_dropdown;

ALTER TABLE data 
DROP COLUMN dropdown_id;

ALTER TABLE data 
DROP COLUMN data_type;

-- Fixed typo from 'data' to match your up-migration target
ALTER TABLE data 
DROP COLUMN data;

ALTER TABLE design 
ALTER COLUMN data_type TYPE VARCHAR USING data_type::VARCHAR;

ALTER TABLE data 
ADD COLUMN options TEXT[];

DROP TYPE data_types;
