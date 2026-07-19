-- migrate:up
CREATE TABLE design (
    id SERIAL PRIMARY KEY,
    form_name VARCHAR NOT NULL,
    control_level_name VARCHAR NOT NULL,
    control_type VARCHAR NOT NULL,
    is_mandatory BOOLEAN NOT NULL,
    sequence INT NOT NULL
);

-- migrate:down
DROP TABLE IF EXISTS design;
