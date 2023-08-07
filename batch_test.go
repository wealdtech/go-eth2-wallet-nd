// Copyright © 2023 Weald Technology Trading.
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
	"testing"

	"github.com/stretchr/testify/require"
	keystorev4 "github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4"
	nd "github.com/wealdtech/go-eth2-wallet-nd/v2"
	scratch "github.com/wealdtech/go-eth2-wallet-store-scratch"
	e2wtypes "github.com/wealdtech/go-eth2-wallet-types/v2"
)

func TestBatch(t *testing.T) {
	ctx := context.Background()
	store := scratch.New()
	encryptor := keystorev4.New()

	// Create a wallet.
	wallet, err := nd.CreateWallet(ctx, "test wallet", store, encryptor)
	require.NoError(t, err)
	require.NoError(t, wallet.(e2wtypes.WalletLocker).Unlock(ctx, nil))

	// Add some accounts.
	account1, err := wallet.(e2wtypes.WalletAccountCreator).CreateAccount(ctx, "account 1", []byte("passphrase"))
	require.NoError(t, err)
	account2, err := wallet.(e2wtypes.WalletAccountCreator).CreateAccount(ctx, "account 2", []byte("passphrase"))
	require.NoError(t, err)

	// Create a batch.
	require.NoError(t, wallet.(e2wtypes.WalletBatchCreator).BatchWallet(ctx, []string{"passphrase"}, "batch passphrase"))

	// Re-open the wallet and fetch the accounts through the batch system.
	wallet, err = nd.OpenWallet(ctx, "test wallet", store, encryptor)
	require.NoError(t, err)
	numAccounts := 0
	for range wallet.Accounts(ctx) {
		numAccounts++
	}
	require.Equal(t, 2, numAccounts)
	obtainedAccount1, err := wallet.(e2wtypes.WalletAccountByNameProvider).AccountByName(ctx, "account 1")
	require.NoError(t, err)
	require.Equal(t, account1.ID(), obtainedAccount1.ID())
	obtainedAccount2, err := wallet.(e2wtypes.WalletAccountByIDProvider).AccountByID(ctx, account2.ID())
	require.NoError(t, err)
	require.Equal(t, account2.Name(), obtainedAccount2.Name())

	// Ensure we can unlock accounts with the batch passphrase.
	require.NoError(t, obtainedAccount1.(e2wtypes.AccountLocker).Unlock(ctx, []byte("batch passphrase")))
	require.NoError(t, obtainedAccount2.(e2wtypes.AccountLocker).Unlock(ctx, []byte("batch passphrase")))

	// Create another account, not in the batch.
	require.NoError(t, wallet.(e2wtypes.WalletLocker).Unlock(ctx, nil))
	account3, err := wallet.(e2wtypes.WalletAccountCreator).CreateAccount(ctx, "account 3", []byte("passphrase"))
	require.NoError(t, err)

	// Re-open the wallet and fetch the non-batch account by name.
	wallet, err = nd.OpenWallet(ctx, "test wallet", store, encryptor)
	require.NoError(t, err)
	numAccounts = 0
	for range wallet.Accounts(ctx) {
		numAccounts++
	}
	require.Equal(t, 2, numAccounts)
	obtainedAccount3, err := wallet.(e2wtypes.WalletAccountByNameProvider).AccountByName(ctx, "account 3")
	require.NoError(t, err)
	require.Equal(t, account3.ID(), obtainedAccount3.ID())

	// Re-open the wallet and fetch the non-batch account by ID.
	wallet, err = nd.OpenWallet(ctx, "test wallet", store, encryptor)
	require.NoError(t, err)
	numAccounts = 0
	for range wallet.Accounts(ctx) {
		numAccounts++
	}
	require.Equal(t, 2, numAccounts)
	obtainedAccount3, err = wallet.(e2wtypes.WalletAccountByIDProvider).AccountByID(ctx, account3.ID())
	require.NoError(t, err)
	require.Equal(t, account3.Name(), obtainedAccount3.Name())

	// Ensure we can unlock account with the account passphrase.
	require.NoError(t, obtainedAccount3.(e2wtypes.AccountLocker).Unlock(ctx, []byte("passphrase")))

	// Recreate the batch.
	require.NoError(t, wallet.(e2wtypes.WalletBatchCreator).BatchWallet(ctx, []string{"passphrase", "batch passphrase"}, "batch passphrase"))

	// Re-open the wallet and fetch the accounts through the batch system.
	wallet, err = nd.OpenWallet(ctx, "test wallet", store, encryptor)
	require.NoError(t, err)
	numAccounts = 0
	for range wallet.Accounts(ctx) {
		numAccounts++
	}
	require.Equal(t, 3, numAccounts)
}
