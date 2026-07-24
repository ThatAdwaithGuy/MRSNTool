-- migrate:up
ALTER TABLE data 
ADD COLUMN form_id INT,
ADD CONSTRAINT fk_form_id
  FOREIGN KEY (form_id)
  REFERENCES forms(id)
  ON DELETE CASCADE;

ALTER TABLE data 
ADD COLUMN form_entry_id INT;


-- migrate:down
ALTER TABLE data 
DROP CONSTRAINT IF EXISTS fk_form_id,
DROP COLUMN IF EXISTS form_id;

ALTER TABLE data 
DROP COLUMN IF EXISTS form_entry_id;
