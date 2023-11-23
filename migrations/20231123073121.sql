-- Disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- Create "new_credentials" table
CREATE TABLE `new_credentials` (
  `id` text NULL,
  `created_at` datetime NULL,
  `updated_at` datetime NULL,
  `deleted_at` datetime NULL,
  `credential_id` blob NOT NULL,
  `public_key` blob NOT NULL,
  `transport` text NOT NULL,
  `attestation_type` text NULL,
  `flags` text NULL,
  `authenticator` text NULL,
  `user_id` text NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_users_credentials` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Copy rows from old table "credentials" to new temporary table "new_credentials"
INSERT INTO `new_credentials` (`id`, `created_at`, `updated_at`, `deleted_at`, `credential_id`, `public_key`, `transport`, `flags`, `authenticator`, `user_id`) SELECT `id`, `created_at`, `updated_at`, `deleted_at`, `credential_id`, `public_key`, `transport`, `flags`, `authenticator`, `user_id` FROM `credentials`;
-- Drop "credentials" table after copying rows
DROP TABLE `credentials`;
-- Rename temporary table "new_credentials" to "credentials"
ALTER TABLE `new_credentials` RENAME TO `credentials`;
-- Create index "credentials_credential_id" to table: "credentials"
CREATE UNIQUE INDEX `credentials_credential_id` ON `credentials` (`credential_id`);
-- Create index "idx_credentials_credential_id" to table: "credentials"
CREATE INDEX `idx_credentials_credential_id` ON `credentials` (`credential_id`);
-- Enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;
