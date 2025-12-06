-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS last_sent (
    chat_id BIGINT NOT NULL,
    token UUID NOT NULL,
    hash VARCHAR NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    PRIMARY KEY (chat_id, token, hash)
);

COMMENT ON TABLE last_sent IS 'Last sent message for each chat';

COMMENT ON COLUMN last_sent.chat_id IS 'Chat ID';
COMMENT ON COLUMN last_sent.token IS 'Token of the project the log belongs to';
COMMENT ON COLUMN last_sent.hash IS 'sha-256 hash of the log data without timestamps used for faster comparison';
COMMENT ON COLUMN last_sent.updated_at IS 'Last update timestamp';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
