-- MySQL dump 10.13  Distrib 8.0.42, for Linux (x86_64)
--
-- Host: localhost    Database: liansangyu
-- ------------------------------------------------------
-- Server version	8.0.42-0ubuntu0.24.04.1

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `elders`
--

DROP TABLE IF EXISTS `elders`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `elders` (
  `openid` varchar(255) NOT NULL,
  `created_at` datetime(3) NOT NULL COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL COMMENT '更新时间',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT '删除时间',
  `disease` longtext,
  `longitude` double NOT NULL COMMENT '经度',
  `latitude` double NOT NULL COMMENT '纬度',
  PRIMARY KEY (`openid`),
  KEY `idx_elders_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_elders_user` FOREIGN KEY (`openid`) REFERENCES `users` (`openid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `elders`
--

LOCK TABLES `elders` WRITE;
/*!40000 ALTER TABLE `elders` DISABLE KEYS */;
INSERT INTO `elders` VALUES ('elder','2025-06-30 19:05:22.534','2025-06-30 19:05:22.534',NULL,'asda',123,456);
/*!40000 ALTER TABLE `elders` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `monitors`
--

DROP TABLE IF EXISTS `monitors`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `monitors` (
  `openid` varchar(255) NOT NULL,
  `created_at` datetime(3) NOT NULL COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL COMMENT '更新时间',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT '删除时间',
  `elder_openid` varchar(255) NOT NULL,
  `passed` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`openid`),
  KEY `idx_monitors_deleted_at` (`deleted_at`),
  KEY `fk_monitors_elder` (`elder_openid`),
  CONSTRAINT `fk_monitors_elder` FOREIGN KEY (`elder_openid`) REFERENCES `elders` (`openid`),
  CONSTRAINT `fk_monitors_user` FOREIGN KEY (`openid`) REFERENCES `users` (`openid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `monitors`
--

LOCK TABLES `monitors` WRITE;
/*!40000 ALTER TABLE `monitors` DISABLE KEYS */;
/*!40000 ALTER TABLE `monitors` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `notifications`
--

DROP TABLE IF EXISTS `notifications`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `notifications` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `created_at` datetime(3) NOT NULL COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL COMMENT '更新时间',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT '删除时间',
  `user_openid` varchar(255) NOT NULL,
  `title` longtext,
  `content` longtext,
  PRIMARY KEY (`id`),
  KEY `idx_notifications_deleted_at` (`deleted_at`),
  KEY `fk_notifications_user` (`user_openid`),
  CONSTRAINT `fk_notifications_user` FOREIGN KEY (`user_openid`) REFERENCES `users` (`openid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `notifications`
--

LOCK TABLES `notifications` WRITE;
/*!40000 ALTER TABLE `notifications` DISABLE KEYS */;
/*!40000 ALTER TABLE `notifications` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `organization_admins`
--

DROP TABLE IF EXISTS `organization_admins`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `organization_admins` (
  `organization_openid` varchar(255) NOT NULL,
  `admin_openid` varchar(255) NOT NULL,
  `passed` tinyint(1) NOT NULL DEFAULT '0',
  KEY `fk_organization_admins_organization` (`organization_openid`),
  KEY `fk_organization_admins_user` (`admin_openid`),
  CONSTRAINT `fk_organization_admins_organization` FOREIGN KEY (`organization_openid`) REFERENCES `organizations` (`openid`),
  CONSTRAINT `fk_organization_admins_user` FOREIGN KEY (`admin_openid`) REFERENCES `users` (`openid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `organization_admins`
--

LOCK TABLES `organization_admins` WRITE;
/*!40000 ALTER TABLE `organization_admins` DISABLE KEYS */;
/*!40000 ALTER TABLE `organization_admins` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `organization_elders`
--

DROP TABLE IF EXISTS `organization_elders`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `organization_elders` (
  `organization_openid` varchar(255) NOT NULL,
  `elder_openid` varchar(255) NOT NULL,
  `passed` tinyint(1) NOT NULL DEFAULT '0',
  KEY `fk_organization_elders_organization` (`organization_openid`),
  KEY `fk_organization_elders_elder` (`elder_openid`),
  CONSTRAINT `fk_organization_elders_elder` FOREIGN KEY (`elder_openid`) REFERENCES `elders` (`openid`),
  CONSTRAINT `fk_organization_elders_organization` FOREIGN KEY (`organization_openid`) REFERENCES `organizations` (`openid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `organization_elders`
--

LOCK TABLES `organization_elders` WRITE;
/*!40000 ALTER TABLE `organization_elders` DISABLE KEYS */;
/*!40000 ALTER TABLE `organization_elders` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `organization_volunteers`
--

DROP TABLE IF EXISTS `organization_volunteers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `organization_volunteers` (
  `organization_openid` varchar(255) NOT NULL,
  `volunteer_openid` varchar(255) NOT NULL,
  `passed` tinyint(1) NOT NULL DEFAULT '0',
  KEY `fk_organization_volunteers_organization` (`organization_openid`),
  KEY `fk_organization_volunteers_volunteer` (`volunteer_openid`),
  CONSTRAINT `fk_organization_volunteers_organization` FOREIGN KEY (`organization_openid`) REFERENCES `organizations` (`openid`),
  CONSTRAINT `fk_organization_volunteers_volunteer` FOREIGN KEY (`volunteer_openid`) REFERENCES `volunteers` (`openid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `organization_volunteers`
--

LOCK TABLES `organization_volunteers` WRITE;
/*!40000 ALTER TABLE `organization_volunteers` DISABLE KEYS */;
/*!40000 ALTER TABLE `organization_volunteers` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `organizations`
--

DROP TABLE IF EXISTS `organizations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `organizations` (
  `openid` varchar(255) NOT NULL,
  `created_at` datetime(3) NOT NULL COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL COMMENT '更新时间',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT '删除时间',
  `name` longtext NOT NULL,
  `logo` longtext,
  PRIMARY KEY (`openid`),
  KEY `idx_organizations_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_organizations_user` FOREIGN KEY (`openid`) REFERENCES `users` (`openid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `organizations`
--

LOCK TABLES `organizations` WRITE;
/*!40000 ALTER TABLE `organizations` DISABLE KEYS */;
INSERT INTO `organizations` VALUES ('organization','2025-06-30 19:08:13.061','2025-06-30 19:08:13.061',NULL,'asdas',''),('organization_2','2025-06-30 19:08:42.664','2025-06-30 19:08:42.664',NULL,'bbbb','');
/*!40000 ALTER TABLE `organizations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `task_participants`
--

DROP TABLE IF EXISTS `task_participants`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `task_participants` (
  `task_id` bigint NOT NULL COMMENT '主键',
  `volunteer_openid` varchar(255) NOT NULL,
  PRIMARY KEY (`task_id`,`volunteer_openid`),
  KEY `fk_task_participants_volunteer` (`volunteer_openid`),
  CONSTRAINT `fk_task_participants_task` FOREIGN KEY (`task_id`) REFERENCES `tasks` (`id`),
  CONSTRAINT `fk_task_participants_volunteer` FOREIGN KEY (`volunteer_openid`) REFERENCES `volunteers` (`openid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `task_participants`
--

LOCK TABLES `task_participants` WRITE;
/*!40000 ALTER TABLE `task_participants` DISABLE KEYS */;
/*!40000 ALTER TABLE `task_participants` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `tasks`
--

DROP TABLE IF EXISTS `tasks`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `tasks` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `created_at` datetime(3) NOT NULL COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL COMMENT '更新时间',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT '删除时间',
  `title` varchar(40) NOT NULL,
  `start_time` datetime(3) NOT NULL,
  `end_time` datetime(3) NOT NULL,
  `longitude` double NOT NULL,
  `latitude` double NOT NULL,
  `desc` longtext NOT NULL,
  `publisher` longtext NOT NULL,
  `publisher_type` longtext NOT NULL,
  `number` smallint unsigned NOT NULL,
  `already` smallint unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `idx_tasks_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tasks`
--

LOCK TABLES `tasks` WRITE;
/*!40000 ALTER TABLE `tasks` DISABLE KEYS */;
/*!40000 ALTER TABLE `tasks` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `openid` varchar(255) NOT NULL COMMENT 'wx openid',
  `created_at` datetime(3) NOT NULL COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL COMMENT '更新时间',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT '删除时间',
  `name` varchar(40) NOT NULL,
  `phone` varchar(11) NOT NULL,
  `user_type` tinyint unsigned NOT NULL,
  `is_volunteer` tinyint(1) NOT NULL DEFAULT '0',
  `is_elder` tinyint(1) NOT NULL DEFAULT '0',
  `is_monitor` tinyint(1) NOT NULL DEFAULT '0',
  `is_organization` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`openid`),
  KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES ('elder','2025-06-30 19:05:16.697','2025-06-30 19:05:16.697',NULL,'hello-world','12345678901',0,0,1,0,0),('normal1','2025-06-30 19:10:09.768','2025-06-30 19:10:09.768',NULL,'hello-world','12345678901',0,0,0,0,0),('normal2','2025-06-30 19:10:17.771','2025-06-30 19:10:17.771',NULL,'hello-world','12345678901',0,0,0,0,0),('normal3','2025-06-30 19:10:22.631','2025-06-30 19:10:22.631',NULL,'hello-world','12345678901',0,0,0,0,0),('organization','2025-06-30 19:08:10.347','2025-06-30 19:08:10.347',NULL,'hello-world','12345678901',0,0,0,0,1),('organization_2','2025-06-30 19:08:32.225','2025-06-30 19:08:32.225',NULL,'hello-world','12345678901',0,0,0,0,1),('volunteer','2025-06-30 19:06:19.225','2025-06-30 19:06:19.225',NULL,'hello-world','12345678901',0,1,0,0,0);
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `volunteers`
--

DROP TABLE IF EXISTS `volunteers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `volunteers` (
  `openid` varchar(255) NOT NULL,
  `created_at` datetime(3) NOT NULL COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL COMMENT '更新时间',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT '删除时间',
  `school` varchar(40) NOT NULL,
  `clazz` varchar(10) NOT NULL COMMENT '班级',
  `skills` longtext,
  `hours` smallint unsigned NOT NULL DEFAULT '0',
  `start_time` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`openid`),
  KEY `idx_volunteers_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_volunteers_user` FOREIGN KEY (`openid`) REFERENCES `users` (`openid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `volunteers`
--

LOCK TABLES `volunteers` WRITE;
/*!40000 ALTER TABLE `volunteers` DISABLE KEYS */;
INSERT INTO `volunteers` VALUES ('volunteer','2025-06-30 19:06:21.833','2025-06-30 19:06:21.833',NULL,'XJTU','1234','abc',0,NULL);
/*!40000 ALTER TABLE `volunteers` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-06-30 19:34:54
