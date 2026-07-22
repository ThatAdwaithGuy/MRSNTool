-- name: NewDesign :one
WITH target_form AS (
  INSERT INTO forms (form_name)
  VALUES ($1)
  ON CONFLICT (form_name) DO UPDATE 
    SET form_name = EXCLUDED.form_name 
  RETURNING id
),
inserted AS (
  INSERT INTO design (
    form_id, 
    label_name,
    data_type,
    is_mandatory,
    sequence, 
    dropdown_id
  ) VALUES (
    (SELECT id FROM target_form),
    $2,
    $3::data_types, 
    $4,
    $5,
    $6
  )
  RETURNING *
)
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

-- name: GetDesignByFormID :many
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
WHERE f.id = $1
ORDER BY d.sequence;

-- name: GetAllFormNames :many
SELECT f.id, f.form_name, f.enterable
FROM forms f
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

-- name: SetFormEnterable :exec
UPDATE forms 
SET enterable = true
WHERE form_name = $1;

-- name: GetAllEnterableForms :many
SELECT f.id, f.form_name, f.enterable
FROM forms f
WHERE f.enterable = true;

-- name: NewDropDown :one
INSERT INTO dropdown (
  options
) VALUES (
  $1
) 
RETURNING *;
