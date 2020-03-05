CREATE TABLE `users` (
    `user_id` STRING(36) NOT NULL,
	`name` STRING(MAX) NOT NULL,
	`uid` STRING(MAX) NOT NULL,
	`created_at` TIMESTAMP NOT NULL,
	`updated_at` TIMESTAMP NOT NULL
) PRIMARY KEY  (`user_id`);
CREATE UNIQUE INDEX idx_users_uid ON users (uid);
CREATE TABLE `friends` (
	`friend_id` STRING(36) NOT NULL,
	`user_id` STRING(36) NOT NULL,
	`to_id` STRING(MAX) NOT NULL,
	`created_at` TIMESTAMP NOT NULL,
	`updated_at` TIMESTAMP NOT NULL
) PRIMARY KEY  (`user_id`, `friend_id`),
INTERLEAVE IN PARENT `users` ;
CREATE UNIQUE INDEX idx_friends_user_id_to_id ON friends (user_id, to_id);
CREATE  INDEX idx_friends_to_id ON friends (to_id);
CREATE TABLE user_items (
  user_item_id STRING(36) NOT NULL,
  user_id      STRING(36) NOT NULL,
  item_id      STRING(36) NOT NULL,
  created_at   TIMESTAMP NOT NULL,
  updated_at   TIMESTAMP NOT NULL,
) PRIMARY KEY(user_item_id, user_id),
  INTERLEAVE IN PARENT users ON DELETE CASCADE;
CREATE TABLE items (
  item_id      STRING(36) NOT NULL,
  name STRING(MAX) NOT NULL,
  created_at   TIMESTAMP NOT NULL,
  updated_at   TIMESTAMP NOT NULL,
) PRIMARY KEY(item_id);
CREATE TABLE item_skus (
  item_sku_id      STRING(36) NOT NULL,
  item_id      STRING(36) NOT NULL,
  name STRING(MAX) NOT NULL,
  created_at   TIMESTAMP NOT NULL,
  updated_at   TIMESTAMP NOT NULL,
) PRIMARY KEY(item_sku_id, item_id),
  INTERLEAVE IN PARENT items ON DELETE CASCADE;
