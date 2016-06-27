
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE episodes (
	id int NOT NULL,
	title varchar(128) NOT NULL,
	status tinyint(4) NOT NULL DEFAULT 0,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE episodes;
