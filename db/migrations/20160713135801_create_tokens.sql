
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `tokens` (
	`id` int NOT NULL AUTO_INCREMENT,
	`token` varchar(255) NOT NULL,
	`created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (`id`)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `tokens`;

