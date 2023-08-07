// Copyright 2019 - 2023 Weald Technology Trading.
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
	"bytes"
	"context"
	"encoding/json"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	e2types "github.com/wealdtech/go-eth2-types/v2"
	e2wtypes "github.com/wealdtech/go-eth2-wallet-types/v2"
)

// account contains the details of the account.
type account struct {
	id        uuid.UUID
	name      string
	publicKey e2types.PublicKey
	crypto    map[string]any
	unlocked  bool
	secretKey e2types.PrivateKey
	version   uint
	wallet    *wallet
	encryptor e2wtypes.Encryptor
	mutex     sync.Mutex
}

// newAccount creates a new account.
func newAccount() *account {
	return &account{}
}

// ID provides the ID for the account.
func (a *account) ID() uuid.UUID {
	return a.id
}

// Name provides the ID for the account.
func (a *account) Name() string {
	return a.name
}

// PublicKey provides the public key for the account.
func (a *account) PublicKey() e2types.PublicKey {
	return a.publicKey
}

// PrivateKey provides the private key for the account.
func (a *account) PrivateKey(_ context.Context) (e2types.PrivateKey, error) {
	if !a.unlocked {
		return nil, errors.New("cannot provide private key when account is locked")
	}

	return a.secretKey, nil
}

// Wallet provides the wallet for the account.
func (a *account) Wallet() e2wtypes.Wallet {
	return a.wallet
}

// Lock locks the account.  A locked account cannot sign data.
func (a *account) Lock(_ context.Context) error {
	a.mutex.Lock()
	a.unlocked = false
	a.mutex.Unlock()

	return nil
}

// Unlock unlocks the account.  An unlocked account can sign data.
func (a *account) Unlock(ctx context.Context, passphrase []byte) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// If the account is already unlocked then nothing to do.
	if a.unlocked {
		return nil
	}

	if a.secretKey == nil {
		// First time unlocking, need to decrypt.
		if a.crypto == nil {
			// This is a batch account, decrypt the batch.
			if err := a.wallet.batchDecrypt(ctx, passphrase); err != nil {
				return errors.Wrap(err, "failed to decrypt batch")
			}
		} else {
			// This is an individual account, decrypt the account.
			privateKeyBytes, err := a.encryptor.Decrypt(a.crypto, string(passphrase))
			if err != nil {
				return errors.New("incorrect passphrase")
			}
			privateKey, err := e2types.BLSPrivateKeyFromBytes(privateKeyBytes)
			if err != nil {
				return errors.Wrap(err, "failed to obtain private key")
			}
			a.secretKey = privateKey
		}

		// Ensure the private key is correct.
		publicKey := a.secretKey.PublicKey()
		if !bytes.Equal(publicKey.Marshal(), a.publicKey.Marshal()) {
			a.secretKey = nil
			return errors.New("private key does not correspond to public key")
		}
	}

	a.unlocked = true

	return nil
}

// IsUnlocked returns true if the account is unlocked.
func (a *account) IsUnlocked(_ context.Context) (bool, error) {
	return a.unlocked, nil
}

// Path returns "" as non-deterministic accounts are not derived.
func (a *account) Path() string {
	return ""
}

// Sign signs data.
func (a *account) Sign(ctx context.Context, data []byte) (e2types.Signature, error) {
	a.mutex.Lock()
	unlocked, err := a.IsUnlocked(ctx)
	if err != nil {
		a.mutex.Unlock()
		return nil, err
	}
	if !unlocked {
		a.mutex.Unlock()
		return nil, errors.New("cannot sign when account is locked")
	}
	if a.secretKey == nil {
		a.mutex.Unlock()
		return nil, errors.New("missing private key for unlocked account")
	}
	a.mutex.Unlock()

	return a.secretKey.Sign(data), nil
}

// storeAccount stores the account.
func (a *account) storeAccount(ctx context.Context) error {
	data, err := json.Marshal(a)
	if err != nil {
		return errors.Wrap(err, "failed to marshal account")
	}
	if err := a.wallet.storeAccountsIndex(); err != nil {
		return errors.Wrap(err, "failed to store account index")
	}
	if err := a.wallet.store.StoreAccount(a.wallet.ID(), a.ID(), data); err != nil {
		return errors.Wrap(err, "failed to store account")
	}

	// Check to ensure the created account can be retrieved.
	if _, err = a.wallet.AccountByName(ctx, a.name); err != nil {
		return errors.Wrap(err, "failed to confirm account when retrieving by name")
	}
	if _, err = a.wallet.AccountByID(ctx, a.id); err != nil {
		return errors.Wrap(err, "failed to confirm account when retrieveing by ID")
	}

	return nil
}

// deserializeAccount deserializes account data to an account.
func deserializeAccount(w *wallet, data []byte) (*account, error) {
	a := newAccount()
	a.wallet = w
	a.encryptor = w.encryptor
	if err := json.Unmarshal(data, a); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal account")
	}

	return a, nil
}
