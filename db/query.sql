-- name: NewDesign :one
INSERT INTO design (
  form_name ,
  control_level_name ,
  control_type ,
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


