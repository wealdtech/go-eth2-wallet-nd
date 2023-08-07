# go-eth2-wallet-nd

[![Tag](https://img.shields.io/github/tag/wealdtech/go-eth2-wallet-nd.svg)](https://github.com/wealdtech/go-eth2-wallet-nd/releases/)
[![License](https://img.shields.io/github/license/wealdtech/go-eth2-wallet-nd.svg)](LICENSE)
[![GoDoc](https://godoc.org/github.com/wealdtech/go-eth2-wallet-nd?status.svg)](https://godoc.org/github.com/wealdtech/go-eth2-wallet-nd)
[![Go Report Card](https://goreportcard.com/badge/github.com/wealdtech/go-eth2-wallet-nd)](https://goreportcard.com/report/github.com/wealdtech/go-eth2-wallet-nd)

Non-deterministic [Ethereum 2 wallet](https://github.com/wealdtech/go-eth2-wallet).


## Table of Contents

- [Install](#install)
- [Usage](#usage)
- [Maintainers](#maintainers)
- [Contribute](#contribute)
- [License](#license)

## Install

`go-eth2-wallet-nd` is a standard Go module which can be installed with:

```sh
go get github.com/wealdtech/go-eth2-wallet-nd
```

## Usage

Access to the `wallet` is usually via [go-eth2-wallet](https://github.com/wealdtech/go-eth2-wallet); the first two examples below shows how this can be achieved.

This wallet generates keys non-deterministically, _i.e._ there is no relationship between keys or idea of a "seed".

Wallet and account names may be composed of any valid UTF-8 characters; the only restriction is they can not start with the underscore (`_`) character.

Note that although non-deterministic wallets do not have passphrases they still need to be unlocked before accounts can be created.  This can be carried out with `walllet.Unlock(nil)`

### Batches

This wallet provides the ability to create account batches.  A batch is a single piece of data that contains all accounts in a wallet at a given point in time, all encrypted with the same key.  This significantly decreases the time to obtain and decrypt accounts, however it does make the wallet less dynamic in that changes to accounts in the wallet will not be reflected in the batch automatically.

Batching is a manual process, and must be triggered by the user calling the `BatchWallet()` function.  It is recommended that batching is called once, after all required accounts in a wallet have been created.  It is possible to run subsequent `BatchWallet()` functions if further accounts have been added, however each call will recreate the batch in its entirety rather than incrementally on top of any existing batch, and as such it can take a significant amount of time to complete.  Wallets are unaware of changes in batches, so any `Wallet` would need to be discarded and re-opened after a call to `BatchWallet()`

### Example

#### Creating a wallet
```go
package main

import (
	e2wallet "github.com/wealdtech/go-eth2-wallet"
)

func main() {

    // Create a wallet
    wallet, err := e2wallet.CreateWallet("My wallet", e2wallet.WithType("non-deterministic"))
    if err != nil {
        panic(err)
    }

    ...
}
```

#### Accessing a wallet
```go
package main

import (
	e2wallet "github.com/wealdtech/go-eth2-wallet"
)

func main() {

    // Open a wallet
    wallet, err := e2wallet.OpenWallet("My wallet")
    if err != nil {
        panic(err)
    }

    ...
}
```

#### Creating an account
```go
package main

import (
	e2wallet "github.com/wealdtech/go-eth2-wallet"
)

func main() {

    // Open a wallet
    wallet, err := e2wallet.OpenWallet("My wallet")
    if err != nil {
        panic(err)
    }

    err = wallet.Unlock(nil)
    if err != nil {
        panic(err)
    }
    // Always immediately defer locking the wallet to ensure it does not remain unlocked outside of the function.
    defer wallet.Lock()
    
    account, err := wallet.CreateAccount("My account", []byte("my account secret"))
    if err != nil {
        panic(err)
    }
    // Wallet should be locked as soon as unlocked operations have finished; it is safe to explicitly call wallet.Lock() as well
    // as defer it as per above.
    wallet.Lock()

    ...
}
```

## Maintainers

Jim McDonald: [@mcdee](https://github.com/mcdee).

## Contribute

Contributions welcome. Please check out [the issues](https://github.com/wealdtech/go-eth2-wallet-nd/issues).

## License

[Apache-2.0](LICENSE) © 2019 Weald Technology Trading Ltd
