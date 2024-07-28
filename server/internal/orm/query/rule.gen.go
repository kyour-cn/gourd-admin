// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"gourd/internal/orm/model"
)

func newRule(db *gorm.DB, opts ...gen.DOOption) rule {
	_rule := rule{}

	_rule.ruleDo.UseDB(db, opts...)
	_rule.ruleDo.UseModel(&model.Rule{})

	tableName := _rule.ruleDo.TableName()
	_rule.ALL = field.NewAsterisk(tableName)
	_rule.ID = field.NewInt32(tableName, "id")
	_rule.AppID = field.NewInt32(tableName, "app_id")
	_rule.Name = field.NewString(tableName, "name")
	_rule.Alias_ = field.NewString(tableName, "alias")
	_rule.Path = field.NewString(tableName, "path")
	_rule.Pid = field.NewInt32(tableName, "pid")
	_rule.Status = field.NewInt32(tableName, "status")
	_rule.Sort = field.NewInt32(tableName, "sort")
	_rule.AddonID = field.NewInt32(tableName, "addon_id")

	_rule.fillFieldMap()

	return _rule
}

// rule 权限规则表
type rule struct {
	ruleDo

	ALL     field.Asterisk
	ID      field.Int32
	AppID   field.Int32  // 应用ID
	Name    field.String // 名字
	Alias_  field.String // 英文别名
	Path    field.String // 规则
	Pid     field.Int32  // 上级Id
	Status  field.Int32  // 状态 0:1
	Sort    field.Int32  // 排序
	AddonID field.Int32  // 插件ID 为0=不验证插件

	fieldMap map[string]field.Expr
}

func (r rule) Table(newTableName string) *rule {
	r.ruleDo.UseTable(newTableName)
	return r.updateTableName(newTableName)
}

func (r rule) As(alias string) *rule {
	r.ruleDo.DO = *(r.ruleDo.As(alias).(*gen.DO))
	return r.updateTableName(alias)
}

func (r *rule) updateTableName(table string) *rule {
	r.ALL = field.NewAsterisk(table)
	r.ID = field.NewInt32(table, "id")
	r.AppID = field.NewInt32(table, "app_id")
	r.Name = field.NewString(table, "name")
	r.Alias_ = field.NewString(table, "alias")
	r.Path = field.NewString(table, "path")
	r.Pid = field.NewInt32(table, "pid")
	r.Status = field.NewInt32(table, "status")
	r.Sort = field.NewInt32(table, "sort")
	r.AddonID = field.NewInt32(table, "addon_id")

	r.fillFieldMap()

	return r
}

func (r *rule) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := r.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (r *rule) fillFieldMap() {
	r.fieldMap = make(map[string]field.Expr, 9)
	r.fieldMap["id"] = r.ID
	r.fieldMap["app_id"] = r.AppID
	r.fieldMap["name"] = r.Name
	r.fieldMap["alias"] = r.Alias_
	r.fieldMap["path"] = r.Path
	r.fieldMap["pid"] = r.Pid
	r.fieldMap["status"] = r.Status
	r.fieldMap["sort"] = r.Sort
	r.fieldMap["addon_id"] = r.AddonID
}

func (r rule) clone(db *gorm.DB) rule {
	r.ruleDo.ReplaceConnPool(db.Statement.ConnPool)
	return r
}

func (r rule) replaceDB(db *gorm.DB) rule {
	r.ruleDo.ReplaceDB(db)
	return r
}

type ruleDo struct{ gen.DO }

type IRuleDo interface {
	gen.SubQuery
	Debug() IRuleDo
	WithContext(ctx context.Context) IRuleDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IRuleDo
	WriteDB() IRuleDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IRuleDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IRuleDo
	Not(conds ...gen.Condition) IRuleDo
	Or(conds ...gen.Condition) IRuleDo
	Select(conds ...field.Expr) IRuleDo
	Where(conds ...gen.Condition) IRuleDo
	Order(conds ...field.Expr) IRuleDo
	Distinct(cols ...field.Expr) IRuleDo
	Omit(cols ...field.Expr) IRuleDo
	Join(table schema.Tabler, on ...field.Expr) IRuleDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IRuleDo
	RightJoin(table schema.Tabler, on ...field.Expr) IRuleDo
	Group(cols ...field.Expr) IRuleDo
	Having(conds ...gen.Condition) IRuleDo
	Limit(limit int) IRuleDo
	Offset(offset int) IRuleDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IRuleDo
	Unscoped() IRuleDo
	Create(values ...*model.Rule) error
	CreateInBatches(values []*model.Rule, batchSize int) error
	Save(values ...*model.Rule) error
	First() (*model.Rule, error)
	Take() (*model.Rule, error)
	Last() (*model.Rule, error)
	Find() ([]*model.Rule, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Rule, err error)
	FindInBatches(result *[]*model.Rule, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.Rule) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IRuleDo
	Assign(attrs ...field.AssignExpr) IRuleDo
	Joins(fields ...field.RelationField) IRuleDo
	Preload(fields ...field.RelationField) IRuleDo
	FirstOrInit() (*model.Rule, error)
	FirstOrCreate() (*model.Rule, error)
	FindByPage(offset int, limit int) (result []*model.Rule, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IRuleDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (r ruleDo) Debug() IRuleDo {
	return r.withDO(r.DO.Debug())
}

func (r ruleDo) WithContext(ctx context.Context) IRuleDo {
	return r.withDO(r.DO.WithContext(ctx))
}

func (r ruleDo) ReadDB() IRuleDo {
	return r.Clauses(dbresolver.Read)
}

func (r ruleDo) WriteDB() IRuleDo {
	return r.Clauses(dbresolver.Write)
}

func (r ruleDo) Session(config *gorm.Session) IRuleDo {
	return r.withDO(r.DO.Session(config))
}

func (r ruleDo) Clauses(conds ...clause.Expression) IRuleDo {
	return r.withDO(r.DO.Clauses(conds...))
}

func (r ruleDo) Returning(value interface{}, columns ...string) IRuleDo {
	return r.withDO(r.DO.Returning(value, columns...))
}

func (r ruleDo) Not(conds ...gen.Condition) IRuleDo {
	return r.withDO(r.DO.Not(conds...))
}

func (r ruleDo) Or(conds ...gen.Condition) IRuleDo {
	return r.withDO(r.DO.Or(conds...))
}

func (r ruleDo) Select(conds ...field.Expr) IRuleDo {
	return r.withDO(r.DO.Select(conds...))
}

func (r ruleDo) Where(conds ...gen.Condition) IRuleDo {
	return r.withDO(r.DO.Where(conds...))
}

func (r ruleDo) Order(conds ...field.Expr) IRuleDo {
	return r.withDO(r.DO.Order(conds...))
}

func (r ruleDo) Distinct(cols ...field.Expr) IRuleDo {
	return r.withDO(r.DO.Distinct(cols...))
}

func (r ruleDo) Omit(cols ...field.Expr) IRuleDo {
	return r.withDO(r.DO.Omit(cols...))
}

func (r ruleDo) Join(table schema.Tabler, on ...field.Expr) IRuleDo {
	return r.withDO(r.DO.Join(table, on...))
}

func (r ruleDo) LeftJoin(table schema.Tabler, on ...field.Expr) IRuleDo {
	return r.withDO(r.DO.LeftJoin(table, on...))
}

func (r ruleDo) RightJoin(table schema.Tabler, on ...field.Expr) IRuleDo {
	return r.withDO(r.DO.RightJoin(table, on...))
}

func (r ruleDo) Group(cols ...field.Expr) IRuleDo {
	return r.withDO(r.DO.Group(cols...))
}

func (r ruleDo) Having(conds ...gen.Condition) IRuleDo {
	return r.withDO(r.DO.Having(conds...))
}

func (r ruleDo) Limit(limit int) IRuleDo {
	return r.withDO(r.DO.Limit(limit))
}

func (r ruleDo) Offset(offset int) IRuleDo {
	return r.withDO(r.DO.Offset(offset))
}

func (r ruleDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IRuleDo {
	return r.withDO(r.DO.Scopes(funcs...))
}

func (r ruleDo) Unscoped() IRuleDo {
	return r.withDO(r.DO.Unscoped())
}

func (r ruleDo) Create(values ...*model.Rule) error {
	if len(values) == 0 {
		return nil
	}
	return r.DO.Create(values)
}

func (r ruleDo) CreateInBatches(values []*model.Rule, batchSize int) error {
	return r.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (r ruleDo) Save(values ...*model.Rule) error {
	if len(values) == 0 {
		return nil
	}
	return r.DO.Save(values)
}

func (r ruleDo) First() (*model.Rule, error) {
	if result, err := r.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Rule), nil
	}
}

func (r ruleDo) Take() (*model.Rule, error) {
	if result, err := r.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Rule), nil
	}
}

func (r ruleDo) Last() (*model.Rule, error) {
	if result, err := r.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Rule), nil
	}
}

func (r ruleDo) Find() ([]*model.Rule, error) {
	result, err := r.DO.Find()
	return result.([]*model.Rule), err
}

func (r ruleDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Rule, err error) {
	buf := make([]*model.Rule, 0, batchSize)
	err = r.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (r ruleDo) FindInBatches(result *[]*model.Rule, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return r.DO.FindInBatches(result, batchSize, fc)
}

func (r ruleDo) Attrs(attrs ...field.AssignExpr) IRuleDo {
	return r.withDO(r.DO.Attrs(attrs...))
}

func (r ruleDo) Assign(attrs ...field.AssignExpr) IRuleDo {
	return r.withDO(r.DO.Assign(attrs...))
}

func (r ruleDo) Joins(fields ...field.RelationField) IRuleDo {
	for _, _f := range fields {
		r = *r.withDO(r.DO.Joins(_f))
	}
	return &r
}

func (r ruleDo) Preload(fields ...field.RelationField) IRuleDo {
	for _, _f := range fields {
		r = *r.withDO(r.DO.Preload(_f))
	}
	return &r
}

func (r ruleDo) FirstOrInit() (*model.Rule, error) {
	if result, err := r.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Rule), nil
	}
}

func (r ruleDo) FirstOrCreate() (*model.Rule, error) {
	if result, err := r.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Rule), nil
	}
}

func (r ruleDo) FindByPage(offset int, limit int) (result []*model.Rule, count int64, err error) {
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

func (r ruleDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = r.Count()
	if err != nil {
		return
	}

	err = r.Offset(offset).Limit(limit).Scan(result)
	return
}

func (r ruleDo) Scan(result interface{}) (err error) {
	return r.DO.Scan(result)
}

func (r ruleDo) Delete(models ...*model.Rule) (result gen.ResultInfo, err error) {
	return r.DO.Delete(models)
}

func (r *ruleDo) withDO(do gen.Dao) *ruleDo {
	r.DO = *do.(*gen.DO)
	return r
}
