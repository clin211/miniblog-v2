// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/clin211/miniblog-v2/pkg/db"
	"github.com/spf13/pflag"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

// 帮助信息文本.
const helpText = `Usage: main [flags] arg [arg...]

This is a pflag example.

Flags:
`

// Querier 定义了数据库查询接口.
type Querier interface {
	// FilterWithNameAndRole 按名称和角色查询记录
	FilterWithNameAndRole(name string) ([]gen.T, error)
}

// GenerateConfig 保存代码生成的配置.
type GenerateConfig struct {
	ModelPackagePath string
	GenerateFunc     func(g *gen.Generator)
}

// 预定义的生成配置.
var generateConfigs = map[string]GenerateConfig{
	"mb": {ModelPackagePath: "../../internal/apiserver/model", GenerateFunc: GenerateMiniBlogModels},
}

// 命令行参数.
var (
	addr       = pflag.StringP("addr", "a", "127.0.0.1:3306", "MySQL host address.")
	username   = pflag.StringP("username", "u", "miniblog", "Username to connect to the database.")
	password   = pflag.StringP("password", "p", "miniblog1234", "Password to use when connecting to the database.")
	database   = pflag.StringP("db", "d", "miniblog_v2", "Database name to connect to.")
	modelPath  = pflag.String("model-pkg-path", "", "Generated model code's package name.")
	components = pflag.StringSlice("component", []string{"mb"}, "Generated model code's for specified component.")
	help       = pflag.BoolP("help", "h", false, "Show this help message.")
)

func main() {
	// 设置自定义的使用说明函数
	pflag.Usage = func() {
		fmt.Printf("%s", helpText)
		pflag.PrintDefaults()
	}
	pflag.Parse()

	// 如果设置了帮助标志，则显示帮助信息并退出
	if *help {
		pflag.Usage()
		return
	}

	// 初始化数据库连接
	dbInstance, err := initializeDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 处理组件并生成代码
	for _, component := range *components {
		processComponent(component, dbInstance)
	}
}

// initializeDatabase 创建并返回一个数据库连接.
func initializeDatabase() (*gorm.DB, error) {
	dbOptions := &db.MySQLOptions{
		Addr:     *addr,
		Username: *username,
		Password: *password,
		Database: *database,
	}

	// 创建并返回数据库连接
	return db.NewMySQL(dbOptions)
}

// processComponent 处理单个组件以生成代码.
func processComponent(component string, dbInstance *gorm.DB) {
	config, ok := generateConfigs[component]
	if !ok {
		log.Printf("Component '%s' not found in configuration. Skipping.", component)
		return
	}

	// 解析模型包路径
	modelPkgPath := resolveModelPackagePath(config.ModelPackagePath)

	// 创建生成器实例
	generator := createGenerator(modelPkgPath)
	generator.UseDB(dbInstance)

	// 应用自定义生成器选项
	applyGeneratorOptions(generator)

	// 使用指定的函数生成模型
	config.GenerateFunc(generator)

	// 执行代码生成
	generator.Execute()
}

// resolveModelPackagePath 确定模型生成的包路径.
func resolveModelPackagePath(defaultPath string) string {
	if *modelPath != "" {
		return *modelPath
	}
	absPath, err := filepath.Abs(defaultPath)
	if err != nil {
		log.Printf("Error resolving path: %v", err)
		return defaultPath
	}
	return absPath
}

// createGenerator 初始化并返回一个新的生成器实例.
func createGenerator(packagePath string) *gen.Generator {
	return gen.NewGenerator(gen.Config{
		Mode:              gen.WithDefaultQuery | gen.WithQueryInterface | gen.WithoutContext,
		ModelPkgPath:      packagePath,
		WithUnitTest:      true,
		FieldNullable:     true,  // 对于数据库中可空的字段，使用指针类型。
		FieldSignable:     false, // 禁用无符号属性以提高兼容性。
		FieldWithIndexTag: false, // 不包含 GORM 的索引标签。
		FieldWithTypeTag:  false, // 不包含 GORM 的类型标签。
	})
}

// applyGeneratorOptions 设置自定义生成器选项.
func applyGeneratorOptions(g *gen.Generator) {
	// 为特定字段自定义 GORM 标签
	g.WithOpts(
		gen.FieldGORMTag("created_at", func(tag field.GormTag) field.GormTag {
			tag.Set("default", "current_timestamp")
			return tag
		}),
		gen.FieldGORMTag("updated_at", func(tag field.GormTag) field.GormTag {
			tag.Set("default", "current_timestamp")
			return tag
		}),
	)
}

// GenerateMiniBlogModels 为 miniblog 组件生成模型.
func GenerateMiniBlogModels(g *gen.Generator) {
	// 用户表模型生成
	g.GenerateModelAs(
		"user",
		"UserM",
		gen.FieldIgnore("placeholder"),
		gen.FieldRename("user_id", "UserID"),
		gen.FieldRename("password_updated_at", "PasswordUpdatedAt"),
		gen.FieldRename("email_verified", "EmailVerified"),
		gen.FieldRename("phone_verified", "PhoneVerified"),
		gen.FieldRename("failed_login_attempts", "FailedLoginAttempts"),
		gen.FieldRename("last_login_at", "LastLoginAt"),
		gen.FieldRename("last_login_ip", "LastLoginIP"),
		gen.FieldRename("last_login_device", "LastLoginDevice"),
		gen.FieldRename("is_risk", "IsRisk"),
		gen.FieldRename("register_source", "RegisterSource"),
		gen.FieldRename("register_ip", "RegisterIP"),
		gen.FieldRename("wechat_openid", "WechatOpenID"),
		gen.FieldRename("created_at", "CreatedAt"),
		gen.FieldRename("updated_at", "UpdatedAt"),
		gen.FieldRename("deleted_at", "DeletedAt"),
		gen.FieldGORMTag("user_id", func(tag field.GormTag) field.GormTag {
			tag.Set("uniqueIndex", "uk_user_id")
			return tag
		}),
		gen.FieldGORMTag("username", func(tag field.GormTag) field.GormTag {
			tag.Set("uniqueIndex", "uk_username")
			return tag
		}),
		gen.FieldGORMTag("email", func(tag field.GormTag) field.GormTag {
			tag.Set("uniqueIndex", "uk_email")
			return tag
		}),
		gen.FieldGORMTag("phone", func(tag field.GormTag) field.GormTag {
			tag.Set("uniqueIndex", "uk_phone")
			return tag
		}),
		gen.FieldGORMTag("wechat_openid", func(tag field.GormTag) field.GormTag {
			tag.Set("uniqueIndex", "uk_wechat_openid")
			return tag
		}),
		gen.FieldGORMTag("status", func(tag field.GormTag) field.GormTag {
			tag.Set("index", "idx_status")
			return tag
		}),
		gen.FieldGORMTag("deleted_at", func(tag field.GormTag) field.GormTag {
			tag.Set("index", "idx_deleted_at")
			return tag
		}),
	)

	// 分类表模型生成
	g.GenerateModelAs(
		"category",
		"CategoryM",
		gen.FieldIgnore("placeholder"),
		gen.FieldRename("parent_id", "ParentID"),
		gen.FieldRename("sort_order", "SortOrder"),
		gen.FieldRename("is_active", "IsActive"),
		gen.FieldRename("created_at", "CreatedAt"),
		gen.FieldRename("updated_at", "UpdatedAt"),
		gen.FieldRename("deleted_at", "DeletedAt"),
		gen.FieldGORMTag("parent_id", func(tag field.GormTag) field.GormTag {
			tag.Set("index", "idx_parent_id")
			return tag
		}),
		gen.FieldGORMTag("sort_order", func(tag field.GormTag) field.GormTag {
			tag.Set("index", "idx_sort_order")
			return tag
		}),
	)

	// 标签表模型生成
	g.GenerateModelAs(
		"tag",
		"TagM",
		gen.FieldIgnore("placeholder"),
		gen.FieldRename("created_at", "CreatedAt"),
		gen.FieldRename("updated_at", "UpdatedAt"),
		gen.FieldRename("deleted_at", "DeletedAt"),
		gen.FieldGORMTag("name", func(tag field.GormTag) field.GormTag {
			tag.Set("uniqueIndex", "uk_tag_name")
			return tag
		}),
	)

	// 文章表模型生成 - 注意这里改为 user_id
	g.GenerateModelAs(
		"post",
		"PostM",
		gen.FieldIgnore("placeholder"),
		gen.FieldRename("post_id", "PostID"),
		gen.FieldRename("user_id", "UserID"),
		gen.FieldRename("category_id", "CategoryID"),
		gen.FieldRename("post_type", "PostType"),
		gen.FieldRename("original_author", "OriginalAuthor"),
		gen.FieldRename("original_source", "OriginalSource"),
		gen.FieldRename("original_author_intro", "OriginalAuthorIntro"),
		gen.FieldRename("view_count", "ViewCount"),
		gen.FieldRename("like_count", "LikeCount"),
		gen.FieldRename("published_at", "PublishedAt"),
		gen.FieldRename("created_at", "CreatedAt"),
		gen.FieldRename("updated_at", "UpdatedAt"),
		gen.FieldRename("deleted_at", "DeletedAt"),
		gen.FieldGORMTag("post_id", func(tag field.GormTag) field.GormTag {
			tag.Set("uniqueIndex", "uk_post_id")
			return tag
		}),
		gen.FieldGORMTag("user_id", func(tag field.GormTag) field.GormTag {
			tag.Set("index", "idx_user_id")
			return tag
		}),
		gen.FieldGORMTag("category_id", func(tag field.GormTag) field.GormTag {
			tag.Set("index", "idx_category_id")
			return tag
		}),
		gen.FieldGORMTag("post_type", func(tag field.GormTag) field.GormTag {
			tag.Set("index", "idx_post_type")
			return tag
		}),
		gen.FieldGORMTag("status", func(tag field.GormTag) field.GormTag {
			tag.Set("index", "idx_status")
			return tag
		}),
		gen.FieldGORMTag("deleted_at", func(tag field.GormTag) field.GormTag {
			tag.Set("index", "idx_deleted_at")
			return tag
		}),
	)

	// 文章标签关联表模型生成
	g.GenerateModelAs(
		"post_tag",
		"PostTagM",
		gen.FieldIgnore("placeholder"),
		gen.FieldRename("post_id", "PostID"),
		gen.FieldRename("tag_id", "TagID"),
		gen.FieldRename("created_at", "CreatedAt"),
		gen.FieldRename("updated_at", "UpdatedAt"),
		gen.FieldRename("deleted_at", "DeletedAt"),
		gen.FieldGORMTag("post_id", func(tag field.GormTag) field.GormTag {
			tag.Set("primaryKey", "true")
			tag.Set("index", "idx_post_id")
			return tag
		}),
		gen.FieldGORMTag("tag_id", func(tag field.GormTag) field.GormTag {
			tag.Set("primaryKey", "true")
			tag.Set("index", "idx_tag_id")
			return tag
		}),
	)

	// Casbin 规则表模型生成
	g.GenerateModelAs(
		"casbin_rule",
		"CasbinRuleM",
		gen.FieldRename("ptype", "PType"),
		gen.FieldIgnore("placeholder"),
	)
}
