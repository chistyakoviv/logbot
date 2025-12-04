-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS logs (
    id SERIAL,
    token UUID NOT NULL,
    data VARCHAR NOT NULL,
    label VARCHAR NOT NULL,
    hash VARCHAR NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS logs_token_idx ON logs (token);
CREATE INDEX IF NOT EXISTS logs_label_idx ON logs (label);
CREATE INDEX IF NOT EXISTS logs_hash_idx ON logs (hash);

COMMENT ON TABLE logs IS 'Logs table';

COMMENT ON COLUMN logs.id IS 'ID';
COMMENT ON COLUMN logs.token IS 'Token of the project the log belongs to';
COMMENT ON COLUMN logs.data IS 'Log data';
COMMENT ON COLUMN logs.label IS 'Label of a service the log belongs to';
COMMENT ON COLUMN logs.hash IS 'sha-256 hash of the log data without timestamps used for faster comparison';
COMMENT ON COLUMN logs.created_at IS 'Created at';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
