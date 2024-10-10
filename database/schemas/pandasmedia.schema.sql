SET NAMES utf8mb4 ;
SET @MYSQLDUMP_TEMP_LOG_BIN = @@SESSION.SQL_LOG_BIN;
SET @@SESSION.SQL_LOG_BIN= 0;

DROP DATABASE IF EXISTS `pandasmedia`;
CREATE DATABASE `pandasmedia`;
USE `pandasmedia`;

DROP TABLE IF EXISTS `categories_clusters`;
CREATE TABLE `categories_clusters` (
    `uuid` VARCHAR(36) NOT NULL,
    `name` VARCHAR(100) NOT NULL,
    `fs_path` VARCHAR(300) NOT NULL,
    `filter_category` VARCHAR(40) NOT NULL,
    `root_category` VARCHAR(40) NOT NULL,
    PRIMARY KEY (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `categorys`;
SET character_set_client = utf8mb4 ;
CREATE TABLE `categorys` (
    `uuid` varchar(40) NOT NULL,
    `name` varchar(100) NOT NULL,
    `fullpath` varchar(300) NOT NULL,
    `parent` varchar(40) DEFAULT NULL,
    `cluster` varchar(36) DEFAULT NULL,
    PRIMARY KEY (`uuid`),
    KEY `parent_fk` (`parent`),
    KEY `cluster_fk` (`cluster`),
    CONSTRAINT `parent_fk` FOREIGN KEY (`parent`) REFERENCES `categorys` (`uuid`) ON DELETE CASCADE,
    CONSTRAINT `cluster_fk` FOREIGN KEY (`cluster`) REFERENCES `categories_clusters` (`uuid`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Adapting an old db schema the new one
-- See: https://github.com/Gerardo115pp/PandasWorld/issues/14
-- END of adapting an old db schema the new one

DROP TABLE IF EXISTS `medias`;
SET character_set_client = utf8mb4 ;
CREATE TABLE `medias` (`uuid` varchar(40) NOT NULL,
    `name` varchar(200) NOT NULL,
    `last_seen` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `main_category` varchar(40) NOT NULL,
    `type` enum('IMAGE','VIDEO') DEFAULT NULL,
    `downloaded_from` INT,
    PRIMARY KEY (`uuid`),
    KEY `main_category_fk` (`main_category`),
    CONSTRAINT `main_category_fk` FOREIGN KEY (`main_category`) REFERENCES `categorys` (`uuid`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `downloads`;
CREATE TABLE `downloads`(
    `uuid` INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `url` VARCHAR(120) NOT NULL,
    `created_on` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `files_downloaded` VARCHAR(200) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB3;

SET @@SESSION.SQL_LOG_BIN = @MYSQLDUMP_TEMP_LOG_BIN;

ALTER TABLE `medias` ADD CONSTRAINT `fk_download` FOREIGN KEY (`downloaded_from`) REFERENCES `downloads`(`uuid`);

CREATE VIEW media_paths AS SELECT medias.uuid, CONCAT(categorys.fullpath, '\\', medias.name) AS path FROM categorys, medias WHERE categorys.uuid=medias.main_category AND categorys.fullpath NOT LIKE '%\\';

DROP PROCEDURE IF EXISTS `delete_tags_by_name`;
DELIMITER //
CREATE PROCEDURE `delete_tags_by_name` ( IN `tag_name` VARCHAR(120)) 
BEGIN
  DELETE FROM mediastags WHERE tag=(SELECT id FROM tags WHERE name=`tag_name`);
  DELETE FROM tags WHERE name=`tag_name`;
END//
DELIMITER ;


DROP PROCEDURE IF EXISTS `delete_media_by_id`;
DELIMITER //
CREATE PROCEDURE `delete_media_by_id` ( IN `media_uuid` VARCHAR(40))
BEGIN
  DELETE FROM `mediastags` WHERE media=`media_uuid`;
  DELETE FROM `medias` WHERE uuid=`media_uuid` LIMIT 1;
END//
DELIMITER ;

DROP PROCEDURE IF EXISTS `delete_category_f`;
DELIMITER //
CREATE PROCEDURE `delete_category_f` (IN `category_uuid` VARCHAR(40))
BEGIN
    DELETE FROM `mediastags` WHERE `media` IN (SELECT uuid FROM `medias` WHERE main_category=`category_uuid`);
    DELETE FROM `medias` WHERE `main_category`=`category_uuid`;
    DELETE FROM `categorys` WHERE `uuid`=`category_uuid` LIMIT 1;
END//
DELIMITER ;
