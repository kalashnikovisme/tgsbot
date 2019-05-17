-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE `standups` (
    `id` INTEGER NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `created` DATETIME NOT NULL,
    `modified` DATETIME NOT NULL,
    `username` VARCHAR(255) NOT NULL,
    `comment` VARCHAR(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `groupid` BIGINT NOT NULL,
    KEY (`created`, `username`)
);
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE `standups`;