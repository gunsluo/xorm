package xorm

import (
	"encoding/xml"
	"fmt"
)

type XmlIf struct {
	Test  string `xml:"test,attr"`
	Value string `xml:",chardata"`
}

type XmlSet struct {
	Ifs []XmlIf `xml:"if"`
}

type XmlSql struct {
	Id    string `xml:"id,attr"`
	Value string `xml:",chardata"`
	Set   XmlSet `xml:"set"`
}

type XmlInclude struct {
	RefId string `xml:"refid,attr"`
}

type XmlCrud struct {
	Id       string       `xml:"id,attr"`
	Value    string       `xml:",chardata"`
	Includes []XmlInclude `xml:"include"`
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
	return fmt.Sprintf("%s%s{{- end}}", this.Test, this.Value)
}

func (this *XmlSql) Parse() *Sql {
	sql := new(Sql)
	sql.Id = this.Id
	sql._type = SqlType.Define

	if len(this.Set.Ifs) > 0 {
		sql._isAnalysis = true
		for _, xmlif := range this.Set.Ifs {
			sql.Value += xmlif.Parse()
		}
		sql.Value += this.Value
	} else {
		sql.Value = this.Value
	}

	return sql
}

func (this *XmlCrud) Parse(fn func(string) *Sql) *Sql {

	sql := new(Sql)
	sql.Id = this.Id

	sql.Value = this.Value
	for _, include := range this.Includes {
		if include.RefId == "" {
			continue
		}
		defSql := fn(include.RefId)
		if defSql == nil {
			continue
		}

		if defSql._isAnalysis {
			sql._isAnalysis = true
		}

		sql.Value += defSql.Value
	}

	return sql
}

func (this *XmlMapper) ParseMapper() *Mapper {

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
		s := sql.Parse(findDefSqlFn)
		s._type = SqlType.Insert
		mapper.Sqls[sql.Id] = s
	}

	//prepare parse delete sql
	for _, sql := range this.Deletes {
		s := sql.Parse(findDefSqlFn)
		s._type = SqlType.Delete
		mapper.Sqls[sql.Id] = s
	}

	//prepare parse update sql
	for _, sql := range this.Updates {
		s := sql.Parse(findDefSqlFn)
		s._type = SqlType.Update
		mapper.Sqls[sql.Id] = s
	}

	//prepare parse delete sql
	for _, sql := range this.Selects {
		s := sql.Parse(findDefSqlFn)
		s._type = SqlType.Select
		mapper.Sqls[sql.Id] = s
	}

	return mapper
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
		return xml.ParseMapper(), nil
	}
}
