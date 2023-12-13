// Copyright 2021 FerretDB Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sqlite

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/FerretDB/FerretDB/internal/backends"
	"github.com/FerretDB/FerretDB/internal/clientconn/conninfo"
	"github.com/FerretDB/FerretDB/internal/util/state"
	"github.com/FerretDB/FerretDB/internal/util/testutil"
)

func TestListDatabases(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in -short mode")
	}

	t.Parallel()

	ctx := conninfo.Ctx(testutil.Ctx(t), conninfo.New())

	sp, err := state.NewProvider("")
	require.NoError(t, err)

	params := NewBackendParams{
		URI: testutil.TestPostgreSQLURI(t, ctx, ""),
		L:   testutil.Logger(t),
		P:   sp,
	}
	b, err := NewBackend(&params)
	require.NoError(t, err)
	t.Cleanup(b.Close)

	dbName := testutil.DatabaseName(t)
	_, err = b.Database(dbName)
	require.NoError(t, err)

	t.Run("ListDatabases with specific name", func(t *testing.T) {
		res, err := b.ListDatabases(ctx, &backends.ListDatabasesParams{
			Name: dbName,
		})
		require.NoError(t, err)
		require.Equal(t, 1, len(res.Databases))
		require.Equal(t, dbName, res.Databases[0].Name)
	})

	t.Run("ListDatabases with wrong name", func(t *testing.T) {
		res, err := b.ListDatabases(ctx, &backends.ListDatabasesParams{
			Name: "not-existing",
		})
		require.NoError(t, err)
		require.Equal(t, 0, len(res.Databases))
	})

	t.Run("ListDatabases with nil param", func(t *testing.T) {
		res, err := b.ListDatabases(ctx, nil)
		require.NoError(t, err)
		require.Equal(t, 1, len(res.Databases))
		require.Equal(t, dbName, res.Databases[0].Name)
	})

	t.Run("ListDatabases with nil param", func(t *testing.T) {
		res, err := b.ListDatabases(ctx, &backends.ListDatabasesParams{})
		require.NoError(t, err)
		require.Equal(t, 1, len(res.Databases))
		require.Equal(t, dbName, res.Databases[0].Name)
	})
}
