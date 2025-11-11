package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"

	"app/cmd/gorm/methods"
	"app/cmd/gorm/tags"
	"app/internal/config"
)

func main() {

	// 初始化数据库
	dbConfig, err := config.GetDBConfig("main")
	if err != nil {
		panic(err)
	}

	mainDB, err := gorm.Open(mysql.Open(dbConfig.GenerateDsn()))
	if err != nil {
		panic("database connect failed: " + err.Error())
	}

	// 公共属性
	comOpts := []gen.ModelOpt{
		// 自动时间戳字段属性
		gen.FieldGORMTag("created_at", tags.CreateField),
		gen.FieldGORMTag("updated_at", tags.UpdateField),

		// 软删除字段属性
		gen.FieldType("deleted_at", "gorm.DeletedAt"),

		// Json序列化
		gen.WithMethod(methods.JsonMethod{}),
	}

	g := gen.NewGenerator(gen.Config{
		OutPath:      "./internal/orm/query",
		ModelPkgPath: "model",
		Mode:         gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	g.UseDB(mainDB)

	// 生成所有表
	//g.ApplyBasic(g.GenerateAllTable(comOpts...)...)

	var allTables []any

	// 后台基础模型
	// App
	appModel := g.GenerateModel("app", comOpts...)
	allTables = append(allTables, appModel)

	// Role
	roleModel := g.GenerateModel("role", append(comOpts,
		gen.FieldRelate(field.HasOne, "App", appModel, &field.RelateConfig{
			GORMTag: field.GormTag{"foreignKey": {"app_id"}, "references": {"id"}},
		}),
	)...)
	allTables = append(allTables, roleModel)

	// UserRole
	userRoleModel := g.GenerateModel("user_role", append(comOpts,
		gen.FieldRelate(field.HasOne, "Role", roleModel, &field.RelateConfig{
			GORMTag: field.GormTag{"foreignKey": {"role_id"}, "references": {"id"}},
		}),
	)...)
	allTables = append(allTables, userRoleModel)

	// User
	userModel := g.GenerateModel("user", append(comOpts,
		gen.FieldRelate(field.HasMany, "UserRole", userRoleModel, &field.RelateConfig{
			GORMTag: field.GormTag{"foreignKey": {"user_id"}, "references": {"id"}},
		}),
	)...)
	allTables = append(allTables, userModel)

	// MenuApi
	menuApiModel := g.GenerateModel("menu_api", comOpts...)
	allTables = append(allTables, menuApiModel)

	// Menu
	menuModel := g.GenerateModel("menu", append(comOpts,
		gen.FieldRelate(field.HasMany, "MenuApi", menuApiModel, &field.RelateConfig{
			GORMTag: field.GormTag{"foreignKey": {"menu_id"}, "references": {"id"}},
		}),
	)...)
	allTables = append(allTables, menuModel)

	// Log/LogType
	logTypeModel := g.GenerateModel("log_type", comOpts...)
	logModel := g.GenerateModel("log", append(comOpts,
		gen.FieldRelate(field.HasOne, "LogType", userModel, &field.RelateConfig{
			GORMTag: field.GormTag{"foreignKey": {"type_id"}, "references": {"id"}},
		}),
	)...)
	allTables = append(allTables, logTypeModel, logModel)

	// File/FileStorage/FileMenu
	fileModel := g.GenerateModel("file", comOpts...)
	fileMenuModel := g.GenerateModel("file_menu", comOpts...)
	fileStorageModel := g.GenerateModel("file_storage", comOpts...)
	allTables = append(allTables, fileModel, fileStorageModel, fileMenuModel)

	// Config
	configModel := g.GenerateModel("config", comOpts...)
	allTables = append(allTables, configModel)

	// Task
	taskModel := g.GenerateModel("task", comOpts...)
	allTables = append(allTables, taskModel)

	// 生成指定表
	g.ApplyBasic(allTables...)

	// 执行生成
	g.Execute()
}
