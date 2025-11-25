-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS labels (
    chat_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    labels VARCHAR[] NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    PRIMARY KEY (chat_id, user_id)
);

-- Create Generalized Inverted Index (GIN) to search users by labels faster
CREATE INDEX idx_users_labels ON labels USING GIN (labels);

COMMENT ON TABLE labels IS 'Labels for each user in each chat';

COMMENT ON COLUMN labels.chat_id IS 'Chat ID the label belongs to';
COMMENT ON COLUMN labels.user_id IS 'User ID the label belongs to';
COMMENT ON COLUMN labels.labels IS 'User labels';
COMMENT ON COLUMN labels.updated_at IS 'Last update timestamp';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
