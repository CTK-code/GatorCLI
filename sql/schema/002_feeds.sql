-- +goose Up
CREATE TABLE feeds (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  name VARCHAR(50) NOT NULL,
  url VARCHAR(100) NOT NULL,
  user_id UUID NOT NULL,
  CONSTRAINT user_fk FOREIGN KEY (user_id)
    REFERENCES users(id) ON DELETE CASCADE,
  UNIQUE(url)
);

-- +goose Down
DROP TABLE feeds;
