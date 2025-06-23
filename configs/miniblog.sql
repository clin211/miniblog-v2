-- Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
-- Use of this source code is governed by a MIT style
-- license that can be found in the LICENSE file. The original repo for
-- this file is https://github.com/clin211/miniblog-v2.git.

-- 创建数据库
CREATE DATABASE IF NOT EXISTS miniblog_v2 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
USE miniblog_v2;

-- 删除已存在的表（按依赖关系逆序删除）
DROP TABLE IF EXISTS post_tag;
DROP TABLE IF EXISTS post;
DROP TABLE IF EXISTS tag;
DROP TABLE IF EXISTS category;
DROP TABLE IF EXISTS user;
DROP TABLE IF EXISTS casbin_rule;

-- 用户表
CREATE TABLE user (
    `id` BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '用户ID',
    `user_id` VARCHAR(32) NOT NULL COMMENT '用户ID',
    `age` INT COMMENT '年龄',
    `avatar` VARCHAR(255) COMMENT '头像URL',
    `username` VARCHAR(100) NOT NULL COMMENT '用户名',
    `password` VARCHAR(255) NOT NULL COMMENT '密码',
    `password_updated_at` TIMESTAMP COMMENT '密码更新时间',
    `email` VARCHAR(100) NOT NULL UNIQUE COMMENT '邮箱',
    `email_verified` TINYINT DEFAULT 0 COMMENT '邮箱是否已验证；1-已验证,0-未验证',
    `phone` VARCHAR(20) NOT NULL UNIQUE COMMENT '手机号',
    `phone_verified` TINYINT DEFAULT 0 COMMENT '手机号是否已验证；1-已验证,0-未验证',
    `gender` TINYINT DEFAULT 0 COMMENT '性别：0-未设置，1-男，2-女，3-其他',
    `status` TINYINT DEFAULT 1 COMMENT '状态：1-正常，0-禁用',
    `failed_login_attempts` INT DEFAULT 0 COMMENT '失败登录次数，超过5次则锁定账户，登录成功后重置',
    `last_login_at` TIMESTAMP NULL COMMENT '最后登录时间',
    `last_login_ip` VARCHAR(45) COMMENT '最后登录IP',
    `last_login_device` VARCHAR(100) COMMENT '最后登录设备',
    `is_risk` TINYINT DEFAULT 0 COMMENT '是否为风险用户；1-是,0-否',
    `register_source` TINYINT DEFAULT 1 COMMENT '注册来源：1-web，2-app，3-wechat，4-qq，5-github，6-google',
    `register_ip` VARCHAR(45) COMMENT '注册IP',
    `wechat_openid` VARCHAR(100) COMMENT '微信OpenID',
    `created_at` TIMESTAMP DEFAULT current_timestamp() COMMENT '创建时间',
    `updated_at` TIMESTAMP DEFAULT current_timestamp() ON UPDATE current_timestamp() COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL COMMENT '删除时间',

    -- 唯一索引（保证数据唯一性，必须保留）
    UNIQUE KEY uk_user_id (`user_id`),
    UNIQUE KEY uk_username (`username`),
    UNIQUE KEY uk_email (`email`),
    UNIQUE KEY uk_phone (`phone`),
    UNIQUE KEY uk_wechat_openid (`wechat_openid`),

    -- 基础查询索引（最常用的）
    INDEX idx_status (`status`),
    INDEX idx_deleted_at (`deleted_at`)
) COMMENT='用户表' ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 分类表
CREATE TABLE category (
    `id` INT AUTO_INCREMENT PRIMARY KEY COMMENT '分类ID',
    `name` VARCHAR(50) NOT NULL COMMENT '分类名称',
    `description` TEXT COMMENT '分类描述',
    `parent_id` INT DEFAULT 0 COMMENT '父分类ID，0表示顶级分类',
    `sort_order` INT DEFAULT 0 COMMENT '排序值',
    `is_active` TINYINT DEFAULT 1 COMMENT '是否激活；1-激活,0-禁用',
    `created_at` TIMESTAMP DEFAULT current_timestamp() COMMENT '创建时间',
    `updated_at` TIMESTAMP DEFAULT current_timestamp() ON UPDATE current_timestamp() COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL COMMENT '删除时间',
    INDEX idx_parent_id (`parent_id`),
    INDEX idx_sort_order (`sort_order`)
) COMMENT='分类表' ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 文章表
CREATE TABLE post (
    `id` BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '主键',
    `post_id` VARCHAR(32) NOT NULL COMMENT '文章ID',
    `title` VARCHAR(200) NOT NULL COMMENT '文章标题',
    `content` LONGTEXT COMMENT '文章内容',
    `cover` VARCHAR(255) COMMENT '文章封面',
    `summary` VARCHAR(500) COMMENT '文章摘要',
    `user_id` VARCHAR(32) NOT NULL COMMENT '用户ID（发布者用户ID）',
    `category_id` INT COMMENT '分类ID',
    `post_type` TINYINT DEFAULT 1 COMMENT '文章类型：1-原创，2-转载，3-投稿',
    `original_author` VARCHAR(100) DEFAULT NULL COMMENT '原作者姓名（转载时/投稿时使用）',
    `original_source` VARCHAR(500) DEFAULT NULL COMMENT '原文链接或来源（转载时/投稿时使用）',
    `original_author_intro` VARCHAR(500) DEFAULT NULL COMMENT '原作者简介（转载时/投稿时使用）',
    `position` INT DEFAULT 0 COMMENT '文章排序，0-默认排序，1-置顶，数字越大越靠前',
    `view_count` INT DEFAULT 0 COMMENT '阅读次数',
    `like_count` INT DEFAULT 0 COMMENT '点赞数',
    `status` TINYINT DEFAULT 1 COMMENT '文章状态：1-草稿，2-已发布，3-已归档',
    `published_at` TIMESTAMP NULL COMMENT '发布时间',
    `created_at` TIMESTAMP DEFAULT current_timestamp() COMMENT '创建时间',
    `updated_at` TIMESTAMP DEFAULT current_timestamp() ON UPDATE current_timestamp() COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL COMMENT '删除时间',

    -- 唯一索引
    UNIQUE KEY uk_post_id (`post_id`),

    -- 基础查询索引（最常用的）
    INDEX idx_user_id (`user_id`),
    INDEX idx_category_id (`category_id`),
    INDEX idx_post_type (`post_type`),
    INDEX idx_status (`status`),
    INDEX idx_deleted_at (`deleted_at`)
) COMMENT='文章表' ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 标签表
CREATE TABLE tag (
    `id` INT AUTO_INCREMENT PRIMARY KEY COMMENT '标签ID',
    `name` VARCHAR(50) NOT NULL UNIQUE COMMENT '标签名称',
    `color` VARCHAR(7) DEFAULT '#007bff' COMMENT '标签颜色(hex 格式)',
    `created_at` TIMESTAMP DEFAULT current_timestamp() COMMENT '创建时间',
    `updated_at` TIMESTAMP DEFAULT current_timestamp() ON UPDATE current_timestamp() COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL COMMENT '删除时间'
) COMMENT='标签表' ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 文章标签关联表（多对多关系）
CREATE TABLE post_tag (
    `post_id` VARCHAR(32) NOT NULL COMMENT '文章ID',
    `tag_id` INT NOT NULL COMMENT '标签ID',
    `created_at` TIMESTAMP DEFAULT current_timestamp() COMMENT '创建时间',
    `updated_at` TIMESTAMP DEFAULT current_timestamp() ON UPDATE current_timestamp() COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL COMMENT '删除时间',
    PRIMARY KEY (`post_id`, `tag_id`),
    INDEX idx_post_id (`post_id`),
    INDEX idx_tag_id (`tag_id`)
) COMMENT='文章标签关联表' ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `casbin_rule` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `ptype` varchar(100) DEFAULT NULL,
  `v0` varchar(100) DEFAULT NULL,
  `v1` varchar(100) DEFAULT NULL,
  `v2` varchar(100) DEFAULT NULL,
  `v3` varchar(100) DEFAULT NULL,
  `v4` varchar(100) DEFAULT NULL,
  `v5` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_casbin_rule` (`ptype`,`v0`,`v1`,`v2`,`v3`,`v4`,`v5`)
) COMMENT='casbin_rule' ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1 COLLATE=latin1_swedish_ci;
