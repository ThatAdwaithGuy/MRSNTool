-- migrate:up
ALTER TYPE data_types RENAME VALUE 'text_input' TO 'text';
ALTER TYPE data_types RENAME VALUE 'numeric_field' TO 'number';
ALTER TYPE data_types RENAME VALUE 'checkbox_switch' TO 'checkbox';
ALTER TYPE data_types RENAME VALUE 'text_area' TO 'textarea';
ALTER TYPE data_types RENAME VALUE 'drop_down' TO 'select';

-- migrate:down
ALTER TYPE data_types RENAME VALUE   'text'  TO 'text_input';
ALTER TYPE data_types RENAME VALUE   'number'  TO 'numeric_field';
ALTER TYPE data_types RENAME VALUE   'checkbox'  TO 'checkbox_switch';
ALTER TYPE data_types RENAME VALUE   'textarea'  TO 'text_area';
ALTER TYPE data_types RENAME VALUE   'select'  TO 'drop_down';
