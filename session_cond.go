// Copyright 2017 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xorm

import "github.com/go-xorm/builder"

// Sql provides raw sql input parameter. When you have a complex SQL statement
// and cannot use Where, Id, In and etc. Methods to describe, you can use SQL.
//
// Deprecated: use SQL instead.
func (session *Session) Sql(query string, args ...interface{}) *Session {
	return session.SQL(query, args...)
}

// SQL provides raw sql input parameter. When you have a complex SQL statement
// and cannot use Where, Id, In and etc. Methods to describe, you can use SQL.
func (session *Session) SQL(query interface{}, args ...interface{}) *Session {
	session.IsSqlFunc = true
	session.Statement.SQL(query, args...)
	return session
}

// Where provides custom query condition.
func (session *Session) Where(query interface{}, args ...interface{}) *Session {
	session.Statement.Where(query, args...)
	return session
}

// And provides custom query condition.
func (session *Session) And(query interface{}, args ...interface{}) *Session {
	session.Statement.And(query, args...)
	return session
}

// Or provides custom query condition.
func (session *Session) Or(query interface{}, args ...interface{}) *Session {
	session.Statement.Or(query, args...)
	return session
}

// Id provides converting id as a query condition
//
// Deprecated: use ID instead
func (session *Session) Id(id interface{}) *Session {
	return session.ID(id)
}

// ID provides converting id as a query condition
func (session *Session) ID(id interface{}) *Session {
	session.Statement.ID(id)
	return session
}

// In provides a query string like "id in (1, 2, 3)"
func (session *Session) In(column string, args ...interface{}) *Session {
	session.Statement.In(column, args...)
	return session
}

// NotIn provides a query string like "id in (1, 2, 3)"
func (session *Session) NotIn(column string, args ...interface{}) *Session {
	session.Statement.NotIn(column, args...)
	return session
}

// Conds returns session query conditions
func (session *Session) Conds() builder.Cond {
	return session.Statement.cond
}