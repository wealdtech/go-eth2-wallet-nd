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
	"context"
	"encoding/hex"
	"fmt"
	"math/rand"
	"sync"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	keystorev4 "github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4"
	nd "github.com/wealdtech/go-eth2-wallet-nd/v2"
	scratch "github.com/wealdtech/go-eth2-wallet-store-scratch"
	e2wtypes "github.com/wealdtech/go-eth2-wallet-types/v2"
)

// _byteArray is a helper to turn a string in to a byte array
func _byteArray(input string) []byte {
	x, _ := hex.DecodeString(input)
	return x
}

func TestAccountInterfaces(t *testing.T) {
	store := scratch.New()
	encryptor := keystorev4.New()
	wallet, err := nd.CreateWallet(context.Background(), "test wallet", store, encryptor)
	require.NoError(t, err)
	locker, isLocker := wallet.(e2wtypes.WalletLocker)
	require.True(t, isLocker)
	err = locker.Unlock(context.Background(), nil)
	require.NoError(t, err)

	account, err := wallet.(e2wtypes.WalletAccountCreator).CreateAccount(context.Background(), "account", []byte("test"))
	require.NoError(t, err)

	require.Implements(t, (*e2wtypes.Account)(nil), account)
	require.Implements(t, (*e2wtypes.AccountIDProvider)(nil), account)
	require.NotEmpty(t, account.(e2wtypes.AccountIDProvider).ID())
	require.Implements(t, (*e2wtypes.AccountNameProvider)(nil), account)
	require.NotEmpty(t, account.(e2wtypes.AccountNameProvider).Name())
	require.Implements(t, (*e2wtypes.AccountPublicKeyProvider)(nil), account)
	require.NotEmpty(t, account.(e2wtypes.AccountPublicKeyProvider).PublicKey())
	require.Implements(t, (*e2wtypes.AccountPathProvider)(nil), account)
	// ND-wallet returns an empty path.
	require.Empty(t, account.(e2wtypes.AccountPathProvider).Path())
	require.Implements(t, (*e2wtypes.AccountWalletProvider)(nil), account)
	require.NotEmpty(t, account.(e2wtypes.AccountWalletProvider).Wallet())
	require.Implements(t, (*e2wtypes.AccountLocker)(nil), account)
	require.Implements(t, (*e2wtypes.AccountSigner)(nil), account)
	require.Implements(t, (*e2wtypes.AccountPrivateKeyProvider)(nil), account)
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
	wallet, err := nd.CreateWallet(context.Background(), "test wallet", store, encryptor)
	require.Nil(t, err)

	// Try to create without unlocking the wallet; should fail.
	_, err = wallet.(e2wtypes.WalletAccountCreator).CreateAccount(context.Background(), "attempt", []byte("test"))
	assert.NotNil(t, err)

	locker, isLocker := wallet.(e2wtypes.WalletLocker)
	require.True(t, isLocker)
	err = locker.Unlock(context.Background(), nil)
	require.Nil(t, err)
	defer func(t *testing.T) {
		require.NoError(t, locker.Lock(context.Background()))
	}(t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			account, err := wallet.(e2wtypes.WalletAccountCreator).CreateAccount(context.Background(), test.accountName, test.passphrase)
			if test.err != nil {
				require.NotNil(t, err)
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				require.Nil(t, err)
				assert.Equal(t, test.accountName, account.Name())
			}
		})
	}
}

func TestAccountUnlockLock(t *testing.T) {
	store := scratch.New()
	encryptor := keystorev4.New()
	wallet, err := nd.CreateWallet(context.Background(), "test wallet", store, encryptor)
	require.NoError(t, err)
	require.NoError(t, wallet.(e2wtypes.WalletLocker).Unlock(context.Background(), nil))

	account, err := wallet.(e2wtypes.WalletAccountCreator).CreateAccount(context.Background(), "account", []byte("pass"))
	require.NoError(t, err)
	unlocked, err := account.(e2wtypes.AccountLocker).IsUnlocked(context.Background())
	require.NoError(t, err)
	require.False(t, unlocked)
	// Unlock account.
	require.NoError(t, account.(e2wtypes.AccountLocker).Unlock(context.Background(), []byte("pass")))
	unlocked, err = account.(e2wtypes.AccountLocker).IsUnlocked(context.Background())
	require.NoError(t, err)
	require.True(t, unlocked)
	// Unlock account again.
	require.NoError(t, account.(e2wtypes.AccountLocker).Unlock(context.Background(), []byte("pass")))
	unlocked, err = account.(e2wtypes.AccountLocker).IsUnlocked(context.Background())
	require.NoError(t, err)
	require.True(t, unlocked)
	// Lock account.
	require.NoError(t, account.(e2wtypes.AccountLocker).Lock(context.Background()))
	unlocked, err = account.(e2wtypes.AccountLocker).IsUnlocked(context.Background())
	require.NoError(t, err)
	require.False(t, unlocked)
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
	wallet, err := nd.CreateWallet(context.Background(), "test wallet", store, encryptor)
	require.Nil(t, err)

	// Try to import without unlocking the wallet; should fail
	_, err = wallet.(e2wtypes.WalletAccountImporter).ImportAccount(context.Background(), "attempt", _byteArray("220091d10843519cd1c452a4ec721d378d7d4c5ece81c4b5556092d410e5e0e1"), []byte("test"))
	assert.NotNil(t, err)

	locker, isLocker := wallet.(e2wtypes.WalletLocker)
	require.True(t, isLocker)
	err = locker.Unlock(context.Background(), nil)
	require.Nil(t, err)
	defer func(t *testing.T) {
		require.NoError(t, locker.Lock(context.Background()))
	}(t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			account, err := wallet.(e2wtypes.WalletAccountImporter).ImportAccount(context.Background(), test.accountName, test.key, test.passphrase)
			if test.err != nil {
				require.NotNil(t, err)
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				require.Nil(t, err)
				assert.Equal(t, test.accountName, account.Name())
				// Should not be able to obtain private key from a locked account
				_, err = account.(e2wtypes.AccountPrivateKeyProvider).PrivateKey(context.Background())
				assert.NotNil(t, err)
				locker, isLocker := account.(e2wtypes.AccountLocker)
				require.True(t, isLocker)
				err = locker.Unlock(context.Background(), test.passphrase)
				require.Nil(t, err)
				_, err := account.(e2wtypes.AccountPrivateKeyProvider).PrivateKey(context.Background())
				assert.Nil(t, err)
			}
		})
	}
}

func TestConcurrentCreate(t *testing.T) {
	store := scratch.New()
	encryptor := keystorev4.New()
	wallet, err := nd.CreateWallet(context.Background(), "test wallet", store, encryptor)
	require.NoError(t, err)
	locker, isLocker := wallet.(e2wtypes.WalletLocker)
	require.True(t, isLocker)
	require.NoError(t, locker.Unlock(context.Background(), nil))

	// Create a number of runners that will try to create accounts simultaneously.
	var runWG sync.WaitGroup
	var setupWG sync.WaitGroup
	starter := make(chan any)
	numAccounts := 64
	for i := 0; i < numAccounts; i++ {
		setupWG.Add(1)
		runWG.Add(1)
		go func() {
			id := rand.Uint32()
			name := fmt.Sprintf("Test account %d", id)
			setupWG.Done()

			<-starter

			account, err := wallet.(e2wtypes.WalletAccountCreator).CreateAccount(context.Background(), name, []byte("test"))
			require.NoError(t, err)
			require.NotNil(t, account)
			runWG.Done()
		}()
	}

	// Wait for setup to complete.
	setupWG.Wait()

	// Start the jobs by closing the channel.
	close(starter)

	// Wait for run to complete
	runWG.Wait()

	// Confirm that all accounts have been created.
	wallet, err = nd.OpenWallet(context.Background(), "test wallet", store, encryptor)
	require.NoError(t, err)
	accounts := 0
	for range wallet.Accounts(context.Background()) {
		accounts++
	}
	require.Equal(t, numAccounts, accounts)
}
