-- 创建数据库
CREATE DATABASE IF NOT EXISTS blog;
USE blog;

-- 删除已存在的表（按依赖关系逆序删除）
DROP TABLE IF EXISTS post_tags;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS users;

-- 用户表
CREATE TABLE users (
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
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
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
CREATE TABLE categories (
    `id` INT AUTO_INCREMENT PRIMARY KEY COMMENT '分类ID',
    `name` VARCHAR(50) NOT NULL COMMENT '分类名称',
    `description` TEXT COMMENT '分类描述',
    `parent_id` INT DEFAULT 0 COMMENT '父分类ID，0表示顶级分类',
    `sort_order` INT DEFAULT 0 COMMENT '排序值',
    `is_active` TINYINT DEFAULT 1 COMMENT '是否激活；1-激活,0-禁用',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL COMMENT '删除时间',
    INDEX idx_parent_id (`parent_id`),
    INDEX idx_sort_order (`sort_order`)
) COMMENT='分类表' ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 文章表
CREATE TABLE posts (
    `id` BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '主键',
    `post_id` VARCHAR(32) NOT NULL COMMENT '文章ID',
    `title` VARCHAR(200) NOT NULL COMMENT '文章标题',
    `content` LONGTEXT COMMENT '文章内容',
    `cover` VARCHAR(255) COMMENT '文章封面',
    `summary` VARCHAR(500) COMMENT '文章摘要',
    `author_id` VARCHAR(32) NOT NULL COMMENT '作者ID（发布者用户ID）',
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
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL COMMENT '删除时间',

    -- 唯一索引
    UNIQUE KEY uk_post_id (`post_id`),

    -- 外键约束
    FOREIGN KEY (`author_id`) REFERENCES users(`user_id`) ON DELETE CASCADE,
    FOREIGN KEY (`category_id`) REFERENCES categories(`id`) ON DELETE SET NULL,

    -- 基础查询索引（最常用的）
    INDEX idx_author_id (`author_id`),
    INDEX idx_post_type (`post_type`),
    INDEX idx_status (`status`),
    INDEX idx_deleted_at (`deleted_at`)
) COMMENT='文章表' ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 标签表
CREATE TABLE tags (
    `id` INT AUTO_INCREMENT PRIMARY KEY COMMENT '标签ID',
    `name` VARCHAR(50) NOT NULL UNIQUE COMMENT '标签名称',
    `color` VARCHAR(7) DEFAULT '#007bff' COMMENT '标签颜色(hex 格式)',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL COMMENT '删除时间'
) COMMENT='标签表' ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 文章标签关联表（多对多关系）
CREATE TABLE post_tags (
    `post_id` VARCHAR(32) NOT NULL COMMENT '文章ID',
    `tag_id` INT NOT NULL COMMENT '标签ID',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL COMMENT '删除时间',
    PRIMARY KEY (`post_id`, `tag_id`),
    FOREIGN KEY (`post_id`) REFERENCES posts(`post_id`) ON DELETE CASCADE,
    FOREIGN KEY (`tag_id`) REFERENCES tags(`id`) ON DELETE CASCADE
) COMMENT='文章标签关联表' ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 插入示例数据

-- 插入用户数据
INSERT INTO users (`user_id`, `username`, `password`, `email`, `phone`, `gender`, `age`, `avatar`, `status`, `email_verified`, `phone_verified`, `register_source`, `register_ip`, `last_login_at`, `last_login_ip`) VALUES
('user_001', 'zhangsan', '$2a$10$N.zmdr9k7uOCQb376NoUnuTUDg8hL0/Ow8MPEM6kDN5.ynOXrYhga', 'zhangsan@example.com', '13800138001', 1, 25, '/avatars/zhangsan.jpg', 1, 1, 1, 1, '192.168.1.100', NOW(), '192.168.1.100'),
('user_002', 'lisi', '$2a$10$N.zmdr9k7uOCQb376NoUnuTUDg8hL0/Ow8MPEM6kDN5.ynOXrYhga', 'lisi@example.com', '13800138002', 2, 30, '/avatars/lisi.jpg', 1, 1, 0, 2, '192.168.1.101', DATE_SUB(NOW(), INTERVAL 2 DAY), '192.168.1.101'),
('user_003', 'wangwu', '$2a$10$N.zmdr9k7uOCQb376NoUnuTUDg8hL0/Ow8MPEM6kDN5.ynOXrYhga', 'wangwu@example.com', '13800138003', 1, 28, '/avatars/wangwu.jpg', 1, 0, 1, 3, '192.168.1.102', DATE_SUB(NOW(), INTERVAL 7 DAY), '192.168.1.102'),
('user_004', 'test_user', '$2a$10$N.zmdr9k7uOCQb376NoUnuTUDg8hL0/Ow8MPEM6kDN5.ynOXrYhga', 'test@example.com', '13800138004', 3, 22, NULL, 0, 0, 0, 1, '192.168.1.103', NULL, NULL);

-- 插入分类数据
INSERT INTO categories (`name`, `description`, `parent_id`, `sort_order`) VALUES
('技术', '技术相关文章', 0, 1),
('生活', '生活感悟文章', 0, 2),
('Go语言', 'Go语言相关技术文章', 1, 1),
('前端技术', '前端开发相关文章', 1, 2);

-- 插入标签数据
INSERT INTO tags (`name`, `color`) VALUES
('Go', '#00ADD8'),
('GORM', '#FF6B6B'),
('数据库', '#4ECDC4'),
('教程', '#45B7D1');

-- 插入文章数据
INSERT INTO posts (`post_id`, `title`, `content`, `summary`, `author_id`, `category_id`, `post_type`, `status`, `published_at`) VALUES
-- 原创文章
('post_001', 'Gorm Gen 入门教程', '这是一篇关于 Gorm Gen 的详细教程...', '详细介绍了 Gorm Gen 的基础使用方法', 'user_001', 3, 1, 2, NOW()),
('post_002', 'Go 语言并发编程', '本文介绍 Go 语言的并发编程模式...', 'Go 语言并发编程的最佳实践', 'user_001', 3, 1, 2, NOW());

-- 转载文章示例
INSERT INTO posts (`post_id`, `title`, `content`, `summary`, `author_id`, `category_id`, `post_type`, `original_author`, `original_source`, `status`, `published_at`) VALUES
('post_003', '深入理解Go语言内存管理', '本文转载自官方博客，详细解释了Go语言的内存管理机制...', 'Go语言内存管理机制深度解析', 'user_002', 3, 2, 'Go Team', 'https://go.dev/blog/gc-guide', 2, NOW());

-- 投稿文章示例
INSERT INTO posts (`post_id`, `title`, `content`, `summary`, `author_id`, `category_id`, `post_type`, `original_author`, `original_author_intro`, `status`, `published_at`) VALUES
('post_004', 'Vue3 组合式API最佳实践', '这是一篇关于Vue3组合式API的投稿文章...', '分享Vue3组合式API的实际使用经验', 'user_001', 4, 3, '王五', '一名资深前端开发工程师，专注于Vue生态系统', 2, NOW());

-- 插入文章标签关联数据
INSERT INTO post_tags (`post_id`, `tag_id`) VALUES
('post_001', 1), ('post_001', 2), ('post_001', 4),
('post_002', 1), ('post_002', 4),
('post_003', 1), ('post_003', 3),
('post_004', 4);