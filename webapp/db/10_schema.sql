USE `performania`;

DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id`           BIGINT       NOT NULL AUTO_INCREMENT,
  `login`        VARCHAR(255) NOT NULL,
  `display_name` VARCHAR(255) NOT NULL DEFAULT '',
  `icon`         LONGBLOB     NOT NULL,
  `cover`        LONGBLOB     NOT NULL,
  `created_at`   TIMESTAMP    NOT NULL,
  `updated_at`   TIMESTAMP    NOT NULL,

  PRIMARY KEY (`id`),
  UNIQUE KEY `login` (`login`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

DROP TABLE IF EXISTS `posts`;
CREATE TABLE `posts` (
  `id`         BIGINT    NOT NULL AUTO_INCREMENT,
  `user_id`    BIGINT    NOT NULL,
  `body`       TEXT      NOT NULL,
  `photo`      LONGBLOB,
  `created_at` TIMESTAMP NOT NULL,
  `updated_at` TIMESTAMP NOT NULL,

  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

DROP TABLE IF EXISTS `blocks`;
CREATE TABLE `blocks` (
  `id`         BIGINT    NOT NULL AUTO_INCREMENT,
  `blocker_id` BIGINT    NOT NULL,
  `blocked_id` BIGINT    NOT NULL,
  `created_at` TIMESTAMP NOT NULL,
  `updated_at` TIMESTAMP NOT NULL,

  PRIMARY KEY (`id`),
  UNIQUE KEY `idx` (`blocker_id`, `blocked_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

DROP TABLE IF EXISTS `favorites`;
CREATE TABLE `favorites` (
  `id`         BIGINT    NOT NULL AUTO_INCREMENT,
  `user_id`    BIGINT    NOT NULL,
  `post_id`    BIGINT    NOT NULL,
  `created_at` TIMESTAMP NOT NULL,
  `updated_at` TIMESTAMP NOT NULL,

  PRIMARY KEY (`id`),
  UNIQUE KEY `idx` (`user_id`, `post_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

DROP TABLE IF EXISTS `notifications`;
CREATE TABLE `notifications` (
  `id`         BIGINT    NOT NULL AUTO_INCREMENT,
  `post_id`    BIGINT    NOT NULL,
  `read`       BOOLEAN   NOT NULL DEFAULT FALSE,
  `created_at` TIMESTAMP NOT NULL,
  `updated_at` TIMESTAMP NOT NULL,

  PRIMARY KEY (`id`),
  UNIQUE KEY `idx` (`post_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
