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
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP() COMMENT '创建时间',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP() COMMENT '更新时间',
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
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP() COMMENT '创建时间',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP() COMMENT '更新时间',
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


-- 分类表
CREATE TABLE category (
    `id` INT AUTO_INCREMENT PRIMARY KEY COMMENT '分类ID',
    `name` VARCHAR(50) NOT NULL COMMENT '分类名称',
    `description` TEXT COMMENT '分类描述',
    `parent_id` INT DEFAULT 0 COMMENT '父分类ID，0表示顶级分类',
    `sort_order` INT DEFAULT 0 COMMENT '排序值',
    `is_active` TINYINT DEFAULT 1 COMMENT '是否激活；1-激活,0-禁用',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP() COMMENT '创建时间',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP() COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL COMMENT '删除时间',
    INDEX idx_parent_id (`parent_id`),
    INDEX idx_sort_order (`sort_order`)
) COMMENT='分类表' ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 插入分类数据
INSERT INTO category (name, description, parent_id, sort_order, is_active) VALUES
-- 顶级分类
('前端开发', '前端技术相关文章，包括 HTML、CSS、JavaScript、TypeScript、Node.js、Vue、React、Webpack、Vite 等前端框架和工具', 0, 10, 1),
('后端开发', '后端技术相关文章，主要包括 Go 语言、Web 框架、API 设计等服务端开发技术', 0, 20, 1),
('云原生与运维', '云原生技术、DevOps、容器化、微服务架构等现代运维和部署技术', 0, 30, 1),
('数据库技术', '数据库设计、SQL 优化、缓存技术等数据存储和管理相关内容', 0, 40, 1),
('项目实战', '完整项目开发经验、架构设计、最佳实践等实战经验分享', 0, 50, 1),
('学习笔记', '技术学习过程中的笔记、总结和心得体会', 0, 60, 1),
('技术分享', '技术趋势分析、工具推荐、开发经验等技术分享内容', 0, 70, 1),

-- 前端开发子分类 (parent_id = 1)
('基础技术', 'HTML、CSS、JavaScript、TypeScript 等前端基础技术', 1, 11, 1),
('前端框架', 'Vue、React、Angular 等主流前端框架', 1, 12, 1),
('构建工具', 'Webpack、Vite、Rollup 等前端构建和打包工具', 1, 13, 1),
('移动开发', 'React Native、Flutter、小程序 等移动端开发技术', 1, 14, 1),

-- 后端开发子分类 (parent_id = 2)
('Go语言', 'Go 语言基础、进阶技巧和最佳实践', 2, 21, 1),
('Web框架', 'Gin、go-zero、Fiber 等 Go Web 开发框架', 2, 22, 1),
('微服务', '微服务架构设计、gRPC、服务治理等', 2, 23, 1),
('API设计', 'RESTful API、GraphQL、接口设计规范等', 2, 24, 1),

-- 云原生与运维子分类 (parent_id = 3)
('容器技术', 'Docker、容器化部署、镜像优化等', 3, 31, 1),
('编排平台', 'Kubernetes、容器编排、集群管理等', 3, 32, 1),
('DevOps', 'CI/CD、自动化部署、持续集成等', 3, 33, 1),
('监控运维', '系统监控、日志分析、性能调优等', 3, 34, 1),

-- 数据库技术子分类 (parent_id = 4)
('关系数据库', 'MySQL、PostgreSQL 等关系型数据库技术', 4, 41, 1),
('NoSQL', 'MongoDB、Redis等非关系型数据库', 4, 42, 1),
('数据库优化', 'SQL优化、索引设计、性能调优等', 4, 43, 1),

-- 项目实战子分类 (parent_id = 5)
('全栈项目', '前后端完整项目开发经验', 5, 51, 1),
('开源项目', '开源项目贡献、项目维护经验', 5, 52, 1),
('架构设计', '系统架构、技术选型、设计模式等', 5, 53, 1),

-- 学习笔记子分类 (parent_id = 6)
('读书笔记', '技术书籍阅读笔记和总结', 6, 61, 1),
('课程学习', '在线课程、培训学习记录', 6, 62, 1),
('问题解决', '开发过程中遇到的问题及解决方案', 6, 63, 1),

-- 技术分享子分类 (parent_id = 7)
('技术趋势', '行业趋势分析、新技术探索', 7, 71, 1),
('工具推荐', '开发工具、效率工具推荐', 7, 72, 1),
('经验总结', '开发经验、团队协作、职业发展等', 7, 73, 1);

-- 标签表
CREATE TABLE tag (
    `id` INT AUTO_INCREMENT PRIMARY KEY COMMENT '标签ID',
    `name` VARCHAR(50) NOT NULL UNIQUE COMMENT '标签名称',
    `color` VARCHAR(7) DEFAULT '#007bff' COMMENT '标签颜色(hex 格式)',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP() COMMENT '创建时间',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP() COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL COMMENT '删除时间'
) COMMENT='标签表' ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 插入标签数据
INSERT INTO tag (name, color) VALUES
-- 前端技术标签
('HTML', '#E34F26'),           -- HTML 官方橙色
('CSS', '#1572B6'),            -- CSS 官方蓝色
('JavaScript', '#F7DF1E'),     -- JavaScript 官方黄色
('TypeScript', '#3178C6'),     -- TypeScript 官方蓝色
('Vue', '#4FC08D'),            -- Vue 官方绿色
('React', '#61DAFB'),          -- React 官方青色
('Node.js', '#339933'),        -- Node.js 官方绿色
('前端框架', '#8B5CF6'),        -- 紫色
('移动端', '#EC4899'),          -- 粉色
('Webpack', '#8DD6F9'),        -- 浅蓝色
('Vite', '#646CFF'),           -- 紫蓝色

-- 后端技术标签（Go 方向）
('Go', '#00ADD8'),             -- Go 官方青色
('Go Web', '#00ADD8'),         -- Go 青色
('Gin', '#00ADD8'),            -- Gin 框架，使用 Go 色系
('go-zero', '#0EA5E9'),        -- go-zero 框架，蓝色系
('Fiber', '#00D9FF'),          -- Fiber 框架，青色系
('GORM', '#00ADD8'),           -- GORM ORM，Go 色系
('微服务', '#6B7280'),          -- 灰色
('gRPC', '#244C5A'),           -- gRPC 深蓝色
('API设计', '#F59E0B'),         -- 橙色

-- 云原生技术标签
('Docker', '#2496ED'),         -- Docker 官方蓝色
('Kubernetes', '#326CE5'),     -- Kubernetes 官方蓝色
('云原生', '#0EA5E9'),          -- 蓝色系
('DevOps', '#22C55E'),         -- 绿色
('CI/CD', '#10B981'),          -- 绿色系
('监控', '#8B5CF6'),            -- 紫色
('日志', '#A855F7'),            -- 紫色系

-- 数据库相关标签
('MySQL', '#4479A1'),          -- MySQL 官方蓝色
('Redis', '#DC382D'),          -- Redis 官方红色
('MongoDB', '#47A248'),        -- MongoDB 官方绿色
('PostgreSQL', '#336791'),     -- PostgreSQL 官方蓝色
('数据库设计', '#7C3AED'),       -- 紫色
('SQL', '#005C84'),            -- 深蓝色
('缓存', '#EF4444'),            -- 红色系

-- 服务器和运维标签
('服务器', '#6B7280'),          -- 灰色
('Linux', '#FCC624'),          -- Linux 金黄色
('Nginx', '#009639'),          -- Nginx 官方绿色
('运维', '#64748B'),            -- 灰蓝色
('安全', '#DC2626'),            -- 红色

-- 项目和学习标签
('项目实战', '#DC2626'),        -- 红色
('算法', '#8B5CF6'),            -- 紫色
('架构设计', '#1E40AF'),        -- 深蓝色
('性能优化', '#059669'),        -- 绿色
('最佳实践', '#7C2D12'),        -- 棕色
('开源项目', '#16A34A'),        -- 绿色
('技术分享', '#0891B2'),        -- 青色
('学习笔记', '#9333EA');        -- 紫色


-- 文章标签关联表（多对多关系）
CREATE TABLE post_tag (
    `post_id` VARCHAR(32) NOT NULL COMMENT '文章ID',
    `tag_id` INT NOT NULL COMMENT '标签ID',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP() COMMENT '创建时间',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP() COMMENT '更新时间',
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
