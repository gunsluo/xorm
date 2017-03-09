package xorm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"text/template"

	"os"
	"path/filepath"
	"strings"

	"github.com/Unknwon/goconfig"
)

const (
	defaultNamespace = "default"
)

type SqlMap struct {
	SqlMapRootDir string
	Sql           map[string]string
	Mappers       map[string]*Mapper
	Extension     map[string]string
	Capacity      uint
	Cipher        Cipher
}

type SqlMapOptions struct {
	Capacity  uint
	Extension map[string]string
	Cipher    Cipher
}

var SqlType = struct {
	Define int
	Insert int
	Delete int
	Update int
	Select int
}{
	0,
	1,
	2,
	3,
	4,
}

type Sql struct {
	Id          string
	Value       string
	_type       int
	_isAnalysis bool
}

type Mapper struct {
	Namespace string
	Sqls      map[string]*Sql
}

func (engine *Engine) SetSqlMapCipher(cipher Cipher) {
	engine.sqlMap.Cipher = cipher
}

func (engine *Engine) ClearSqlMapCipher() {
	engine.sqlMap.Cipher = nil
}

func (sqlMap *SqlMap) checkNilAndInit() {
	if sqlMap.Sql == nil {
		if sqlMap.Capacity == 0 {
			sqlMap.Sql = make(map[string]string, 100)
			sqlMap.Mappers = make(map[string]*Mapper, 100)
		} else {
			sqlMap.Sql = make(map[string]string, sqlMap.Capacity)
			sqlMap.Mappers = make(map[string]*Mapper, sqlMap.Capacity)
		}

	}
}

func (engine *Engine) InitSqlMap(options ...SqlMapOptions) error {
	var opt SqlMapOptions

	if len(options) > 0 {
		opt = options[0]
	}

	if len(opt.Extension) == 0 {
		opt.Extension = map[string]string{"xml": ".xml", "json": ".json"}
	} else {
		if opt.Extension["xml"] == "" || len(opt.Extension["xml"]) == 0 {
			opt.Extension["xml"] = ".xml"
		}
		if opt.Extension["json"] == "" || len(opt.Extension["json"]) == 0 {
			opt.Extension["json"] = ".json"
		}
	}

	engine.sqlMap.Extension = opt.Extension
	engine.sqlMap.Capacity = opt.Capacity

	engine.sqlMap.Cipher = opt.Cipher

	var err error
	if engine.sqlMap.SqlMapRootDir == "" {
		cfg, err := goconfig.LoadConfigFile("./sql/xormcfg.ini")
		if err != nil {
			return err
		}
		engine.sqlMap.SqlMapRootDir, err = cfg.GetValue("", "SqlMapRootDir")
		if err != nil {
			return err
		}
	}

	err = filepath.Walk(engine.sqlMap.SqlMapRootDir, engine.sqlMap.walkFunc)
	if err != nil {
		return err
	}

	return nil
}

func (engine *Engine) LoadSqlMap(filepath string) error {

	if len(engine.sqlMap.Extension) == 0 {
		engine.sqlMap.Extension = map[string]string{"xml": ".xml", "json": ".json"}
	} else {
		if engine.sqlMap.Extension["xml"] == "" || len(engine.sqlMap.Extension["xml"]) == 0 {
			engine.sqlMap.Extension["xml"] = ".xml"
		}
		if engine.sqlMap.Extension["json"] == "" || len(engine.sqlMap.Extension["json"]) == 0 {
			engine.sqlMap.Extension["json"] = ".json"
		}
	}

	if strings.HasSuffix(filepath, engine.sqlMap.Extension["xml"]) || strings.HasSuffix(filepath, engine.sqlMap.Extension["json"]) {
		err := engine.loadSqlMap(filepath)
		if err != nil {
			return err
		}
	}

	return nil
}

func (engine *Engine) BatchLoadSqlMap(filepathSlice []string) error {
	if len(engine.sqlMap.Extension) == 0 {
		engine.sqlMap.Extension = map[string]string{"xml": ".xml", "json": ".json"}
	} else {
		if engine.sqlMap.Extension["xml"] == "" || len(engine.sqlMap.Extension["xml"]) == 0 {
			engine.sqlMap.Extension["xml"] = ".xml"
		}
		if engine.sqlMap.Extension["json"] == "" || len(engine.sqlMap.Extension["json"]) == 0 {
			engine.sqlMap.Extension["json"] = ".json"
		}
	}

	for _, filepath := range filepathSlice {
		if strings.HasSuffix(filepath, engine.sqlMap.Extension["xml"]) || strings.HasSuffix(filepath, engine.sqlMap.Extension["json"]) {
			err := engine.loadSqlMap(filepath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (engine *Engine) ReloadSqlMap(filepath string) error {
	if len(engine.sqlMap.Extension) == 0 {
		engine.sqlMap.Extension = map[string]string{"xml": ".xml", "json": ".json"}
	} else {
		if engine.sqlMap.Extension["xml"] == "" || len(engine.sqlMap.Extension["xml"]) == 0 {
			engine.sqlMap.Extension["xml"] = ".xml"
		}
		if engine.sqlMap.Extension["json"] == "" || len(engine.sqlMap.Extension["json"]) == 0 {
			engine.sqlMap.Extension["json"] = ".json"
		}
	}

	if strings.HasSuffix(filepath, engine.sqlMap.Extension["xml"]) || strings.HasSuffix(filepath, engine.sqlMap.Extension["json"]) {
		err := engine.reloadSqlMap(filepath)
		if err != nil {
			return err
		}
	}

	return nil
}

func (engine *Engine) BatchReloadSqlMap(filepathSlice []string) error {
	if len(engine.sqlMap.Extension) == 0 {
		engine.sqlMap.Extension = map[string]string{"xml": ".xml", "json": ".json"}
	} else {
		if engine.sqlMap.Extension["xml"] == "" || len(engine.sqlMap.Extension["xml"]) == 0 {
			engine.sqlMap.Extension["xml"] = ".xml"
		}
		if engine.sqlMap.Extension["json"] == "" || len(engine.sqlMap.Extension["json"]) == 0 {
			engine.sqlMap.Extension["json"] = ".json"
		}
	}

	for _, filepath := range filepathSlice {
		if strings.HasSuffix(filepath, engine.sqlMap.Extension["xml"]) || strings.HasSuffix(filepath, engine.sqlMap.Extension["json"]) {
			err := engine.loadSqlMap(filepath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (engine *Engine) loadSqlMap(filepath string) error {

	info, err := os.Lstat(filepath)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}

	err = engine.sqlMap.paresSql(filepath)
	if err != nil {
		return err
	}

	return nil
}

func (engine *Engine) reloadSqlMap(filepath string) error {

	info, err := os.Lstat(filepath)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}
	err = engine.sqlMap.paresSql(filepath)
	if err != nil {
		return err
	}

	return nil
}

func (sqlMap *SqlMap) walkFunc(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}

	if strings.HasSuffix(path, sqlMap.Extension["xml"]) || strings.HasSuffix(path, sqlMap.Extension["json"]) {
		err = sqlMap.paresSql(path)
		if err != nil {
			return err
		}
	}

	return nil
}

func (sqlMap *SqlMap) paresSql(filepath string) error {

	content, err := ioutil.ReadFile(filepath)

	if err != nil {
		return err
	}
	enc := sqlMap.Cipher
	if enc != nil {
		content, err = enc.Decrypt(content)

		if err != nil {
			return err
		}
	}

	sqlMap.checkNilAndInit()

	if strings.HasSuffix(filepath, sqlMap.Extension["xml"]) {

		mapper, err := UnmarshalMapper(content)
		if err != nil {
			return err
		}

		if mapper.Namespace == "" {
			mapper.Namespace = defaultNamespace
		}

		if _, ok := sqlMap.Mappers[mapper.Namespace]; ok {
			// warning: same namespace name
			// sqlMap.Mappers[mapper.Namespace] = mapper
		} else {
			sqlMap.Mappers[mapper.Namespace] = mapper
		}

		return nil
	}

	if strings.HasSuffix(filepath, sqlMap.Extension["json"]) {
		var result map[string]string
		err = json.Unmarshal(content, &result)
		if err != nil {
			return err
		}
		for k := range result {
			sqlMap.Sql[k] = result[k]
		}

		return nil
	}
	return nil

}

func (engine *Engine) AddSql(key string, sql string) {
	engine.sqlMap.addSql(key, sql)
}

func (sqlMap *SqlMap) addSql(key string, sql string) {
	sqlMap.checkNilAndInit()
	sqlMap.Sql[key] = sql
}

func (engine *Engine) UpdateSql(key string, sql string) {
	engine.sqlMap.updateSql(key, sql)
}

func (sqlMap *SqlMap) updateSql(key string, sql string) {
	sqlMap.checkNilAndInit()
	sqlMap.Sql[key] = sql
}

func (engine *Engine) RemoveSql(key string) {
	engine.sqlMap.removeSql(key)
}

func (sqlMap *SqlMap) removeSql(key string) {
	sqlMap.checkNilAndInit()
	delete(sqlMap.Sql, key)
}

func (engine *Engine) BatchAddSql(sqlStrMap map[string]string) {
	engine.sqlMap.batchAddSql(sqlStrMap)
}

func (sqlMap *SqlMap) batchAddSql(sqlStrMap map[string]string) {
	sqlMap.checkNilAndInit()
	for k, v := range sqlStrMap {
		sqlMap.Sql[k] = v
	}
}

func (engine *Engine) BatchUpdateSql(sqlStrMap map[string]string) {
	engine.sqlMap.batchUpdateSql(sqlStrMap)
}

func (sqlMap *SqlMap) batchUpdateSql(sqlStrMap map[string]string) {
	sqlMap.checkNilAndInit()
	for k, v := range sqlStrMap {
		sqlMap.Sql[k] = v
	}
}

func (engine *Engine) BatchRemoveSql(key []string) {
	engine.sqlMap.batchRemoveSql(key)
}

func (sqlMap *SqlMap) batchRemoveSql(key []string) {
	sqlMap.checkNilAndInit()
	for _, v := range key {
		delete(sqlMap.Sql, v)
	}
}

func (engine *Engine) GetSql(key string) string {
	return engine.sqlMap.getSql(key)
}

func (sqlMap *SqlMap) getSql(key string) string {
	return sqlMap.Sql[key]
}

func (engine *Engine) GetSqlMap(keys ...interface{}) map[string]string {
	return engine.sqlMap.getSqlMap(keys...)
}

func (sqlMap *SqlMap) getSqlMap(keys ...interface{}) map[string]string {
	var resultSqlMap map[string]string
	i := len(keys)
	if i == 0 {
		return sqlMap.Sql
	}

	if i == 1 {
		switch keys[0].(type) {
		case string:
			resultSqlMap = make(map[string]string, 1)
		case []string:
			ks := keys[0].([]string)
			n := len(ks)
			resultSqlMap = make(map[string]string, n)
		}
	} else {
		resultSqlMap = make(map[string]string, i)
	}

	for k, _ := range keys {
		switch keys[k].(type) {
		case string:
			key := keys[k].(string)
			resultSqlMap[key] = sqlMap.Sql[key]
		case []string:
			ks := keys[k].([]string)
			for _, v := range ks {
				resultSqlMap[v] = sqlMap.Sql[v]
			}
		}
	}

	return resultSqlMap
}

func parseTagName(sqlTagName string) (string, string) {
	idx := strings.Index(sqlTagName, ".")
	if idx == -1 {
		return defaultNamespace, sqlTagName
	}

	buf := []byte(sqlTagName)
	namespace := string(buf[:idx])
	id := string(buf[idx+1:])
	if namespace == "" {
		return defaultNamespace, id
	}

	return namespace, id
}

func (sqlMap *SqlMap) getMapperSql(sqlTagName string) *Sql {
	namespace, id := parseTagName(sqlTagName)
	mapper, ok := sqlMap.Mappers[namespace]
	if !ok {
		return &Sql{Value: fmt.Sprintf("namespace[%s] not found", namespace)}
	}

	sql, ok := mapper.Sqls[id]
	if !ok {
		return &Sql{Value: fmt.Sprintf("[%s.%s] not found", namespace, id)}
	}

	return sql
}

var (
	defaultTpl = template.New("mapper")
)

func formatTemplate(pattern string, args ...interface{}) string {
	if len(args) == 0 {
		return ""
	}

	var buf bytes.Buffer
	tpl := template.Must(defaultTpl.Parse(pattern))
	err := tpl.Execute(&buf, args[0])
	if err != nil {
		return ""
	}

	return buf.String()
}

func (this *Sql) Format(args ...interface{}) string {
	if this._isAnalysis {
		return formatTemplate(this.Value, args...)
	}

	if len(args) > 0 {
		return formatTemplate(this.Value, args...)
	}

	return this.Value
}
