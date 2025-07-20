CREATE TABLE IF NOT EXISTS worklogs (
    jira_key varchar(256) UNIQUE PRIMARY KEY NOT NULL,
    jira_data text NOT NULL,
    duration integer NOT NULL,
    running boolean NOT NULL,
    log_text text NOT NULL,
    updated_at datetime DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Create trigger to automatically update the timestamp on UPDATE
CREATE TRIGGER IF NOT EXISTS update_worklogs_timestamp
  AFTER UPDATE ON worklogs
BEGIN
    UPDATE worklogs
    SET updated_at = CURRENT_TIMESTAMP
    WHERE jira_key = NEW.jira_key;
END;
