package blst

import bind "github.com/supranational/blst/bindings/go"

// Internal types for blst.
type blstPublicKey = bind.P1Affine
type blstSignature = bind.P2Affine
type blstAggregatePublicKey = bind.P1Aggregate
type blstAggregateSignature = bind.P2Aggregate
