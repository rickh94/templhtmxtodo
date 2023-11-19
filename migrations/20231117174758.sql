-- Create "users" table
CREATE TABLE `users` (
  `id` text NULL,
  `created_at` datetime NULL,
  `updated_at` datetime NULL,
  `deleted_at` datetime NULL,
  `name` text NULL,
  `email` text NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "users_email" to table: "users"
CREATE UNIQUE INDEX `users_email` ON `users` (`email`);
-- Create index "idx_users_email" to table: "users"
CREATE INDEX `idx_users_email` ON `users` (`email`);
-- Create "credentials" table
CREATE TABLE `credentials` (
  `id` text NULL,
  `created_at` datetime NULL,
  `updated_at` datetime NULL,
  `deleted_at` datetime NULL,
  `credential_id` blob NOT NULL,
  `public_key` blob NOT NULL,
  `transport` text NOT NULL,
  `attenestation_type` text NULL,
  `flags` text NULL,
  `authenticator` text NULL,
  `user_id` text NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_users_credentials` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "credentials_credential_id" to table: "credentials"
CREATE UNIQUE INDEX `credentials_credential_id` ON `credentials` (`credential_id`);
-- Create index "idx_credentials_credential_id" to table: "credentials"
CREATE INDEX `idx_credentials_credential_id` ON `credentials` (`credential_id`);
-- Create "todos" table
CREATE TABLE `todos` (
  `id` text NULL,
  `created_at` datetime NULL,
  `updated_at` datetime NULL,
  `deleted_at` datetime NULL,
  `text` text NOT NULL,
  `completed` numeric NULL DEFAULT false,
  `user_id` text NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_users_todos` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
);
