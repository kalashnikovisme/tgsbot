-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE `standupers` (
    `id` INTEGER NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `username` VARCHAR(255) NOT NULL,
    `groupid` BIGINT NOT NULL,
);
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE `standupers`;