-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS labels (
    chat_id BIGINT NOT NULL,
    username VARCHAR NOT NULL,
    labels VARCHAR[] NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    PRIMARY KEY (chat_id, username)
);

-- Create Generalized Inverted Index (GIN) to search users by labels faster
CREATE INDEX idx_users_labels ON labels USING GIN (labels);

-- Index users separately from composite key index so that it would be possible
-- to search by user and label combining btree and gin indexes using bitmap index
CREATE INDEX idx_labels_user ON labels (username);

COMMENT ON TABLE labels IS 'Labels for each user in each chat';

COMMENT ON COLUMN labels.chat_id IS 'Chat ID the label belongs to';
COMMENT ON COLUMN labels.username IS 'User username the label belongs to';
COMMENT ON COLUMN labels.labels IS 'User labels';
COMMENT ON COLUMN labels.updated_at IS 'Last update timestamp';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
