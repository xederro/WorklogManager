package db

import _ "embed"

//go:embed "worklogs_schema.sql"
var WorklogsSchema string
