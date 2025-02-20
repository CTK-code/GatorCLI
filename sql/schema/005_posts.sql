-- +goose UP
CREATE TABLE posts (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  title TEXT NOT NULL,
  url TEXT NOT NULL,
  description TEXT NOT NULL,
  published_at TIMESTAMP NOT NULL,
  feed_id UUID NOT NULL, 
  CONSTRAINT feed_fk FOREIGN KEY (feed_id)
    REFERENCES feeds(id),
  UNIQUE(url)
);

-- +goose Down
DROP TABLE posts;
