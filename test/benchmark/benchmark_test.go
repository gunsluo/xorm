package main

import (
	"testing"

	"github.com/gunsluo/xorm"

	_ "github.com/go-sql-driver/mysql"
)

func engine() (*xorm.Engine, error) {
	engine, err := xorm.NewEngine("mysql", "root:password@tcp(127.0.0.1:3306)/example?charset=utf8&loc=Asia%2FShanghai")
	if err != nil {
		return nil, err
	}

	engine.SetSqlMapRootDir("./sqlmap")
	err = engine.InitSqlMap()
	if err != nil {
		return nil, err
	}

	return engine, nil
}

func TestSelectEntryList(t *testing.T) {
	e, err := engine()
	if err != nil {
		t.Error(err.Error())
	}
	_, err = selectEntryList(e)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestSelectEntryAll(t *testing.T) {
	e, err := engine()
	if err != nil {
		t.Error(err.Error())
	}
	_, err = selectEntryAll(e)
	if err != nil {
		t.Error(err.Error())
	}
}

func BenchmarkSelectEntryList(b *testing.B) {
	b.ReportAllocs()
	e, err := engine()
	if err != nil {
		b.Error(err.Error())
	}

	for i := 0; i < b.N; i++ {
		_, err = selectEntryList(e)
		if err != nil {
			b.Error(err.Error())
		}
	}
}

func BenchmarkSelectEntryList2(b *testing.B) {
	b.ReportAllocs()
	e, err := engine()
	if err != nil {
		b.Error(err.Error())
	}

	for i := 0; i < b.N; i++ {
		_, err = selectEntryList2(e)
		if err != nil {
			b.Error(err.Error())
		}
	}
}

func BenchmarkSelectEntryAll(b *testing.B) {
	b.ReportAllocs()
	e, err := engine()
	if err != nil {
		b.Error(err.Error())
	}

	for i := 0; i < b.N; i++ {
		_, err = selectEntryAll(e)
		if err != nil {
			b.Error(err.Error())
		}
	}
}

func BenchmarkSelectEntryAll2(b *testing.B) {
	b.ReportAllocs()
	e, err := engine()
	if err != nil {
		b.Error(err.Error())
	}

	for i := 0; i < b.N; i++ {
		_, err = selectEntryAll2(e)
		if err != nil {
			b.Error(err.Error())
		}
	}
}
