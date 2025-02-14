package vertigo

// Copyright (c) 2019 Micro Focus or one of its affiliates.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
)

type tx struct {
	conn    *connection
	context context.Context
}

func (t *tx) Commit() error {
	stmt, err := t.conn.PrepareContext(t.context, "COMMIT")

	if err != nil {
		return err
	}

	_, err = stmt.Exec([]driver.Value{})

	if err != nil {
		return err
	}

	return nil
}

func (t *tx) Rollback() error {
	stmt, err := t.conn.PrepareContext(t.context, "ROLLBACK")

	if err != nil {
		return err
	}

	_, err = stmt.Exec([]driver.Value{})

	if err != nil {
		return err
	}

	return nil

}

func newTransaction(ctx context.Context, c *connection, opts driver.TxOptions) (*tx, error) {

	queryStr := "START TRANSACTION"

	switch sql.IsolationLevel(opts.Isolation) {
	case sql.LevelReadUncommitted:
		queryStr += " READ UNCOMMITTED"
	case sql.LevelReadCommitted:
		queryStr += " READ COMMITTED"
	case sql.LevelSerializable:
		queryStr += " SERIALIZABLE"
	case sql.LevelRepeatableRead:
		queryStr += " REPEATABLE READ"
	case sql.LevelDefault:
		break
	default:
		return nil, fmt.Errorf("unsupported transaction isolation level: %s", sql.IsolationLevel(opts.Isolation))
	}

	if opts.ReadOnly {
		queryStr += " READ ONLY"
	} else {
		queryStr += " READ WRITE"
	}

	stmt, err := c.PrepareContext(ctx, queryStr)

	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec([]driver.Value{})

	if err != nil {
		return nil, err
	}

	res := &tx{
		conn:    c,
		context: ctx,
	}

	return res, nil
}
