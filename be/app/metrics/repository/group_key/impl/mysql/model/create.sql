CREATE TABLE `dashboard_data_group_tab` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `display_name` varchar(255) NOT NULL,
    `data_group_status` int unsigned NOT NULL,
    `mtime` int unsigned NOT NULL,
    `ctime` int unsigned NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci



CREATE TABLE `dashboard_group_key_tab` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `display_name` varchar(255) NOT NULL,
    `root_id` varchar(255) NOT NULL,
    `parent_id` bigint unsigned NOT NULL,
    `remark` varchar(255) DEFAULT NULL,
    `group_key_status` int unsigned NOT NULL,
    `order_index` smallint unsigned NOT NULL,
    `mtime` int unsigned NOT NULL,
    `ctime` int unsigned NOT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_root_id` (`root_id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci


CREATE TABLE `dashboard_relation_tab` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `group_key_id` bigint unsigned NOT NULL,
    `data_group_id` bigint unsigned NOT NULL,
    `relation_status` int unsigned NOT NULL,
    `mtime` int unsigned NOT NULL,
    `ctime` int unsigned NOT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_group_key_status` (`group_key_id`,`relation_status`),
    KEY `idx_data_group_status` (`data_group_id`,`relation_status`),
    KEY `idx_group_key_data_group_id` (`group_key_id`,`data_group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
