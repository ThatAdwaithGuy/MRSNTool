-- migrate:up

ALTER TABLE data 
ALTER COLUMN form_id SET NOT NULL;

ALTER TABLE data 
ALTER COLUMN form_entry_id SET NOT NULL;

-- migrate:down

ALTER TABLE data 
ALTER COLUMN form_id DROP NOT NULL;

ALTER TABLE data 
ALTER COLUMN form_entry_id DROP NOT NULL;

