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

package nd_test

import (
	"encoding/hex"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	keystorev4 "github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4"
	nd "github.com/wealdtech/go-eth2-wallet-nd/v2"
	scratch "github.com/wealdtech/go-eth2-wallet-store-scratch"
	wtypes "github.com/wealdtech/go-eth2-wallet-types/v2"
)

// _byteArray is a helper to turn a string in to a byte array
func _byteArray(input string) []byte {
	x, _ := hex.DecodeString(input)
	return x
}

func TestCreateAccount(t *testing.T) {
	tests := []struct {
		name        string
		accountName string
		passphrase  []byte
		err         error
	}{
		{
			name:        "Empty",
			accountName: "",
			err:         errors.New("account name missing"),
		},
		{
			name:        "Invalid",
			accountName: "_bad",
			err:         errors.New(`invalid account name "_bad"`),
		},
		{
			name:        "Good",
			accountName: "test",
		},
		{
			name:        "Duplicate",
			accountName: "test",
			err:         errors.New(`account with name "test" already exists`),
		},
	}

	store := scratch.New()
	encryptor := keystorev4.New()
	wallet, err := nd.CreateWallet("test wallet", store, encryptor)
	require.Nil(t, err)

	// Try to create without unlocking the wallet; should fail
	_, err = wallet.(wtypes.WalletAccountCreator).CreateAccount("attempt", []byte("test"))
	assert.NotNil(t, err)

	err = wallet.Unlock(nil)
	require.Nil(t, err)
	defer wallet.Lock()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			account, err := wallet.(wtypes.WalletAccountCreator).CreateAccount(test.accountName, test.passphrase)
			if test.err != nil {
				require.NotNil(t, err)
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				require.Nil(t, err)
				assert.Equal(t, test.accountName, account.Name())
				assert.Equal(t, "", account.Path())
				//				assert.Equal(t, test.id, account.ID())
				//				assert.Equal(t, test.version, account.Version())
				//				assert.Equal(t, test.walletType, account.Type())
			}
		})
	}
}

func TestImportAccount(t *testing.T) {
	tests := []struct {
		name        string
		accountName string
		key         []byte
		passphrase  []byte
		err         error
	}{
		{
			name:        "Empty",
			accountName: "",
			err:         errors.New("account name missing"),
		},
		{
			name:        "Invalid",
			accountName: "_bad",
			err:         errors.New(`invalid account name "_bad"`),
		},
		{
			name:        "Good",
			key:         _byteArray("220091d10843519cd1c452a4ec721d378d7d4c5ece81c4b5556092d410e5e0e1"),
			accountName: "test",
		},
		{
			name:        "Duplicate",
			accountName: "test",
			err:         errors.New(`account with name "test" already exists`),
		},
	}

	store := scratch.New()
	encryptor := keystorev4.New()
	wallet, err := nd.CreateWallet("test wallet", store, encryptor)
	require.Nil(t, err)

	// Try to import without unlocking the wallet; should fail
	_, err = wallet.(wtypes.WalletAccountImporter).ImportAccount("attempt", _byteArray("220091d10843519cd1c452a4ec721d378d7d4c5ece81c4b5556092d410e5e0e1"), []byte("test"))
	assert.NotNil(t, err)

	err = wallet.Unlock(nil)
	require.Nil(t, err)
	defer wallet.Lock()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			account, err := wallet.(wtypes.WalletAccountImporter).ImportAccount(test.accountName, test.key, test.passphrase)
			if test.err != nil {
				require.NotNil(t, err)
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				require.Nil(t, err)
				assert.Equal(t, test.accountName, account.Name())
				assert.Equal(t, "", account.Path())
				// Should not be able to obtain private key from a locked account
				_, err = account.(wtypes.AccountPrivateKeyProvider).PrivateKey()
				assert.NotNil(t, err)
				err = account.Unlock(test.passphrase)
				require.Nil(t, err)
				_, err := account.(wtypes.AccountPrivateKeyProvider).PrivateKey()
				assert.Nil(t, err)
			}
		})
	}
}
