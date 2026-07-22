-- migrate:up
CREATE TABLE forms (
    id SERIAL PRIMARY KEY,
    form_name VARCHAR NOT NULL 
);

CREATE TABLE dropdown (
    id SERIAL PRIMARY KEY,
    options TEXT[]
);

ALTER TABLE design 
ADD COLUMN dropdown_id INT;

ALTER TABLE design 
ADD CONSTRAINT fk_dropdown
FOREIGN KEY (dropdown_id) REFERENCES dropdown(id) 
ON DELETE CASCADE;

CREATE TABLE data (
    id SERIAL PRIMARY KEY,
    options TEXT[]
);

ALTER TABLE design 
ADD COLUMN form_id INT;

ALTER TABLE design 
ALTER COLUMN form_id SET NOT NULL;

ALTER TABLE design 
ADD CONSTRAINT fk_form 
FOREIGN KEY (form_id) REFERENCES forms(id) 
ON DELETE CASCADE;

ALTER TABLE design 
DROP COLUMN form_name;

-- migrate:down
-- Step 1: Add the old column back
ALTER TABLE design 
ADD COLUMN form_name VARCHAR;

-- Step 2: Restore the string values BEFORE dropping the source table!
UPDATE design 
SET form_name = forms.form_name 
FROM forms 
WHERE design.form_id = forms.id;

ALTER TABLE design 
ALTER COLUMN form_name SET NOT NULL;

-- Step 3: Clean up structural relationships
ALTER TABLE design 
DROP CONSTRAINT fk_form;

ALTER TABLE design 
DROP COLUMN form_id;

ALTER TABLE design 
DROP CONSTRAINT fk_dropdown;

ALTER TABLE design 
DROP COLUMN dropdown_id;

-- Step 4: Drop the tables safely now that dependencies are gone
DROP TABLE IF EXISTS data;
DROP TABLE IF EXISTS dropdown;
DROP TABLE IF EXISTS forms;
