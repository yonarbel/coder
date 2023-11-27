-- name: GetReplicasUpdatedAfter :many
SELECT * FROM replicas WHERE updated_at > $1 AND stopped_at IS NULL;

-- name: GetReplicaByID :one
SELECT * FROM replicas WHERE id = $1;

-- name: InsertReplica :one
INSERT INTO replicas (
    id,
    created_at,
    started_at,
    updated_at,
    hostname,
    region_id,
    relay_address,
    version,
    database_latency,
	"primary"
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING *;

-- name: UpsertReplica :one
INSERT INTO replicas (
    id,
    updated_at,
    created_at,
    started_at,
    hostname,
    region_id,
    relay_address,
    version,
    database_latency,
	"primary"
) VALUES ($1, $2, $2, $3, $5, $6, $7, $8, $10, $11) ON CONFLICT (id) DO UPDATE SET
    updated_at = $2,
    started_at = $3,
    stopped_at = $4,
    hostname = $5,
    region_id = $6,
    relay_address = $7,
    version = $8,
    error = $9,
    database_latency = $10,
	"primary" = $11
RETURNING *;

-- name: DeleteReplicasUpdatedBefore :exec
DELETE FROM replicas WHERE updated_at < $1;
