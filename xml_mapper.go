package xorm

import (
	"encoding/xml"
	"fmt"
	"strings"
)

type XmlIf struct {
	Test  string `xml:"test,attr"`
	Value string `xml:",cdata"`
}

type XmlSet struct {
	XMLName xml.Name `xml:""`
	Ifs     string   `xml:",innerxml"`
	//Ifs []XmlIf `xml:"if"`
}

type XmlSql struct {
	Id    string `xml:"id,attr"`
	Value string `xml:",cdata"`
	Set   XmlSet `xml:",any"`
}

type XmlInclude struct {
	RefId string `xml:"refid,attr"`
}

type XmlCrud struct {
	Id    string `xml:"id,attr"`
	Value string `xml:",innerxml"`
	//Includes []XmlInclude `xml:"include"`
}

type XmlInsert struct {
	XmlCrud
}

type XmlDelete struct {
	XmlCrud
}

type XmlSelect struct {
	XmlCrud
}

type XmlUpdate struct {
	XmlCrud
}

type XmlMapper struct {
	Namespace string      `xml:"namespace,attr"`
	Sqls      []XmlSql    `xml:"sql"`
	Inserts   []XmlInsert `xml:"insert"`
	Deletes   []XmlDelete `xml:"delete"`
	Selects   []XmlSelect `xml:"select"`
	Updates   []XmlUpdate `xml:"update"`
}

func (this *XmlIf) Parse() string {
	return fmt.Sprintf("%s%s{%% endif %%} ", this.Test, this.Value)
}

func (this *XmlSet) Parse() (res string, flag bool) {

	// sort by data & element in xml file
	vs, err := parseValue(this.Ifs)
	if err != nil {
		return
	}

	if this.XMLName.Local == "where" && len(vs) > 0 {
		res = stringJoin(res, "where 1=1 ")
	}

	for _, v := range vs {
		switch v._type {
		case ValueType.Char:
			res = stringJoin(res, v.Val, " ")
		case ValueType.Cdata:
			res = stringJoin(res, v.Val, " ")
		case ValueType.If:
			flag = true
			res = stringJoin(res, v.Val, " ")
		case ValueType.Include:
		case ValueType.Foreach:
		}
	}

	return
}

func (this *XmlSql) Parse() *Sql {
	sql := new(Sql)
	sql.Id = this.Id
	sql._type = SqlType.Define

	if len(this.Set.Ifs) > 0 {
		this.Set.Parse()
		sql.Value, sql._isAnalysis = this.Set.Parse()
	} else {
		sql.Value = stringJoin(strings.TrimSpace(this.Value), " ")
	}

	return sql
}

func stringJoin(args ...string) string {
	if len(args) == 0 {
		return ""
	}

	var buf []byte
	for _, v := range args {
		buf = append(buf, []byte(v)...)
	}

	return string(buf)
}

func (this *XmlCrud) Parse(fn func(string) *Sql) (*Sql, error) {

	sql := new(Sql)
	sql.Id = this.Id

	// sort by data & element in xml file
	vs, err := parseValue(this.Value)
	if err != nil {
		return nil, err
	}

	for _, v := range vs {
		switch v._type {
		case ValueType.Char:
			sql.Value = stringJoin(sql.Value, v.Val, " ")
		case ValueType.Cdata:
			sql.Value = stringJoin(sql.Value, v.Val, " ")
		case ValueType.Include:
			if v.Val == "" {
				break
			}
			defSql := fn(v.Val)
			if defSql == nil {
				break
			}
			if defSql._isAnalysis {
				sql._isAnalysis = true
			}
			sql.Value = stringJoin(sql.Value, defSql.Value, " ")
		case ValueType.Foreach:
		}
	}

	return sql, nil
}

func (this *XmlMapper) ParseMapper() (*Mapper, error) {

	mapper := new(Mapper)
	mapper.Namespace = this.Namespace
	capacity := len(this.Inserts) + len(this.Deletes) + len(this.Updates) + len(this.Selects)
	mapper.Sqls = make(map[string]*Sql, capacity)

	// prepare parse define sql (only one)
	defSqls := make(map[string]*Sql, len(this.Sqls))
	for _, sql := range this.Sqls {
		defSqls[sql.Id] = sql.Parse()
	}

	findDefSqlFn := func(refId string) *Sql {
		if sql, ok := defSqls[refId]; ok {
			return sql
		}

		return nil
	}

	//prepare parse insert sql
	for _, sql := range this.Inserts {
		s, err := sql.Parse(findDefSqlFn)
		if err != nil {
			return nil, err
		}
		s._type = SqlType.Insert
		mapper.Sqls[sql.Id] = s
	}

	//prepare parse delete sql
	for _, sql := range this.Deletes {
		s, err := sql.Parse(findDefSqlFn)
		if err != nil {
			return nil, err
		}
		s._type = SqlType.Delete
		mapper.Sqls[sql.Id] = s
	}

	//prepare parse update sql
	for _, sql := range this.Updates {
		s, err := sql.Parse(findDefSqlFn)
		if err != nil {
			return nil, err
		}
		s._type = SqlType.Update
		mapper.Sqls[sql.Id] = s
	}

	//prepare parse delete sql
	for _, sql := range this.Selects {
		s, err := sql.Parse(findDefSqlFn)
		if err != nil {
			return nil, err
		}
		s._type = SqlType.Select
		mapper.Sqls[sql.Id] = s
	}

	return mapper, nil
}

func UnmarshalXml(content []byte) (*XmlMapper, error) {
	mapper := new(XmlMapper)
	if err := xml.Unmarshal(content, mapper); err != nil {
		return nil, err
	}

	return mapper, nil
}

func UnmarshalMapper(content []byte) (*Mapper, error) {
	if xml, err := UnmarshalXml(content); err != nil {
		return nil, err
	} else {
		return xml.ParseMapper()
	}
}

var ValueType = struct {
	Char    int
	Cdata   int
	Include int
	If      int
	Foreach int
}{
	0,
	1,
	2,
	3,
	4,
}

type Value struct {
	Val   string
	_type int
}

/*
 *  min position form two tag
 *
 *  author: jerrylou, <gunsluo@gmail.com>
 *  since:  2017-03-08 14:05:44
 */
func minPos(s, e int) int {
	if s == -1 && e == -1 {
		return -1
	}

	if s >= 0 && e >= 0 {
		if s > e {
			return e
		} else {
			return s
		}
	}

	if s >= 0 {
		return s
	}

	return e
}

func parseValue(innerxml string) ([]*Value, error) {

	var idx int
	var start int
	var end int
	var typ int
	var vs []*Value

	var str string = innerxml
	for len(str) > 0 {
		start = strings.Index(str, "<![CDATA[")
		typ = ValueType.Cdata

		idx = strings.Index(str, "<include ")
		start = minPos(start, idx)
		if start == 0 && start == idx {
			typ = ValueType.Include
		}

		idx = strings.Index(str, "<foreach ")
		start = minPos(start, idx)
		if start == 0 && start == idx {
			typ = ValueType.Foreach
		}

		idx = strings.Index(str, "<if ")
		start = minPos(start, idx)
		if start == 0 && start == idx {
			typ = ValueType.If
		}

		if start == -1 {
			start = 0
			end = len(str)
			typ = ValueType.Char
		}

		if start > 0 {
			end = start
			start = 0
			typ = ValueType.Char
		}

		v := new(Value)
		v._type = typ

		switch typ {
		case ValueType.Char:
			val := []byte(str)
			v.Val = strings.TrimSpace(string(val[start : end-1]))
			if v.Val != "" {
				vs = append(vs, v)
			}
			str = string(val[end:])
		case ValueType.Cdata:
			end = strings.Index(str, "]]>")
			val := []byte(str)
			v.Val = strings.TrimSpace(string(val[start+9 : end]))
			if v.Val != "" {
				vs = append(vs, v)
			}
			str = string(val[end+3:])
		case ValueType.Include:
			end = strings.Index(str, "/>")
			val := []byte(str)
			var x XmlInclude
			if err := xml.Unmarshal(val[start:end+2], &x); err != nil {
				return nil, err
			}
			v.Val = x.RefId
			if v.Val != "" {
				vs = append(vs, v)
			}
			str = string(val[end+2:])
		case ValueType.If:
			end = strings.Index(str, "</if>")
			val := []byte(str)
			x := new(XmlIf)
			if err := xml.Unmarshal(val[start:end+5], x); err != nil {
				return nil, err
			}
			v.Val = x.Parse()
			if v.Val != "" {
				vs = append(vs, v)
			}
			str = string(val[end+5:])
		case ValueType.Foreach:
			end = strings.Index(str, "</foreach>")
			val := []byte(str)
			str = string(val[end+10:])
		}
	}

	return vs, nil
}
