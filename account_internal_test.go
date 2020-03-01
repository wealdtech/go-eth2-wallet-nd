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
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalAccount(t *testing.T) {
	tests := []struct {
		name       string
		input      []byte
		err        error
		id         uuid.UUID
		version    uint
		walletType string
		publicKey  []byte
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
			name:  "Blank",
			input: []byte(""),
			err:   errors.New("unexpected end of JSON input"),
		},
		{
			name:  "NotJSON",
			input: []byte(`bad`),
			err:   errors.New(`invalid character 'b' looking for beginning of value`),
		},
		{
			name:  "MissingID",
			input: []byte(`{"name":"test account","pubkey":"a99a76ed7796f7be22d5b7e85deeb7c5677e88e511e0b337618f8c4eb61349b4bf2d153f649f7b53359fe8b94a38e44c","version":4,"crypto":{"checksum":{"function":"sha256","message":"09b65fda487a021900003a8b2081694b15ca73e0e59a5c79a5126f6818a2f171","params":{}},"cipher":{"function":"aes-128-ctr","message":"8386db98fbe002c02de9bc122b7680078045bf6c5c9ac2f7e8b53afbea0d3e15","params":{"iv":"45092570c625ad5e8decfcd991464740"}},"kdf":{"function":"pbkdf2","message":"","params":{"c":16,"dklen":32,"prf":"hmac-sha256","salt":"ae6433afd822e6d99dfaa1a0d73d2ee263efdf62f858ba0c422cf27982d09c8a"}}}}`),
			err:   errors.New("account ID missing"),
		},
		{
			name:  "WrongID",
			input: []byte(`{"uuid":1,"name":"test account","pubkey":"a99a76ed7796f7be22d5b7e85deeb7c5677e88e511e0b337618f8c4eb61349b4bf2d153f649f7b53359fe8b94a38e44c","version":4,"crypto":{"checksum":{"function":"sha256","message":"09b65fda487a021900003a8b2081694b15ca73e0e59a5c79a5126f6818a2f171","params":{}},"cipher":{"function":"aes-128-ctr","message":"8386db98fbe002c02de9bc122b7680078045bf6c5c9ac2f7e8b53afbea0d3e15","params":{"iv":"45092570c625ad5e8decfcd991464740"}},"kdf":{"function":"pbkdf2","message":"","params":{"c":16,"dklen":32,"prf":"hmac-sha256","salt":"ae6433afd822e6d99dfaa1a0d73d2ee263efdf62f858ba0c422cf27982d09c8a"}}}}`),
			err:   errors.New("account ID invalid"),
		},
		{
			name:  "BadID",
			input: []byte(`{"uuid":"c99","name":"test account","pubkey":"a99a76ed7796f7be22d5b7e85deeb7c5677e88e511e0b337618f8c4eb61349b4bf2d153f649f7b53359fe8b94a38e44c","version":4,"crypto":{"checksum":{"function":"sha256","message":"09b65fda487a021900003a8b2081694b15ca73e0e59a5c79a5126f6818a2f171","params":{}},"cipher":{"function":"aes-128-ctr","message":"8386db98fbe002c02de9bc122b7680078045bf6c5c9ac2f7e8b53afbea0d3e15","params":{"iv":"45092570c625ad5e8decfcd991464740"}},"kdf":{"function":"pbkdf2","message":"","params":{"c":16,"dklen":32,"prf":"hmac-sha256","salt":"ae6433afd822e6d99dfaa1a0d73d2ee263efdf62f858ba0c422cf27982d09c8a"}}}}`),
			err:   errors.New("invalid UUID length: 3"),
		},
		{
			name:  "WrongOldID",
			input: []byte(`{"id":1,"name":"test account","pubkey":"a99a76ed7796f7be22d5b7e85deeb7c5677e88e511e0b337618f8c4eb61349b4bf2d153f649f7b53359fe8b94a38e44c","version":4,"crypto":{"checksum":{"function":"sha256","message":"09b65fda487a021900003a8b2081694b15ca73e0e59a5c79a5126f6818a2f171","params":{}},"cipher":{"function":"aes-128-ctr","message":"8386db98fbe002c02de9bc122b7680078045bf6c5c9ac2f7e8b53afbea0d3e15","params":{"iv":"45092570c625ad5e8decfcd991464740"}},"kdf":{"function":"pbkdf2","message":"","params":{"c":16,"dklen":32,"prf":"hmac-sha256","salt":"ae6433afd822e6d99dfaa1a0d73d2ee263efdf62f858ba0c422cf27982d09c8a"}}}}`),
			err:   errors.New("account ID invalid"),
		},
		{
			name:  "BadOldID",
			input: []byte(`{"id":"c99","name":"test account","pubkey":"a99a76ed7796f7be22d5b7e85deeb7c5677e88e511e0b337618f8c4eb61349b4bf2d153f649f7b53359fe8b94a38e44c","version":4,"crypto":{"checksum":{"function":"sha256","message":"09b65fda487a021900003a8b2081694b15ca73e0e59a5c79a5126f6818a2f171","params":{}},"cipher":{"function":"aes-128-ctr","message":"8386db98fbe002c02de9bc122b7680078045bf6c5c9ac2f7e8b53afbea0d3e15","params":{"iv":"45092570c625ad5e8decfcd991464740"}},"kdf":{"function":"pbkdf2","message":"","params":{"c":16,"dklen":32,"prf":"hmac-sha256","salt":"ae6433afd822e6d99dfaa1a0d73d2ee263efdf62f858ba0c422cf27982d09c8a"}}}}`),
			err:   errors.New("invalid UUID length: 3"),
		},
		{
			name:  "MissingName",
			input: []byte(`{"id":"c9958061-63d4-4a80-bcf3-25f3dda22340","pubkey":"a99a76ed7796f7be22d5b7e85deeb7c5677e88e511e0b337618f8c4eb61349b4bf2d153f649f7b53359fe8b94a38e44c","version":4,"crypto":{"checksum":{"function":"sha256","message":"09b65fda487a021900003a8b2081694b15ca73e0e59a5c79a5126f6818a2f171","params":{}},"cipher":{"function":"aes-128-ctr","message":"8386db98fbe002c02de9bc122b7680078045bf6c5c9ac2f7e8b53afbea0d3e15","params":{"iv":"45092570c625ad5e8decfcd991464740"}},"kdf":{"function":"pbkdf2","message":"","params":{"c":16,"dklen":32,"prf":"hmac-sha256","salt":"ae6433afd822e6d99dfaa1a0d73d2ee263efdf62f858ba0c422cf27982d09c8a"}}}}`),
			err:   errors.New("account name missing"),
		},
		{
			name:  "WrongName",
			input: []byte(`{"uuid":"c9958061-63d4-4a80-bcf3-25f3dda22340","name":true,"pubkey":"a99a76ed7796f7be22d5b7e85deeb7c5677e88e511e0b337618f8c4eb61349b4bf2d153f649f7b53359fe8b94a38e44c","version":4,"crypto":{"checksum":{"function":"sha256","message":"09b65fda487a021900003a8b2081694b15ca73e0e59a5c79a5126f6818a2f171","params":{}},"cipher":{"function":"aes-128-ctr","message":"8386db98fbe002c02de9bc122b7680078045bf6c5c9ac2f7e8b53afbea0d3e15","params":{"iv":"45092570c625ad5e8decfcd991464740"}},"kdf":{"function":"pbkdf2","message":"","params":{"c":16,"dklen":32,"prf":"hmac-sha256","salt":"ae6433afd822e6d99dfaa1a0d73d2ee263efdf62f858ba0c422cf27982d09c8a"}}}}`),
			err:   errors.New("account name invalid"),
		},
		{
			name:  "MissingCrypto",
			input: []byte(`{"uuid":"c9958061-63d4-4a80-bcf3-25f3dda22340","name":"test account","pubkey":"a99a76ed7796f7be22d5b7e85deeb7c5677e88e511e0b337618f8c4eb61349b4bf2d153f649f7b53359fe8b94a38e44c","version":4}`),
			err:   errors.New("account crypto missing"),
		},
		{
			name:  "BadCrypto",
			input: []byte(`{"uuid":"c9958061-63d4-4a80-bcf3-25f3dda22340","name":"test account","pubkey":"a99a76ed7796f7be22d5b7e85deeb7c5677e88e511e0b337618f8c4eb61349b4bf2d153f649f7b53359fe8b94a38e44c","version":4,"crypto":2}`),
			err:   errors.New("account crypto invalid"),
		},
		{
			name:  "MissingPubKey",
			input: []byte(`{"uuid":"c9958061-63d4-4a80-bcf3-25f3dda22340","name":"test account","version":4,"crypto":{"checksum":{"function":"sha256","message":"09b65fda487a021900003a8b2081694b15ca73e0e59a5c79a5126f6818a2f171","params":{}},"cipher":{"function":"aes-128-ctr","message":"8386db98fbe002c02de9bc122b7680078045bf6c5c9ac2f7e8b53afbea0d3e15","params":{"iv":"45092570c625ad5e8decfcd991464740"}},"kdf":{"function":"pbkdf2","message":"","params":{"c":16,"dklen":32,"prf":"hmac-sha256","salt":"ae6433afd822e6d99dfaa1a0d73d2ee263efdf62f858ba0c422cf27982d09c8a"}}}}`),
			err:   errors.New("account pubkey missing"),
		},
		{
			name:  "InvalidPubKey",
			input: []byte(`{"uuid":"c9958061-63d4-4a80-bcf3-25f3dda22340","name":"test account","pubkey":true,"version":4,"crypto":{"checksum":{"function":"sha256","message":"09b65fda487a021900003a8b2081694b15ca73e0e59a5c79a5126f6818a2f171","params":{}},"cipher":{"function":"aes-128-ctr","message":"8386db98fbe002c02de9bc122b7680078045bf6c5c9ac2f7e8b53afbea0d3e15","params":{"iv":"45092570c625ad5e8decfcd991464740"}},"kdf":{"function":"pbkdf2","message":"","params":{"c":16,"dklen":32,"prf":"hmac-sha256","salt":"ae6433afd822e6d99dfaa1a0d73d2ee263efdf62f858ba0c422cf27982d09c8a"}}}}`),
			err:   errors.New("account pubkey invalid"),
		},
		{
			name:  "BadPubKey",
			input: []byte(`{"uuid":"c9958061-63d4-4a80-bcf3-25f3dda22340","name":"test account","pubkey":"a99a76ed7796f7be22d5b7e85deeb7c5677e88e511e0b337618f8c4eb61349b4bf2d153f649f7b53359fe8b94a38e44h","version":4,"crypto":{"checksum":{"function":"sha256","message":"09b65fda487a021900003a8b2081694b15ca73e0e59a5c79a5126f6818a2f171","params":{}},"cipher":{"function":"aes-128-ctr","message":"8386db98fbe002c02de9bc122b7680078045bf6c5c9ac2f7e8b53afbea0d3e15","params":{"iv":"45092570c625ad5e8decfcd991464740"}},"kdf":{"function":"pbkdf2","message":"","params":{"c":16,"dklen":32,"prf":"hmac-sha256","salt":"ae6433afd822e6d99dfaa1a0d73d2ee263efdf62f858ba0c422cf27982d09c8a"}}}}`),
			err:   errors.New(`encoding/hex: invalid byte: U+0068 'h'`),
		},
		{
			name:  "BadPubKey2",
			input: []byte(`{"uuid":"c9958061-63d4-4a80-bcf3-25f3dda22340","name":"test account","pubkey":"a99a76ed7796f7be22d5b7e85deeb7c5677e88e511e0b337618f8c4eb61349b4bf2d153f649f7b53359fe8b94a38e44c4c","version":4,"crypto":{"checksum":{"function":"sha256","message":"09b65fda487a021900003a8b2081694b15ca73e0e59a5c79a5126f6818a2f171","params":{}},"cipher":{"function":"aes-128-ctr","message":"8386db98fbe002c02de9bc122b7680078045bf6c5c9ac2f7e8b53afbea0d3e15","params":{"iv":"45092570c625ad5e8decfcd991464740"}},"kdf":{"function":"pbkdf2","message":"","params":{"c":16,"dklen":32,"prf":"hmac-sha256","salt":"ae6433afd822e6d99dfaa1a0d73d2ee263efdf62f858ba0c422cf27982d09c8a"}}}}`),
			err:   errors.New(`public key must be 48 bytes`),
		},
		{
			name:  "MissingVersion",
			input: []byte(`{"uuid":"c9958061-63d4-4a80-bcf3-25f3dda22340","name":"test account","pubkey":"a99a76ed7796f7be22d5b7e85deeb7c5677e88e511e0b337618f8c4eb61349b4bf2d153f649f7b53359fe8b94a38e44c","crypto":{"checksum":{"function":"sha256","message":"09b65fda487a021900003a8b2081694b15ca73e0e59a5c79a5126f6818a2f171","params":{}},"cipher":{"function":"aes-128-ctr","message":"8386db98fbe002c02de9bc122b7680078045bf6c5c9ac2f7e8b53afbea0d3e15","params":{"iv":"45092570c625ad5e8decfcd991464740"}},"kdf":{"function":"pbkdf2","message":"","params":{"c":16,"dklen":32,"prf":"hmac-sha256","salt":"ae6433afd822e6d99dfaa1a0d73d2ee263efdf62f858ba0c422cf27982d09c8a"}}}}`),
			err:   errors.New(`account version missing`),
		},
		{
			name:  "BadVersion",
			input: []byte(`{"uuid":"c9958061-63d4-4a80-bcf3-25f3dda22340","name":"test account","pubkey":"a99a76ed7796f7be22d5b7e85deeb7c5677e88e511e0b337618f8c4eb61349b4bf2d153f649f7b53359fe8b94a38e44c","version":true,"crypto":{"checksum":{"function":"sha256","message":"09b65fda487a021900003a8b2081694b15ca73e0e59a5c79a5126f6818a2f171","params":{}},"cipher":{"function":"aes-128-ctr","message":"8386db98fbe002c02de9bc122b7680078045bf6c5c9ac2f7e8b53afbea0d3e15","params":{"iv":"45092570c625ad5e8decfcd991464740"}},"kdf":{"function":"pbkdf2","message":"","params":{"c":16,"dklen":32,"prf":"hmac-sha256","salt":"ae6433afd822e6d99dfaa1a0d73d2ee263efdf62f858ba0c422cf27982d09c8a"}}}}`),
			err:   errors.New(`account version invalid`),
		},
		{
			name:  "WrongVersion",
			input: []byte(`{"uuid":"c9958061-63d4-4a80-bcf3-25f3dda22340","name":"test account","pubkey":"a99a76ed7796f7be22d5b7e85deeb7c5677e88e511e0b337618f8c4eb61349b4bf2d153f649f7b53359fe8b94a38e44c","version":3,"crypto":{"checksum":{"function":"sha256","message":"09b65fda487a021900003a8b2081694b15ca73e0e59a5c79a5126f6818a2f171","params":{}},"cipher":{"function":"aes-128-ctr","message":"8386db98fbe002c02de9bc122b7680078045bf6c5c9ac2f7e8b53afbea0d3e15","params":{"iv":"45092570c625ad5e8decfcd991464740"}},"kdf":{"function":"pbkdf2","message":"","params":{"c":16,"dklen":32,"prf":"hmac-sha256","salt":"ae6433afd822e6d99dfaa1a0d73d2ee263efdf62f858ba0c422cf27982d09c8a"}}}}`),
			err:   errors.New(`unsupported keystore version`),
		},
		{
			name:       "Good",
			input:      []byte(`{"uuid":"c9958061-63d4-4a80-bcf3-25f3dda22340","name":"test account","pubkey":"a99a76ed7796f7be22d5b7e85deeb7c5677e88e511e0b337618f8c4eb61349b4bf2d153f649f7b53359fe8b94a38e44c","version":4,"crypto":{"checksum":{"function":"sha256","message":"09b65fda487a021900003a8b2081694b15ca73e0e59a5c79a5126f6818a2f171","params":{}},"cipher":{"function":"aes-128-ctr","message":"8386db98fbe002c02de9bc122b7680078045bf6c5c9ac2f7e8b53afbea0d3e15","params":{"iv":"45092570c625ad5e8decfcd991464740"}},"kdf":{"function":"pbkdf2","message":"","params":{"c":16,"dklen":32,"prf":"hmac-sha256","salt":"ae6433afd822e6d99dfaa1a0d73d2ee263efdf62f858ba0c422cf27982d09c8a"}}}}`),
			walletType: "non-deterministic",
			id:         uuid.MustParse("c9958061-63d4-4a80-bcf3-25f3dda22340"),
			publicKey:  []byte{0xa9, 0x9a, 0x76, 0xed, 0x77, 0x96, 0xf7, 0xbe, 0x22, 0xd5, 0xb7, 0xe8, 0x5d, 0xee, 0xb7, 0xc5, 0x67, 0x7e, 0x88, 0xe5, 0x11, 0xe0, 0xb3, 0x37, 0x61, 0x8f, 0x8c, 0x4e, 0xb6, 0x13, 0x49, 0xb4, 0xbf, 0x2d, 0x15, 0x3f, 0x64, 0x9f, 0x7b, 0x53, 0x35, 0x9f, 0xe8, 0xb9, 0x4a, 0x38, 0xe4, 0x4c},
			version:    4,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := newAccount()
			err := json.Unmarshal(test.input, output)
			if test.err != nil {
				require.NotNil(t, err)
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				require.Nil(t, err)
				assert.Equal(t, test.id, output.ID())
				assert.Equal(t, test.publicKey, output.PublicKey().Marshal())
				//				assert.Equal(t, test.version, output.Version())
				//				assert.Equal(t, test.walletType, output.Type())
			}
		})
	}
}

func TestUnlock(t *testing.T) {
	tests := []struct {
		name       string
		account    []byte
		passphrase []byte
		err        error
	}{
		{
			name:       "PublicKeyMismatch",
			account:    []byte(`{"uuid":"c9958061-63d4-4a80-bcf3-25f3dda22340","name":"test account","pubkey":"b89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","version":4,"crypto":{"checksum":{"function":"sha256","message":"09b65fda487a021900003a8b2081694b15ca73e0e59a5c79a5126f6818a2f171","params":{}},"cipher":{"function":"aes-128-ctr","message":"8386db98fbe002c02de9bc122b7680078045bf6c5c9ac2f7e8b53afbea0d3e15","params":{"iv":"45092570c625ad5e8decfcd991464740"}},"kdf":{"function":"pbkdf2","message":"","params":{"c":16,"dklen":32,"prf":"hmac-sha256","salt":"ae6433afd822e6d99dfaa1a0d73d2ee263efdf62f858ba0c422cf27982d09c8a"}}}}`),
			passphrase: []byte("test passphrase"),
			err:        errors.New("secret key does not correspond to public key"),
		},
		{
			name:       "Keystore",
			account:    []byte(`{"uuid":"c9958061-63d4-4a80-bcf3-25f3dda22340","name":"test account","pubkey":"a99a76ed7796f7be22d5b7e85deeb7c5677e88e511e0b337618f8c4eb61349b4bf2d153f649f7b53359fe8b94a38e44c","version":4,"crypto":{"checksum":{"function":"sha256","message":"09b65fda487a021900003a8b2081694b15ca73e0e59a5c79a5126f6818a2f171","params":{}},"cipher":{"function":"aes-128-ctr","message":"8386db98fbe002c02de9bc122b7680078045bf6c5c9ac2f7e8b53afbea0d3e15","params":{"iv":"45092570c625ad5e8decfcd991464740"}},"kdf":{"function":"pbkdf2","message":"","params":{"c":16,"dklen":32,"prf":"hmac-sha256","salt":"ae6433afd822e6d99dfaa1a0d73d2ee263efdf62f858ba0c422cf27982d09c8a"}}}}`),
			passphrase: []byte("test passphrase"),
		},
		{
			name:       "BadPassphrase",
			account:    []byte(`{"uuid":"c9958061-63d4-4a80-bcf3-25f3dda22340","name":"test account","pubkey":"a99a76ed7796f7be22d5b7e85deeb7c5677e88e511e0b337618f8c4eb61349b4bf2d153f649f7b53359fe8b94a38e44c","version":4,"crypto":{"checksum":{"function":"sha256","message":"09b65fda487a021900003a8b2081694b15ca73e0e59a5c79a5126f6818a2f171","params":{}},"cipher":{"function":"aes-128-ctr","message":"8386db98fbe002c02de9bc122b7680078045bf6c5c9ac2f7e8b53afbea0d3e15","params":{"iv":"45092570c625ad5e8decfcd991464740"}},"kdf":{"function":"pbkdf2","message":"","params":{"c":16,"dklen":32,"prf":"hmac-sha256","salt":"ae6433afd822e6d99dfaa1a0d73d2ee263efdf62f858ba0c422cf27982d09c8a"}}}}`),
			passphrase: []byte("wrong passphrase"),
			err:        errors.New("incorrect passphrase"),
		},
		{
			name:       "EmptyPassphrase",
			account:    []byte(`{"uuid":"c9958061-63d4-4a80-bcf3-25f3dda22340","name":"test account","pubkey":"a99a76ed7796f7be22d5b7e85deeb7c5677e88e511e0b337618f8c4eb61349b4bf2d153f649f7b53359fe8b94a38e44c","version":4,"crypto":{"checksum":{"function":"sha256","message":"4a67cc6a4ff5e81235393c677652213cc96488d68f17d045f99f9cef8acc81a1","params":{}},"cipher":{"function":"aes-128-ctr","message":"ce7c1d11cd71adb604c055a2d198336387e0579275c4d2d45c184ed54631ebdd","params":{"iv":"c752efc43ca0651bb06adccf4b8651b8"}},"kdf":{"function":"pbkdf2","message":"","params":{"c":16,"dklen":32,"prf":"hmac-sha256","salt":"b49107e74e59a80ce5ac1624e6d27e7305aa22f5ffba4f602dd4dfe34fdf8640"}}}}`),
			passphrase: []byte(""),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			account := newAccount()
			err := json.Unmarshal(test.account, account)
			require.Nil(t, err)

			// Try to sign something - should fail because locked
			_, err = account.Sign([]byte("test"))
			assert.NotNil(t, err)

			err = account.Unlock(test.passphrase)
			if test.err != nil {
				require.NotNil(t, err)
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				require.Nil(t, err)

				// Try to sign something - should succeed because unlocked
				signature, err := account.Sign([]byte("test"))
				assert.Nil(t, err)

				verified := signature.Verify([]byte("test"), account.PublicKey())
				assert.Equal(t, true, verified)

				account.Lock()

				// Try to sign something - should fail because locked (again)
				_, err = account.Sign([]byte("test"))
				assert.NotNil(t, err)
			}
		})
	}
}
