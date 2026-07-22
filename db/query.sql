-- name: NewDesign :one
WITH target_form AS (
  -- 1. Try to insert the form; if it already exists, do nothing but still grab the ID
  INSERT INTO forms (form_name)
  VALUES ($1)
  ON CONFLICT (form_name) DO UPDATE 
    SET form_name = EXCLUDED.form_name -- Trick to force returning the ID even on conflict
  RETURNING id
),
inserted AS (
  -- 2. Insert into the design table using the ID fetched from the first CTE
  INSERT INTO design (
    form_id, 
    label_name,
    data_type,
    is_mandatory,
    sequence 
  ) VALUES (
    (SELECT id FROM target_form),
    $2,
    $3::data_types, -- Cast explicitly if using the custom ENUM type
    $4,
    $5
  )
  RETURNING *
)
-- 3. Return the final selection
SELECT 
  i.id, 
  $1::VARCHAR AS form_name,
  i.label_name, 
  i.data_type, 
  i.is_mandatory, 
  i.sequence,
  i.dropdown_id
FROM inserted i;

-- name: GetDesignByFormName :many
SELECT 
  d.id, 
  f.form_name, 
  d.form_id,
  d.label_name, 
  d.data_type, 
  d.is_mandatory, 
  d.sequence,
  d.dropdown_id
FROM design d
JOIN forms f ON d.form_id = f.id
WHERE f.form_name = $1
ORDER BY d.sequence;

-- name: GetAllFormNames :many
SELECT form_name
FROM forms
ORDER BY form_name;

-- name: ListAllDesigns :many
SELECT 
  d.id, 
  f.form_name, 
  d.form_id,
  d.label_name, 
  d.data_type, 
  d.is_mandatory, 
  d.sequence,
  d.dropdown_id
FROM design d
JOIN forms f ON d.form_id = f.id
ORDER BY f.form_name, d.sequence;

-- name: SetFormEnterable :one 
UPDATE forms
SET enterable = true
WHERE id = $1
RETURNING *;

-- name: GetAllEnterableForms :many
SELECT f.id, f.form_name, f.enterable
FROM forms f
WHERE f.enterable = true;

