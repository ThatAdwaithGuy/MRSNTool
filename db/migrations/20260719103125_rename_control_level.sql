-- migrate:up
ALTER TABLE design
    RENAME COLUMN control_level_name TO label_name;

ALTER TABLE design
    RENAME COLUMN control_type TO data_type;

-- migrate:down
ALTER TABLE design
    RENAME COLUMN label_name TO control_level_name;

ALTER TABLE design
    RENAME COLUMN data_type TO control_type;
