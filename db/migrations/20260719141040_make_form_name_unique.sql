-- migrate:up
ALTER TABLE forms 
ADD CONSTRAINT uni_form_name UNIQUE (form_name);
-- migrate:down
ALTER TABLE forms 
DROP CONSTRAINT uni_form_name;
