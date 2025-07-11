// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"app/internal/orm/model"
)

func newRole(db *gorm.DB, opts ...gen.DOOption) role {
	_role := role{}

	_role.roleDo.UseDB(db, opts...)
	_role.roleDo.UseModel(&model.Role{})

	tableName := _role.roleDo.TableName()
	_role.ALL = field.NewAsterisk(tableName)
	_role.ID = field.NewInt32(tableName, "id")
	_role.AppID = field.NewInt32(tableName, "app_id")
	_role.Name = field.NewString(tableName, "name")
	_role.Rules = field.NewString(tableName, "rules")
	_role.RulesCheckd = field.NewString(tableName, "rules_checkd")
	_role.CreateTime = field.NewInt32(tableName, "create_time")
	_role.UpdateTime = field.NewInt32(tableName, "update_time")
	_role.Remark = field.NewString(tableName, "remark")
	_role.Status = field.NewInt32(tableName, "status")
	_role.Sort = field.NewInt32(tableName, "sort")
	_role.IsAdmin = field.NewInt32(tableName, "is_admin")
	_role.App = roleHasOneApp{
		db: db.Session(&gorm.Session{}),

		RelationField: field.NewRelation("App", "model.App"),
	}

	_role.fillFieldMap()

	return _role
}

// role 用户角色
type role struct {
	roleDo

	ALL         field.Asterisk
	ID          field.Int32
	AppID       field.Int32  // 应用ID
	Name        field.String // 角色名称
	Rules       field.String // 权限ID ,分割a
	RulesCheckd field.String // 权限树选中的字节点ID
	CreateTime  field.Int32  // 创建时间
	UpdateTime  field.Int32  // 更新时间
	Remark      field.String // 简介
	Status      field.Int32  // 状态
	Sort        field.Int32  // 排序
	IsAdmin     field.Int32  // 是否为管理员（所有权限）
	App         roleHasOneApp

	fieldMap map[string]field.Expr
}

func (r role) Table(newTableName string) *role {
	r.roleDo.UseTable(newTableName)
	return r.updateTableName(newTableName)
}

func (r role) As(alias string) *role {
	r.roleDo.DO = *(r.roleDo.As(alias).(*gen.DO))
	return r.updateTableName(alias)
}

func (r *role) updateTableName(table string) *role {
	r.ALL = field.NewAsterisk(table)
	r.ID = field.NewInt32(table, "id")
	r.AppID = field.NewInt32(table, "app_id")
	r.Name = field.NewString(table, "name")
	r.Rules = field.NewString(table, "rules")
	r.RulesCheckd = field.NewString(table, "rules_checkd")
	r.CreateTime = field.NewInt32(table, "create_time")
	r.UpdateTime = field.NewInt32(table, "update_time")
	r.Remark = field.NewString(table, "remark")
	r.Status = field.NewInt32(table, "status")
	r.Sort = field.NewInt32(table, "sort")
	r.IsAdmin = field.NewInt32(table, "is_admin")

	r.fillFieldMap()

	return r
}

func (r *role) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := r.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (r *role) fillFieldMap() {
	r.fieldMap = make(map[string]field.Expr, 12)
	r.fieldMap["id"] = r.ID
	r.fieldMap["app_id"] = r.AppID
	r.fieldMap["name"] = r.Name
	r.fieldMap["rules"] = r.Rules
	r.fieldMap["rules_checkd"] = r.RulesCheckd
	r.fieldMap["create_time"] = r.CreateTime
	r.fieldMap["update_time"] = r.UpdateTime
	r.fieldMap["remark"] = r.Remark
	r.fieldMap["status"] = r.Status
	r.fieldMap["sort"] = r.Sort
	r.fieldMap["is_admin"] = r.IsAdmin

}

func (r role) clone(db *gorm.DB) role {
	r.roleDo.ReplaceConnPool(db.Statement.ConnPool)
	r.App.db = db.Session(&gorm.Session{Initialized: true})
	r.App.db.Statement.ConnPool = db.Statement.ConnPool
	return r
}

func (r role) replaceDB(db *gorm.DB) role {
	r.roleDo.ReplaceDB(db)
	r.App.db = db.Session(&gorm.Session{})
	return r
}

type roleHasOneApp struct {
	db *gorm.DB

	field.RelationField
}

func (a roleHasOneApp) Where(conds ...field.Expr) *roleHasOneApp {
	if len(conds) == 0 {
		return &a
	}

	exprs := make([]clause.Expression, 0, len(conds))
	for _, cond := range conds {
		exprs = append(exprs, cond.BeCond().(clause.Expression))
	}
	a.db = a.db.Clauses(clause.Where{Exprs: exprs})
	return &a
}

func (a roleHasOneApp) WithContext(ctx context.Context) *roleHasOneApp {
	a.db = a.db.WithContext(ctx)
	return &a
}

func (a roleHasOneApp) Session(session *gorm.Session) *roleHasOneApp {
	a.db = a.db.Session(session)
	return &a
}

func (a roleHasOneApp) Model(m *model.Role) *roleHasOneAppTx {
	return &roleHasOneAppTx{a.db.Model(m).Association(a.Name())}
}

func (a roleHasOneApp) Unscoped() *roleHasOneApp {
	a.db = a.db.Unscoped()
	return &a
}

type roleHasOneAppTx struct{ tx *gorm.Association }

func (a roleHasOneAppTx) Find() (result *model.App, err error) {
	return result, a.tx.Find(&result)
}

func (a roleHasOneAppTx) Append(values ...*model.App) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Append(targetValues...)
}

func (a roleHasOneAppTx) Replace(values ...*model.App) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a roleHasOneAppTx) Delete(values ...*model.App) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Delete(targetValues...)
}

func (a roleHasOneAppTx) Clear() error {
	return a.tx.Clear()
}

func (a roleHasOneAppTx) Count() int64 {
	return a.tx.Count()
}

func (a roleHasOneAppTx) Unscoped() *roleHasOneAppTx {
	a.tx = a.tx.Unscoped()
	return &a
}

type roleDo struct{ gen.DO }

type IRoleDo interface {
	gen.SubQuery
	Debug() IRoleDo
	WithContext(ctx context.Context) IRoleDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IRoleDo
	WriteDB() IRoleDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IRoleDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IRoleDo
	Not(conds ...gen.Condition) IRoleDo
	Or(conds ...gen.Condition) IRoleDo
	Select(conds ...field.Expr) IRoleDo
	Where(conds ...gen.Condition) IRoleDo
	Order(conds ...field.Expr) IRoleDo
	Distinct(cols ...field.Expr) IRoleDo
	Omit(cols ...field.Expr) IRoleDo
	Join(table schema.Tabler, on ...field.Expr) IRoleDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IRoleDo
	RightJoin(table schema.Tabler, on ...field.Expr) IRoleDo
	Group(cols ...field.Expr) IRoleDo
	Having(conds ...gen.Condition) IRoleDo
	Limit(limit int) IRoleDo
	Offset(offset int) IRoleDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IRoleDo
	Unscoped() IRoleDo
	Create(values ...*model.Role) error
	CreateInBatches(values []*model.Role, batchSize int) error
	Save(values ...*model.Role) error
	First() (*model.Role, error)
	Take() (*model.Role, error)
	Last() (*model.Role, error)
	Find() ([]*model.Role, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Role, err error)
	FindInBatches(result *[]*model.Role, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.Role) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IRoleDo
	Assign(attrs ...field.AssignExpr) IRoleDo
	Joins(fields ...field.RelationField) IRoleDo
	Preload(fields ...field.RelationField) IRoleDo
	FirstOrInit() (*model.Role, error)
	FirstOrCreate() (*model.Role, error)
	FindByPage(offset int, limit int) (result []*model.Role, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Rows() (*sql.Rows, error)
	Row() *sql.Row
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IRoleDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (r roleDo) Debug() IRoleDo {
	return r.withDO(r.DO.Debug())
}

func (r roleDo) WithContext(ctx context.Context) IRoleDo {
	return r.withDO(r.DO.WithContext(ctx))
}

func (r roleDo) ReadDB() IRoleDo {
	return r.Clauses(dbresolver.Read)
}

func (r roleDo) WriteDB() IRoleDo {
	return r.Clauses(dbresolver.Write)
}

func (r roleDo) Session(config *gorm.Session) IRoleDo {
	return r.withDO(r.DO.Session(config))
}

func (r roleDo) Clauses(conds ...clause.Expression) IRoleDo {
	return r.withDO(r.DO.Clauses(conds...))
}

func (r roleDo) Returning(value interface{}, columns ...string) IRoleDo {
	return r.withDO(r.DO.Returning(value, columns...))
}

func (r roleDo) Not(conds ...gen.Condition) IRoleDo {
	return r.withDO(r.DO.Not(conds...))
}

func (r roleDo) Or(conds ...gen.Condition) IRoleDo {
	return r.withDO(r.DO.Or(conds...))
}

func (r roleDo) Select(conds ...field.Expr) IRoleDo {
	return r.withDO(r.DO.Select(conds...))
}

func (r roleDo) Where(conds ...gen.Condition) IRoleDo {
	return r.withDO(r.DO.Where(conds...))
}

func (r roleDo) Order(conds ...field.Expr) IRoleDo {
	return r.withDO(r.DO.Order(conds...))
}

func (r roleDo) Distinct(cols ...field.Expr) IRoleDo {
	return r.withDO(r.DO.Distinct(cols...))
}

func (r roleDo) Omit(cols ...field.Expr) IRoleDo {
	return r.withDO(r.DO.Omit(cols...))
}

func (r roleDo) Join(table schema.Tabler, on ...field.Expr) IRoleDo {
	return r.withDO(r.DO.Join(table, on...))
}

func (r roleDo) LeftJoin(table schema.Tabler, on ...field.Expr) IRoleDo {
	return r.withDO(r.DO.LeftJoin(table, on...))
}

func (r roleDo) RightJoin(table schema.Tabler, on ...field.Expr) IRoleDo {
	return r.withDO(r.DO.RightJoin(table, on...))
}

func (r roleDo) Group(cols ...field.Expr) IRoleDo {
	return r.withDO(r.DO.Group(cols...))
}

func (r roleDo) Having(conds ...gen.Condition) IRoleDo {
	return r.withDO(r.DO.Having(conds...))
}

func (r roleDo) Limit(limit int) IRoleDo {
	return r.withDO(r.DO.Limit(limit))
}

func (r roleDo) Offset(offset int) IRoleDo {
	return r.withDO(r.DO.Offset(offset))
}

func (r roleDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IRoleDo {
	return r.withDO(r.DO.Scopes(funcs...))
}

func (r roleDo) Unscoped() IRoleDo {
	return r.withDO(r.DO.Unscoped())
}

func (r roleDo) Create(values ...*model.Role) error {
	if len(values) == 0 {
		return nil
	}
	return r.DO.Create(values)
}

func (r roleDo) CreateInBatches(values []*model.Role, batchSize int) error {
	return r.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (r roleDo) Save(values ...*model.Role) error {
	if len(values) == 0 {
		return nil
	}
	return r.DO.Save(values)
}

func (r roleDo) First() (*model.Role, error) {
	if result, err := r.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Role), nil
	}
}

func (r roleDo) Take() (*model.Role, error) {
	if result, err := r.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Role), nil
	}
}

func (r roleDo) Last() (*model.Role, error) {
	if result, err := r.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Role), nil
	}
}

func (r roleDo) Find() ([]*model.Role, error) {
	result, err := r.DO.Find()
	return result.([]*model.Role), err
}

func (r roleDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Role, err error) {
	buf := make([]*model.Role, 0, batchSize)
	err = r.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (r roleDo) FindInBatches(result *[]*model.Role, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return r.DO.FindInBatches(result, batchSize, fc)
}

func (r roleDo) Attrs(attrs ...field.AssignExpr) IRoleDo {
	return r.withDO(r.DO.Attrs(attrs...))
}

func (r roleDo) Assign(attrs ...field.AssignExpr) IRoleDo {
	return r.withDO(r.DO.Assign(attrs...))
}

func (r roleDo) Joins(fields ...field.RelationField) IRoleDo {
	for _, _f := range fields {
		r = *r.withDO(r.DO.Joins(_f))
	}
	return &r
}

func (r roleDo) Preload(fields ...field.RelationField) IRoleDo {
	for _, _f := range fields {
		r = *r.withDO(r.DO.Preload(_f))
	}
	return &r
}

func (r roleDo) FirstOrInit() (*model.Role, error) {
	if result, err := r.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Role), nil
	}
}

func (r roleDo) FirstOrCreate() (*model.Role, error) {
	if result, err := r.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Role), nil
	}
}

func (r roleDo) FindByPage(offset int, limit int) (result []*model.Role, count int64, err error) {
	result, err = r.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = r.Offset(-1).Limit(-1).Count()
	return
}

func (r roleDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = r.Count()
	if err != nil {
		return
	}

	err = r.Offset(offset).Limit(limit).Scan(result)
	return
}

func (r roleDo) Scan(result interface{}) (err error) {
	return r.DO.Scan(result)
}

func (r roleDo) Delete(models ...*model.Role) (result gen.ResultInfo, err error) {
	return r.DO.Delete(models)
}

func (r *roleDo) withDO(do gen.Dao) *roleDo {
	r.DO = *do.(*gen.DO)
	return r
}
