-- name: GetWorklog :one
SELECT * FROM worklogs
WHERE jira_key = ? LIMIT 1;

-- name: ListWorklog :many
SELECT * FROM worklogs
ORDER BY jira_key;

-- name: CreateWorklog :exec
INSERT INTO worklogs (jira_key, jira_data, duration, running, log_text)
VALUES (?, ?, ?, ?, ?);

-- name: UpdateWorklog :exec
UPDATE worklogs
set jira_data = ?,
    duration = ?,
    running = ?,
    log_text = ?
WHERE jira_key = ?;

-- name: DeleteWorklog :exec
DELETE FROM worklogs
WHERE jira_key = ?;
