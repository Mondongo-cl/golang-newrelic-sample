
GRANT ALL ON *.* to 'root'@'%' IDENTIFIED BY 'password!';
FLUSH PRIVILEGES;
DROP DATABASE IF EXISTS `customer`;
CREATE DATABASE `customer` 
USE customer;
DROP TABLE IF EXISTS `customers`;
CREATE TABLE `customers` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `first_name` varchar(64) NOT NULL,
  `last_name` varchar(64) DEFAULT NULL,
  `email` varchar(256) NOT NULL,
  `dob` datetime(3) NOT NULL,
  `country_code` varchar(4) NOT NULL,
  `mobile_number` varchar(16) NOT NULL,
  `profile_pic_res_id` varchar(32) DEFAULT NULL,
  `bl_id` varchar(32) DEFAULT NULL,
  `bl_qr` varchar(32) DEFAULT NULL,
  `radixx_id` smallint(6) DEFAULT NULL,
  `iomob_id` smallint(6) DEFAULT NULL,
  `is_verified` tinyint(1) DEFAULT NULL,
  `marketting_com` tinyint(1) DEFAULT NULL,
  `password` varchar(128) DEFAULT NULL,
  `minor_consent` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `email` (`email`),
  KEY `idx_customers_deleted_at` (`deleted_at`),
  KEY `idx_customers_email` (`email`),
  KEY `idx_customers_country_code` (`country_code`),
  KEY `idx_customers_mobile_number` (`mobile_number`)
) ENGINE=InnoDB AUTO_INCREMENT=140663 DEFAULT CHARSET=latin1;
