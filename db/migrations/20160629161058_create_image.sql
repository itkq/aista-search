
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `images` (
	`id` int NOT NULL AUTO_INCREMENT,
	`episode_id` int NOT NULL,
	`path` varchar(128) NOT NULL,
	`url` varchar(255),
	`sentence` varchar(255),
	`created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (`id`)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `images`;
