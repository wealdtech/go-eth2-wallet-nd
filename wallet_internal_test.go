// Copyright © 2019 Weald Technology Trading
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

package nd

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	keystorev4 "github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4"
	filesystem "github.com/wealdtech/go-eth2-wallet-store-filesystem"
	e2wtypes "github.com/wealdtech/go-eth2-wallet-types/v2"
)

func TestUnmarshalWallet(t *testing.T) {
	tests := []struct {
		name       string
		input      []byte
		err        error
		id         uuid.UUID
		version    uint
		walletType string
	}{
		{
			name: "Nil",
			err:  errors.New("unexpected end of JSON input"),
		},
		{
			name:  "Empty",
			input: []byte{},
			err:   errors.New("unexpected end of JSON input"),
		},
		{
			name:  "NotJSON",
			input: []byte(`bad`),
			err:   errors.New(`invalid character 'b' looking for beginning of value`),
		},
		{
			name:  "MissingID",
			input: []byte(`{"name":"Bad","type":"non-deterministic","version":1}`),
			err:   errors.New("wallet ID missing"),
		},
		{
			name:  "WrongID",
			input: []byte(`{"uuid":7,"name":"Bad","type":"non-deterministic","version":1}`),
			err:   errors.New("wallet ID invalid"),
		},
		{
			name:  "BadID",
			input: []byte(`{"uuid":"bad","name":"Bad","type":"non-deterministic","version":1}`),
			err:   errors.New("failed to parse wallet ID: invalid UUID length: 3"),
		},
		{
			name:  "WrongOldID",
			input: []byte(`{"id":7,"name":"Bad","type":"non-deterministic","version":1}`),
			err:   errors.New("wallet ID invalid"),
		},
		{
			name:  "BadOldID",
			input: []byte(`{"id":"bad","name":"Bad","type":"non-deterministic","version":1}`),
			err:   errors.New("failed to parse wallet ID: invalid UUID length: 3"),
		},
		{
			name:  "MissingName",
			input: []byte(`{"id":"c9958061-63d4-4a80-bcf3-25f3dda22340","type":"non-deterministic","version":1}`),
			err:   errors.New("wallet name missing"),
		},
		{
			name:  "WrongName",
			input: []byte(`{"uuid":"c9958061-63d4-4a80-bcf3-25f3dda22340","name":1,"type":"non-deterministic","version":1}`),
			err:   errors.New("wallet name invalid"),
		},
		{
			name:  "MissingType",
			input: []byte(`{"uuid":"c9958061-63d4-4a80-bcf3-25f3dda22340","name":"Bad","version":1}`),
			err:   errors.New("wallet type missing"),
		},
		{
			name:  "WrongType",
			input: []byte(`{"uuid":"c9958061-63d4-4a80-bcf3-25f3dda22340","name":"Bad","type":7,"version":1}`),
			err:   errors.New("wallet type invalid"),
		},
		{
			name:  "BadType",
			input: []byte(`{"uuid":"c9958061-63d4-4a80-bcf3-25f3dda22340","name":"Bad","type":"hd","version":1}`),
			err:   errors.New(`wallet type "hd" unexpected`),
		},
		{
			name:  "MissingVersion",
			input: []byte(`{"uuid":"c9958061-63d4-4a80-bcf3-25f3dda22340","name":"Bad","type":"non-deterministic"}`),
			err:   errors.New("wallet version missing"),
		},
		{
			name:  "WrongVersion",
			input: []byte(`{"uuid":"c9958061-63d4-4a80-bcf3-25f3dda22340","name":"Bad","type":"non-deterministic","version":"1"}`),
			err:   errors.New("wallet version invalid"),
		},
		{
			name:       "Good",
			input:      []byte(`{"uuid":"c9958061-63d4-4a80-bcf3-25f3dda22340","name":"Good","type":"non-deterministic","version":1}`),
			walletType: "non-deterministic",
			id:         uuid.MustParse("c9958061-63d4-4a80-bcf3-25f3dda22340"),
			version:    1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := newWallet()
			err := json.Unmarshal(test.input, output)
			if test.err != nil {
				require.NotNil(t, err)
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				require.Nil(t, err)
				assert.Equal(t, test.id, output.ID())
				assert.Equal(t, test.version, output.Version())
				assert.Equal(t, test.walletType, output.Type())
			}
		})
	}
}

func TestRetrieveAccountsIndex(t *testing.T) {
	// #nosec G404
	path := filepath.Join(os.TempDir(), fmt.Sprintf("TestRetrieveAccountsIndex-%d", rand.Int31()))
	defer os.RemoveAll(path)
	store := filesystem.New(filesystem.WithLocation(path))
	encryptor := keystorev4.New()
	w, err := CreateWallet(context.Background(), "test wallet", store, encryptor)
	require.NoError(t, err)
	require.NoError(t, w.(e2wtypes.WalletLocker).Unlock(context.Background(), nil))

	account1, err := w.(e2wtypes.WalletAccountCreator).CreateAccount(context.Background(), "account1", []byte("test"))
	require.NoError(t, err)

	account2, err := w.(e2wtypes.WalletAccountCreator).CreateAccount(context.Background(), "account2", []byte("test"))
	require.NoError(t, err)

	idx, found := w.(*wallet).index.ID(account1.Name())
	require.True(t, found)
	require.Equal(t, account1.ID(), idx)

	idx, found = w.(*wallet).index.ID(account2.Name())
	require.True(t, found)
	require.Equal(t, account2.ID(), idx)

	_, found = w.(*wallet).index.ID("not present")
	require.False(t, found)

	// Manually delete the wallet index.
	indexPath := filepath.Join(path, w.ID().String(), "index")
	_, err = os.Stat(indexPath)
	require.False(t, os.IsNotExist(err))
	os.Remove(indexPath)
	_, err = os.Stat(indexPath)
	require.True(t, os.IsNotExist(err))

	// Re-open the wallet with a new store, to force re-creation of the index.
	store = filesystem.New(filesystem.WithLocation(path))
	w, err = OpenWallet(context.Background(), "test wallet", store, encryptor)
	require.NoError(t, err)

	require.NoError(t, w.(*wallet).retrieveAccountsIndex(context.Background()))
	idx, found = w.(*wallet).index.ID(account1.Name())
	require.True(t, found)
	require.Equal(t, account1.ID(), idx)

	idx, found = w.(*wallet).index.ID(account2.Name())
	require.True(t, found)
	require.Equal(t, account2.ID(), idx)

	_, found = w.(*wallet).index.ID("not present")
	require.False(t, found)
}
