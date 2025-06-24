CREATE TABLE `dashboard_data_source_tab` (
         `id` bigint unsigned NOT NULL AUTO_INCREMENT,
         `display_name` varchar(255) NOT NULL,
         `tips` varchar(255) DEFAULT NULL,
         `remark` varchar(255) DEFAULT NULL,
         `data_name` varchar(255) DEFAULT NULL,
         `data_source_status` int unsigned NOT NULL,
         `group_key` varchar(255) NOT NULL,
         `enable_group_key_name` tinyint(1) NOT NULL,
         `data_type` varchar(45) NOT NULL,
         `sort_fields` varchar(255) NOT NULL,
         `mtime` int unsigned NOT NULL,
         `ctime` int unsigned NOT NULL,
         PRIMARY KEY (`id`),
         KEY `idx_group_key` (`group_key`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


CREATE TABLE `dashboard_data_source_meta_tab` (
      `id` int unsigned NOT NULL AUTO_INCREMENT,
      `data_source_id` int unsigned NOT NULL,
      `field_name` varchar(255) DEFAULT NULL,
      `display_name` varchar(255) NOT NULL,
      `data_type` varchar(45) NOT NULL,
      `output_data_type` varchar(45) NOT NULL,
      `field_tips` varchar(512) DEFAULT NULL,
      `action` text NOT NULL,
      `enum` text NOT NULL,
      `formatter` varchar(1024) NOT NULL,
      `template` varchar(1024) NOT NULL,
      `mtime` int unsigned NOT NULL,
      `ctime` int unsigned NOT NULL,
      `nested_path` varchar(255) DEFAULT NULL,
      `data_source_meta_status` int unsigned NOT NULL,
      PRIMARY KEY (`id`),
      KEY `idx_data_source_id_status` (`data_source_id`,`data_source_meta_status`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
