-- +migrate Up
CREATE TABLE content_tag (
    content_id TEXT NOT NULL,
    tag_id TEXT NOT NULL,
    PRIMARY KEY (content_id, tag_id),
    FOREIGN KEY (content_id) REFERENCES content(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tag(id) ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE content_tag;
