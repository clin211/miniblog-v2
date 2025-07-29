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

-- 插入的分类数据
INSERT INTO category (name, description, parent_id, sort_order, is_active) VALUES
-- 顶级分类
('前端开发', '前端技术相关文章，包括 React、Next.js、JavaScript、TypeScript、CSS 等现代前端技术栈', 0, 100, 1),
('后端开发', '后端技术相关文章，主要包括 Go 语言生态、Web 框架、微服务架构等服务端开发技术', 0, 200, 1),
('算法与数据结构', '计算机科学基础，包括各种数据结构实现、算法分析、复杂度分析等核心内容', 0, 300, 1),
('系统设计', '软件架构设计、设计模式、微服务架构、系统优化等高级技术主题', 0, 400, 1),
('开发工具', '开发环境配置、编辑器使用、版本控制、Shell 脚本等提升开发效率的工具和技巧', 0, 500, 1),
('安全技术', 'Web 安全、网络安全、安全防护等信息安全相关技术内容', 0, 600, 1),
('学习笔记', '技术学习过程中的笔记、总结、心得体会和个人思考', 0, 700, 1),
('面试准备', '技术面试相关内容，包括面试题解析、求职经验、技术考察重点等', 0, 800, 1),

-- 前端开发子分类 (parent_id = 1)
('React 生态', 'React 框架及其生态系统，包括状态管理、路由、组件库等', 1, 101, 1),
('Next.js 框架', 'Next.js 全栈框架的深入学习，包括 SSR、SSG、路由、优化等', 1, 102, 1),
('JavaScript 基础', 'JavaScript 语言基础、ES6+特性、异步编程等核心概念', 1, 103, 1),
('TypeScript', 'TypeScript 类型系统、高级特性、项目实践等', 1, 104, 1),
('CSS 技术', 'CSS 布局、动画、预处理器、现代 CSS 特性等样式技术', 1, 105, 1),
('前端工程化', '前端项目构建、代码质量、性能优化、最佳实践等工程化内容', 1, 106, 1),
('Node.js 开发', 'Node.js 运行时环境、服务端 JavaScript、npm 生态、后端 API 开发等', 1, 107, 1),

-- 后端开发子分类 (parent_id = 2)
('Go 语言基础', 'Go 语言语法、特性、标准库、并发编程等基础知识', 2, 201, 1),
('Go Web 开发', 'Gin、go-zero、GoFrame 等 Go Web 框架的使用和实践', 2, 202, 1),
('微服务架构', '微服务设计原则、gRPC、服务治理、分布式系统等', 2, 203, 1),
('API 设计', 'RESTful API、GraphQL、接口设计规范、文档管理等', 2, 204, 1),
('Node.js 后端', 'Express、Koa、Nest.js 等 Node.js 后端框架及服务端开发', 2, 205, 1),

-- 算法与数据结构子分类 (parent_id = 3)
('基础数据结构', '数组、链表、栈、队列、树等基本数据结构的实现和应用', 3, 301, 1),
('高级数据结构', '堆、图、哈希表、并查集等复杂数据结构', 3, 302, 1),
('算法分析', '时间复杂度、空间复杂度、算法优化等理论分析', 3, 303, 1),
('算法实现', '排序、搜索、动态规划、贪心等经典算法的多语言实现', 3, 304, 1),

-- 系统设计子分类 (parent_id = 4)
('设计模式', '常用设计模式的理论学习和实际应用', 4, 401, 1),
('架构设计', '系统架构、技术选型、扩展性设计等宏观设计', 4, 402, 1),
('性能优化', '系统性能分析、优化策略、缓存设计等', 4, 403, 1),
('分布式系统', '分布式架构、一致性、容错性等分布式系统核心概念', 4, 404, 1),

-- 开发工具子分类 (parent_id = 5)
('编辑器配置', 'VS Code、Vim 等编辑器的配置和使用技巧', 5, 501, 1),
('Shell 脚本', 'Shell 编程、命令行工具、自动化脚本等', 5, 502, 1),
('版本控制', 'Git 使用技巧、工作流程、团队协作等', 5, 503, 1),
('开发环境', '开发环境搭建、容器化、虚拟化等环境管理', 5, 504, 1),

-- 安全技术子分类 (parent_id = 6)
('Web 安全基础', 'Web 安全基本概念、攻防原理、安全意识等', 6, 601, 1),
('常见漏洞', 'XSS、CSRF、SQL 注入等常见 Web 漏洞的原理和防护', 6, 602, 1),
('安全实践', '安全编码规范、安全测试、渗透测试等实践内容', 6, 603, 1),

-- 学习笔记子分类 (parent_id = 7)
('技术总结', '技术学习过程中的总结和归纳', 7, 701, 1),
('项目实践', '项目开发过程中的经验和教训', 7, 702, 1),
('问题解决', '开发过程中遇到的问题及解决方案记录', 7, 703, 1),
('个人思考', '技术发展趋势、职业规划等个人思考和感悟', 7, 704, 1),

-- 面试准备子分类 (parent_id = 8)
('前端面试', '前端技术面试题目、考察重点、面试技巧等', 8, 801, 1),
('后端面试', '后端技术面试、系统设计、编程题等', 8, 802, 1),
('算法面试', '算法题解析、刷题心得、面试算法技巧等', 8, 803, 1),
('面试经验', '面试流程、求职经验、职业发展建议等', 8, 804, 1);

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
-- 前端框架标签 (蓝色系)
('Next.js', '#0070f3'),
('React', '#61dafb'),
('Vue', '#4fc08d'),
('Angular', '#dd0031'),

-- 前端基础技术标签 (绿色系)
('JavaScript', '#f7df1e'),
('TypeScript', '#3178c6'),
('HTML', '#e34f26'),
('CSS', '#1572b6'),
('Sass', '#cc6699'),
('Less', '#1d365d'),

-- 前端工具和库标签 (紫色系)
('Webpack', '#8dd6f9'),
('Vite', '#646cff'),
('Rollup', '#ec4a3f'),
('ESLint', '#4b32c3'),
('Prettier', '#f7b93e'),
('Babel', '#f9dc3e'),

-- 状态管理标签 (橙色系)
('Zustand', '#ff6b35'),
('Redux', '#764abc'),
('MobX', '#ff9955'),
('TanStack Query', '#ff4154'),
('SWR', '#000000'),

-- 后端技术标签 (红色系)
('Go', '#00add8'),
('Gin', '#00add8'),
('go-zero', '#00add8'),
('GoFrame', '#00add8'),
('gRPC', '#244c5a'),
('Protobuf', '#4285f4'),

-- 数据库标签 (青色系)
('MySQL', '#4479a1'),
('PostgreSQL', '#336791'),
('MongoDB', '#47a248'),
('Redis', '#dc382d'),
('SQLite', '#003b57'),

-- 云原生和运维标签 (灰色系)
('Docker', '#2496ed'),
('Kubernetes', '#326ce5'),
('CI/CD', '#6c757d'),
('DevOps', '#495057'),
('Nginx', '#009639'),

-- 算法和数据结构标签 (黄色系)
('算法', '#ffc107'),
('数据结构', '#fd7e14'),
('复杂度分析', '#e83e8c'),
('排序算法', '#6f42c1'),
('搜索算法', '#20c997'),

-- 设计模式标签 (粉色系)
('设计模式', '#e91e63'),
('策略模式', '#f06292'),
('观察者模式', '#ba68c8'),
('工厂模式', '#9575cd'),
('单例模式', '#7986cb'),

-- 开发工具标签 (棕色系)
('VS Code', '#007acc'),
('Vim', '#019733'),
('Git', '#f05032'),
('Shell', '#89e051'),
('Linux', '#fcc624'),

-- 性能优化标签 (深蓝色系)
('性能优化', '#0d47a1'),
('缓存', '#1565c0'),
('SEO', '#1976d2'),
('图像优化', '#1e88e5'),
('代码分割', '#2196f3'),

-- 安全相关标签 (深红色系)
('Web安全', '#b71c1c'),
('XSS', '#c62828'),
('CSRF', '#d32f2f'),
('SQL注入', '#f44336'),

-- 架构设计标签 (深绿色系)
('微服务', '#2e7d32'),
('RESTful API', '#388e3c'),
('GraphQL', '#43a047'),
('系统架构', '#4caf50'),
('分布式', '#66bb6a'),

-- 测试相关标签 (深紫色系)
('单元测试', '#4a148c'),
('集成测试', '#6a1b9a'),
('E2E测试', '#7b1fa2'),
('Jest', '#8e24aa'),
('Cypress', '#9c27b0'),

-- 移动开发标签 (青绿色系)
('React Native', '#02569b'),
('Flutter', '#61dafb'),
('小程序', '#07c160'),
('PWA', '#5a0fc8'),

-- Node.js 核心技术标签 (绿色系)
('Node.js', '#339933'),
('npm', '#cb3837'),
('yarn', '#2c8ebb'),
('pnpm', '#f69220'),

-- Node.js 后端框架标签 (深绿色系)
('Express', '#000000'),
('Koa', '#33333d'),
('Nest.js', '#e0234e'),
('Fastify', '#000000'),

-- Node.js 工具和库标签 (蓝绿色系)
('Nodemon', '#76d04b'),
('PM2', '#2b037a'),
('Socket.io', '#010101'),
('Passport', '#34e27a'),
('Multer', '#ff6600'),

-- Node.js 数据库相关标签 (青色系)
('Mongoose', '#880000'),
('Sequelize', '#52b0e7'),
('Prisma', '#2d3748'),
('TypeORM', '#fe0803'),

-- Node.js 构建和部署标签 (橙色系)
('esbuild', '#ffcf00'),
('Serverless', '#fd5750'),

-- Node.js 实时通信标签 (红色系)
('WebSocket', '#010101'),
('Server-Sent Events', '#ff4444'),
('GraphQL Subscriptions', '#e10098'),

-- Node.js 微服务标签 (深蓝色系)
('Microservices', '#0066cc'),
('API Gateway', '#ff9900'),
('Service Mesh', '#326ce5'),

-- 学习和分享标签 (暖色系)
('学习笔记', '#ff5722'),
('最佳实践', '#ff7043'),
('经验分享', '#ff8a65'),
('问题解决', '#ffab91'),
('技术趋势', '#ffcc02'),

-- 项目类型标签 (中性色系)
('全栈项目', '#607d8b'),
('开源项目', '#78909c'),
('个人项目', '#90a4ae'),
('企业项目', '#b0bec5'),

-- 渲染模式标签 (特殊色系)
('SSR', '#00bcd4'),
('SSG', '#009688'),
('CSR', '#4db6ac'),
('ISR', '#26a69a'),

-- 国际化和本地化标签
('国际化', '#795548'),
('i18n', '#8d6e63'),
('本地化', '#a1887f'),

-- 内容管理标签
('Markdown', '#000000'),
('CMS', '#6c757d'),
('内容管理', '#495057'),

-- 面试相关标签
('面试题', '#dc3545'),
('技术面试', '#c82333'),
('求职', '#bd2130');

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
