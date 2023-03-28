# Release

Описание релизов foundation (REALEASE NOTES)

## TOC

- [Release](#release)
  - [TOC](#toc)
  - [v3.1.1](#v311)
  - [v3.1.0](#v310)
  - [v3.0.0](#v300)
  - [v2.0.0](#v200)
  - [v1.1.0](#v110)
  - [v1.0.0](#v100)
  - [Links](#links)

## v3.1.1

_Changes:_

* Methods `AddAccountRight()` and `RemoveAccountRight()` in packege `acl` moved to package `mock`, as methods in Wallet type. - https://github.com/tickets-dao/foundation/-/issues/73
* Dependency `fabtic-chaincode-go` updated to last version, new method of ChaincodeStubInterface `PurgePrivateData()` added to `mock/stub` - https://github.com/tickets-dao/foundation/-/issues/73

_Notes:_

There is fully backward compatibility with v3.0.0


## v3.1.0

_Changes:_

Added package `acl` with ACL access matrix wrappers, https://github.com/tickets-dao/foundation/-/issues/68:  
  * `AddAccountRight()`
  * `GetAccountRight()`
  * `GetAccountAllRights()`
  * `RemoveAccountRight()`
  * `GetAccountAllRights()`

More information in _core/acl/accessmatrix.go_.

_Notes:_

There is fully backward compatibility with v3.0.0

## v3.0.0

_Changes:_
- `golang_math_big` replaced with `core/types/big`, https://github.com/tickets-dao/foundation/-/issues/58
- `*types.Sender` and `*types.Address` can be passed as pointers to avoid copying locks.
  , https://github.com/tickets-dao/foundation/-/issues/3
- fixed linter warnings, https://github.com/tickets-dao/foundation/-/issues/11

_Notes:_

There is **NO** backward compatibility with v2.0.0, please be attentive with updating your chaincodes.

If you want to use big numbers in chaincodes use `core/types/big`.  
If you want to use `types.Sender` or `types.Address` in chaincodes pass them via pointer `*types.Sender` or `*types.Address` to avoid copying locks. 

## v2.0.0

_Changes:_

- Added TTL for transactions in chaincode option `TxTTL`, https://github.com/tickets-dao/foundation/-/issues/51
- Added several transactions per user in same time in chaincode option `NonceTTL`, https://github.com/tickets-dao/foundation/-/issues/53
- Package `golang-math-big` changed name to `golang_math_big`, https://github.com/tickets-dao/foundation/-/issues/3
- Changed Foundation API methods names with code name convention, https://github.com/tickets-dao/foundation/-/issues/3:
    * func (bc *BaseContract) GetAllowedMspID() string
    * func (a Address) IsUserIDSame(b Address) bool
- Changed `MockStub` name to `Stub`, https://github.com/tickets-dao/foundation/-/issues/3

_Notes:_

There is **NO** backward compatibility with v1.0.0, please be attentive with updating your chaincodes.

Changed storing format for pending batch transactions to `protobuf` with type `pb.PendingTx`.  
Also, changed nonce storing, now it stored as `protobuf` with `pb.Nonce` type.
Be sure, that you use `observer` & `robot` which support this changes. Check product dependency to foundation v2.x.x.

## v1.1.0

_Changes:_

Added package `acl` with ACL access matrix wrappers, https://github.com/tickets-dao/foundation/-/issues/68:  
  * `AddAccountRight()`
  * `GetAccountRight()`
  * `GetAccountAllRights()`
  * `RemoveAccountRight()`
  * `GetAccountAllRights()`

More information in _core/acl/accessmatrix.go_.

_Notes:_

There is fully backward compatibility with v1.0.0

## v1.0.0

_Changes:_

  - Moved `foundation` codebase from `atmz` to `core` repository, https://github.com/tickets-dao/foundation/-/issues/2
  - Code cleaned, https://github.com/tickets-dao/foundation/-/issues/3
  - Added gitignore, https://github.com/tickets-dao/foundation/-/issues/4
  - Added unit tests, https://github.com/tickets-dao/foundation/-/issues/5
  - Add pipeline for foundation repo, https://github.com/tickets-dao/foundation/-/issues/7
  - Fixed codebase after adding pipeline, fixed linters, https://github.com/tickets-dao/foundation/-/issues/9
  - Written acceptance tests checking the basic functionality only, https://github.com/tickets-dao/foundation/-/issues/10
    - added tests for Given, Token, Allowed balances
    - added tests for TxSetRate and TxSetLimits APIs
  - Integration test, https://github.com/tickets-dao/foundation/-/issues/43
    - created pipeline for integration tests with sandbox
    - added integration tests: add user, metadata, swap, swap back, transfer, multiswap, multiswap back
  - Add sonar properties, https://github.com/tickets-dao/foundation/-/issues/36
  - Removed etc, industrial and fractionalized tokens, https://github.com/tickets-dao/foundation/-/issues/35

_Notes:_

This release fully backward-compatible with release v0.8.2 

## Links

* no