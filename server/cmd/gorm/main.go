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

	g := gen.NewGenerator(gen.Config{
		OutPath:      "./internal/orm/query",
		ModelPkgPath: "model",
		Mode:         gen.WithoutContext | gen.WithDefaultQuery, // generate mode

		// 如果你希望为可为null的字段生成属性为指针类型, 设置 FieldNullable 为 true
		FieldNullable: true,
		// 如果你希望在 `Create` API 中为字段分配默认值, 设置 FieldCoverable 为 true, 参考: https://gorm.io/docs/create.html#Default-Values
		//FieldCoverable: true,
		// 如果你希望生成无符号整数类型字段, 设置 FieldSignable 为 true
		FieldSignable: true,
		// 如果你希望从数据库生成索引标记, 设置 FieldWithIndexTag 为 true
		//FieldWithIndexTag: true,
		// 如果你希望从数据库生成类型标记, 设置 FieldWithTypeTag 为 true
		FieldWithTypeTag: true,
		// 如果你需要对查询代码进行单元测试, 设置 WithUnitTest 为 true
		//WithUnitTest: true,
	})

	g.UseDB(mainDB)

	// 公共模型选项
	comOpts := []gen.ModelOpt{
		// 自动时间戳字段属性
		gen.FieldGORMTag("created_at", tags.CreateField),
		gen.FieldGORMTag("updated_at", tags.UpdateField),

		// 软删除字段属性
		gen.FieldType("deleted_at", "gorm.DeletedAt"),

		// Json序列化
		gen.WithMethod(methods.JsonMethod{}),
	}
	g.WithOpts(comOpts...)

	// 生成所有表
	//g.ApplyBasic(g.GenerateAllTable(comOpts...)...)

	var allTables []any

	// 后台基础模型
	// App
	appModel := g.GenerateModel("app")
	allTables = append(allTables, appModel)

	// Role, UserRole
	roleModel := g.GenerateModel("role", gen.FieldRelate(field.HasOne, "App", appModel, &field.RelateConfig{
		GORMTag: field.GormTag{"foreignKey": {"app_id"}, "references": {"id"}},
	}))
	userRoleModel := g.GenerateModel("user_role",
		gen.FieldRelate(field.HasOne, "Role", roleModel, &field.RelateConfig{
			GORMTag: field.GormTag{"foreignKey": {"role_id"}, "references": {"id"}},
		}),
	)
	allTables = append(allTables, roleModel, userRoleModel)

	// User
	userModel := g.GenerateModel("user",
		gen.FieldRelate(field.HasMany, "UserRole", userRoleModel, &field.RelateConfig{
			GORMTag: field.GormTag{"foreignKey": {"user_id"}, "references": {"id"}},
		}),
	)
	allTables = append(allTables, userModel)

	// MenuApi
	menuApiModel := g.GenerateModel("menu_api")
	allTables = append(allTables, menuApiModel)

	// Menu
	menuModel := g.GenerateModel("menu",
		gen.FieldRelate(field.HasMany, "MenuApi", menuApiModel, &field.RelateConfig{
			GORMTag: field.GormTag{"foreignKey": {"menu_id"}, "references": {"id"}},
		}),
	)
	allTables = append(allTables, menuModel)

	// Log/LogType
	logTypeModel := g.GenerateModel("log_type")
	logModel := g.GenerateModel("log",
		gen.FieldRelate(field.HasOne, "LogType", logTypeModel, &field.RelateConfig{
			GORMTag: field.GormTag{"foreignKey": {"type_id"}, "references": {"id"}},
		}),
	)
	allTables = append(allTables, logTypeModel, logModel)

	// File/FileStorage/FileMenu
	fileModel := g.GenerateModel("file")
	fileMenuModel := g.GenerateModel("file_menu")
	fileStorageModel := g.GenerateModel("file_storage")
	allTables = append(allTables, fileModel, fileStorageModel, fileMenuModel)

	// Config
	configModel := g.GenerateModel("config")
	allTables = append(allTables, configModel)

	// Task
	taskModel := g.GenerateModel("task")
	allTables = append(allTables, taskModel)

	// 生成指定表
	g.ApplyBasic(allTables...)

	// 执行生成
	g.Execute()
}
