-- name: NewDesign :one
INSERT INTO design (
  form_name ,
  label_name,
  data_type,
  is_mandatory ,
  sequence 
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5
)
RETURNING *;

-- name: GetDesignByFormName :many
SELECT *
FROM design
WHERE form_name = $1
ORDER BY sequence;


-- name: GetAllFormNames :many
SELECT DISTINCT form_name
FROM design
ORDER BY form_name;
