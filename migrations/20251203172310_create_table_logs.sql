-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS logs (
    id SERIAL,
    token UUID NOT NULL,
    data VARCHAR NOT NULL,
    service VARCHAR NOT NULL,
    container_name VARCHAR,
    container_id VARCHAR,
    node VARCHAR,
    node_id VARCHAR,
    hash VARCHAR NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS logs_token_idx ON logs (token);
CREATE INDEX IF NOT EXISTS logs_hash_idx ON logs (hash);
CREATE INDEX IF NOT EXISTS logs_token_hash_createdat_idx ON logs (token, hash, created_at);
-- CREATE INDEX IF NOT EXISTS logs_token_hash_createdat_covering_idx
--    ON logs (token, hash, created_at) INCLUDE (data, label)

COMMENT ON TABLE logs IS 'Logs table';

COMMENT ON COLUMN logs.id IS 'ID';
COMMENT ON COLUMN logs.token IS 'Token of the project the log belongs to';
COMMENT ON COLUMN logs.data IS 'Log data';
COMMENT ON COLUMN logs.service IS 'Label of a service the log belongs to';
COMMENT ON COLUMN logs.container_id IS 'Container ID';
COMMENT ON COLUMN logs.container_name IS 'Container name';
COMMENT ON COLUMN logs.node IS 'Hostname of the node';
COMMENT ON COLUMN logs.node_id IS 'Node ID';
COMMENT ON COLUMN logs.hash IS 'sha-256 hash of the log data without timestamps used for faster comparison';
COMMENT ON COLUMN logs.created_at IS 'Created at';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
