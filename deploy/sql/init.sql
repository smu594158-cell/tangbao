-- 创建数据库
CREATE DATABASE IF NOT EXISTS `hztour_db` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `hztour_db`;

-- 1. 用户表
CREATE TABLE IF NOT EXISTS `users` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `username` VARCHAR(64) NOT NULL COMMENT '登录账号',
    `password_hash` VARCHAR(255) NOT NULL COMMENT '哈希加盐密码',
    `nickname` VARCHAR(64) DEFAULT '' COMMENT '用户昵称',
    `role` TINYINT DEFAULT 1 COMMENT '角色: 1-普通用户, 9-管理员',
    `status` TINYINT DEFAULT 1 COMMENT '状态: 1-正常, 0-禁用',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL DEFAULT NULL,
    UNIQUE KEY `idx_username` (`username`),
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户基础信息表';

-- 2. 景区表
CREATE TABLE IF NOT EXISTS `attractions` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `name` VARCHAR(128) NOT NULL COMMENT '景区名称',
    `description` TEXT COMMENT '景区简介',
    `location_lng` DECIMAL(10, 6) COMMENT '经度',
    `location_lat` DECIMAL(10, 6) COMMENT '纬度',
    `address` VARCHAR(255) COMMENT '详细地址',
    `heat_level` INT DEFAULT 0 COMMENT '热度级别',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL DEFAULT NULL,
    INDEX `idx_name` (`name`),
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='杭州各景区基础信息表';

-- 3. 聊天记录表
CREATE TABLE IF NOT EXISTS `chat_histories` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '关联用户ID',
    `session_id` VARCHAR(64) NOT NULL COMMENT '会话ID，用于多轮上下文',
    `role` ENUM('user', 'assistant', 'system') NOT NULL COMMENT '发言角色',
    `content` TEXT NOT NULL COMMENT '对话内容',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '记录时间，策略: 仅保留30天',
    INDEX `idx_user_session` (`user_id`, `session_id`),
    INDEX `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户与AI多轮对话历史表';

-- 4. 生成文本表
CREATE TABLE IF NOT EXISTS `generated_texts` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `attraction_id` BIGINT UNSIGNED NOT NULL COMMENT '关联景区ID',
    `source_url` VARCHAR(512) COMMENT '数据来源(如小红书URL)',
    `original_content` TEXT COMMENT '爬取的原始语料',
    `generated_content` TEXT NOT NULL COMMENT 'AI生成的景点介绍文本',
    `word_count` INT NOT NULL COMMENT '生成字数',
    `plagiarism_score` DECIMAL(5,2) COMMENT '原创度检测评分',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX `idx_attraction` (`attraction_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='景点/人物介绍文本生成表';

