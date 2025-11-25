-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS commands (
    name VARCHAR NOT NULL,
    user_id BIGINT NOT NULL,
    chat_id BIGINT NOT NULL,
    stage INTEGER NOT NULL DEFAULT -1,
    data JSONB DEFAULT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    PRIMARY KEY (user_id, chat_id)
);

COMMENT ON TABLE commands IS 'Current command state for each user in each chat';

COMMENT ON COLUMN commands.name IS 'Command name';
COMMENT ON COLUMN commands.user_id IS 'User ID the command belongs to';
COMMENT ON COLUMN commands.chat_id IS 'Chat ID the command is processed in';
COMMENT ON COLUMN commands.stage IS 'Current step of the command';
COMMENT ON COLUMN commands.data IS 'User data passed between steps';
COMMENT ON COLUMN commands.updated_at IS 'Last update timestamp';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
