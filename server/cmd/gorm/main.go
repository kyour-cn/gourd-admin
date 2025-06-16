package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"gourd/cmd/gorm/gen_tool"
	"gourd/cmd/gorm/methods"
	"gourd/cmd/gorm/tags"
	"gourd/internal/config"
)

// main 模型代码生成
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
		gen.FieldType("create_time", "uint"),
		gen.FieldType("update_time", "uint"),

		// 软删除字段属性
		gen.FieldType("delete_time", "soft_delete.DeletedAt"),

		// Json序列化
		gen.WithMethod(methods.JsonMethod{}),
	}

	g := gen.NewGenerator(gen.Config{
		OutPath: "./internal/orm/query",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	g.UseDB(mysqlDb)

	// 使用工具生成模型
	db := gen_tool.Database{
		Generator: g,
		ComOpts:   &comOpts,
		Tables: []gen_tool.Table{

			// 系统基础数据表
			{Name: "app"},
			{Name: "log"},
			{Name: "log_type"},
			{Name: "log_stat_view"},
			{Name: "menu_api"},
			{Name: "role"},
			{
				Name: "menu",
				Relate: &[]gen_tool.TableRelate{
					{
						TableName:  "menu_api",
						FieldName:  "MenuApi",
						Type:       field.HasMany,
						ForeignKey: "menu_id",
						LocalKey:   "id",
					},
					{
						TableName:  "app",
						FieldName:  "App",
						Type:       field.HasOne,
						ForeignKey: "app_id",
						LocalKey:   "id",
					},
				},
			},
			{
				Name: "user",
				Relate: &[]gen_tool.TableRelate{
					{
						TableName:  "role",
						FieldName:  "Role",
						Type:       field.HasOne,
						ForeignKey: "role_id",
						LocalKey:   "id",
						Relate: &[]gen_tool.TableRelate{
							{
								TableName:  "app",
								FieldName:  "App",
								Type:       field.HasOne,
								ForeignKey: "app_id",
								LocalKey:   "id",
							},
						},
					},
				},
			},

			// 业务数据表

		},
	}

	db.GenTable()
}
