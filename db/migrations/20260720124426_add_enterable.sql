-- migrate:up
ALTER TABLE forms
ADD enterable bool;

-- migrate:down
ALTER TABLE forms 
DROP COLUMN enterable;

