package main

import (
	"app/cmd/gorm/methods"
	"app/cmd/gorm/tags"
	"app/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

func main() {

	// 初始化数据库
	dbConfig, err := config.GetDBConfig("mysql")
	if err != nil {
		panic(err)
	}

	mysqlDb, err := gorm.Open(mysql.Open(dbConfig.GenerateDsn()))
	if err != nil {
		panic("mysql connect failed: " + err.Error())
	}

	// 公共属性
	comOpts := []gen.ModelOpt{
		// 自动时间戳字段属性
		gen.FieldGORMTag("create_time", tags.CreateField),
		gen.FieldGORMTag("update_time", tags.UpdateField),

		// 软删除字段属性
		gen.FieldType("delete_time", "soft_delete.DeletedAt"),

		// Json序列化
		gen.WithMethod(methods.JsonMethod{}),
	}

	g := gen.NewGenerator(gen.Config{
		OutPath:      "./internal/orm/query",
		ModelPkgPath: "model",
		Mode:         gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	g.UseDB(mysqlDb)

	// 生成所有表
	//g.ApplyBasic(g.GenerateAllTable(comOpts...)...)

	var allTables []any

	// App
	appModel := g.GenerateModel("app", comOpts...)
	allTables = append(allTables, appModel)

	// Role
	roleModel := g.GenerateModel("role", append(comOpts,
		gen.FieldRelate(field.HasOne, "App", appModel, &field.RelateConfig{
			GORMTag: field.GormTag{
				"foreignKey": {"app_id"},
				"references": {"id"},
			},
		}),
	)...)
	allTables = append(allTables, roleModel)

	// UserRole
	userRoleModel := g.GenerateModel("user_role", append(comOpts,
		gen.FieldRelate(field.HasOne, "Role", roleModel, &field.RelateConfig{
			GORMTag: field.GormTag{
				"foreignKey": {"role_id"},
				"references": {"id"},
			},
		}),
	)...)
	allTables = append(allTables, userRoleModel)

	// User
	userModel := g.GenerateModel("user", append(comOpts,
		gen.FieldRelate(field.HasMany, "UserRole", userRoleModel, &field.RelateConfig{
			GORMTag: field.GormTag{
				"foreignKey": {"user_id"},
				"references": {"id"},
			},
		}),
	)...)
	allTables = append(allTables, userModel)

	// MenuApi
	menuApiModel := g.GenerateModel("menu_api", comOpts...)
	allTables = append(allTables, menuApiModel)

	// Menu
	menuModel := g.GenerateModel("menu", append(comOpts,
		gen.FieldRelate(field.HasOne, "App", appModel, &field.RelateConfig{
			GORMTag: field.GormTag{
				"foreignKey": {"app_id"},
				"references": {"id"},
			},
		}),
		gen.FieldRelate(field.HasMany, "MenuApi", menuApiModel, &field.RelateConfig{
			GORMTag: field.GormTag{
				"foreignKey": {"menu_id"},
				"references": {"id"},
			},
		}),
	)...)
	allTables = append(allTables, menuModel)

	// Log/LogType
	logTypeModel := g.GenerateModel("log_type", comOpts...)
	allTables = append(allTables, logTypeModel)
	logModel := g.GenerateModel("log", append(comOpts,
		gen.FieldRelate(field.HasOne, "LogType", userModel, &field.RelateConfig{
			GORMTag: field.GormTag{
				"foreignKey": {"type_id"},
				"references": {"id"},
			},
		}),
	)...)
	allTables = append(allTables, logModel)

	// File/FileStorage
	fileModel := g.GenerateModel("file", comOpts...)
	allTables = append(allTables, fileModel)
	fileStorageModel := g.GenerateModel("file_storage", comOpts...)
	allTables = append(allTables, fileStorageModel)

	// 生成指定表
	g.ApplyBasic(allTables...)

	// 执行生成
	g.Execute()
}
