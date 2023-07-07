
-- name: CreateUser :one
INSERT INTO public.users (username, full_name, hash_password, email) VALUES($1, $2, $3, $4) RETURNING *;

-- name: UpdateUserPassword :one
UPDATE public.users SET  hash_password=$1, changed_password_at=now() WHERE username=$2 RETURNING *;

-- name: UpdateUserEmail :one
UPDATE public.users SET  email=$1 WHERE username=$2 RETURNING *;

-- name: UpdateUserFullName :one
UPDATE public.users SET  full_name=$1 WHERE username=$2 RETURNING *;


-- name: GetUser :one
SELECT * FROM public.users WHERE username=$1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM public.users ORDER BY username LIMIT $1 OFFSET $2;

-- name: DeleteUser :exec
DELETE FROM public.users WHERE username=$1;