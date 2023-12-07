package blst

import bind "github.com/supranational/blst/bindings/go"

// Internal types for blst. They determine which group is used for public keys and which for signatures.
// Currently we use the small-and-fast group for public keys and the slow-and-big group for signatures.
// for more context checkout: https://medium.com/nethermind-eth/bls-signatures-withdrawals-bbf38658c242#eea8
type blstPublicKey = bind.P1Affine
type blstSignature = bind.P2Affine
type blstAggregatePublicKey = bind.P1Aggregate
type blstAggregateSignature = bind.P2Aggregate
