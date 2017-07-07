-- +migrate Up
-- START TRANSACTION;

CREATE TABLE `_users` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(32) NOT NULL,
    `usertype` int(11) unsigned NOT NULL DEFAULT 0,
    `email` varchar(36) NOT NULL,
    `ctime` TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    `utime` TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    `dtime` TIMESTAMP(6) NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    UNIQUE KEY (`name`, `dtime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE VIEW `users` AS SELECT * FROM _users WHERE _users.dtime = 0;

CREATE TABLE `_teams` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(32) NOT NULL,
    `displayname` varchar(32) NOT NULL,
    `descp` varchar(128) NOT NULL,
    `ctime` TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    `utime` TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    `dtime` TIMESTAMP(6) NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    UNIQUE KEY (`name`, `dtime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE VIEW `teams` AS SELECT * FROM _teams WHERE _teams.dtime = 0;

CREATE TABLE `_spaces` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(32) NOT NULL,
    `ctime` TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    `utime` TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    `dtime` TIMESTAMP(6) NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    UNIQUE KEY (`name`, `dtime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE VIEW `spaces` AS SELECT * FROM _spaces WHERE _spaces.dtime = 0;

CREATE TABLE `_apps` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(32) NOT NULL,
    `descp` varchar(128) NOT NULL,
    `space_id` int(11) unsigned NOT NULL,
    `ctime` TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    `utime` TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    `dtime` TIMESTAMP(6) NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    UNIQUE KEY (`name`, `space_id`, `dtime`),
    KEY `space_id_idx` (`space_id`),
    CONSTRAINT `apps_space_id_cst` FOREIGN KEY (`space_id`) REFERENCES `_spaces` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE VIEW `apps` AS SELECT * FROM _apps WHERE _apps.dtime = 0;

CREATE TABLE `_services` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(32) NOT NULL,
    `app_id` int(11) unsigned NOT NULL,
    `ctime` TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    `utime` TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    `dtime` TIMESTAMP(6) NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    UNIQUE KEY (`name`, `app_id`, `dtime`),
    KEY `app_id_idx` (`app_id`),
    CONSTRAINT `services_app_id_cst` FOREIGN KEY (`app_id`) REFERENCES `_apps` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE VIEW `services` AS SELECT * FROM _services WHERE _services.dtime = 0;

CREATE TABLE `users_teams` (
    `user_id` int(11) unsigned NOT NULL,
    `team_id` int(11) unsigned NOT NULL,
    `membertype` int(11) unsigned NOT NULL DEFAULT 1,
    PRIMARY KEY (`user_id`, `team_id`),
    KEY `user_id_idx` (`user_id`),
    KEY `team_id_idx` (`team_id`),
    CONSTRAINT `users_teams_user_id_cst` FOREIGN KEY (`user_id`) REFERENCES `_users` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT `users_teams_team_id_cst` FOREIGN KEY (`team_id`) REFERENCES `_teams` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `users_spaces` (
    `user_id` int(11) unsigned NOT NULL,
    `space_id` int(11) unsigned NOT NULL,
    PRIMARY KEY (`user_id`, `space_id`),
    CONSTRAINT `users_spaces_user_id_cst` FOREIGN KEY (`user_id`) REFERENCES `_users` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT `users_spaces_space_id_cst` FOREIGN KEY (`space_id`) REFERENCES `_spaces` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `teams_spaces` (
    `team_id` int(11) unsigned NOT NULL,
    `space_id` int(11) unsigned NOT NULL,
    PRIMARY KEY (`team_id`, `space_id`),
    CONSTRAINT `teams_spaces_team_id_cst` FOREIGN KEY (`team_id`) REFERENCES `_teams` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT `teams_spaces_space_id_cst` FOREIGN KEY (`space_id`) REFERENCES `_spaces` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- COMMIT;
-- +migrate Down   
DROP TABLE teams_spaces;
DROP TABLE users_spaces;
DROP TABLE users_teams; 
DROP VIEW services;
DROP TABLE _services; 
DROP VIEW apps;
DROP TABLE _apps;
DROP VIEW spaces;
DROP TABLE _spaces;
DROP VIEW teams;
DROP TABLE _teams;
DROP VIEW users;
DROP TABLE _users;



 















