package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/big"
	mrand "math/rand"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/crypto/blst"
	blstbind "github.com/supranational/blst/bindings/go"
)

// TODO: verify
// BLS12-381 group order
var q, _ = new(big.Int).SetString(
	"52435875175126190479447740508185965837690552500527637822603658699938581184513", 10,
)

var generalDST = []byte("BLS_SIG_BLS12381G2_XMD:SHA-256_SSWU_RO_POP_")

func randomNonZeroScalar() (*big.Int, error) {
	r := new(big.Int)
	var err error
	for r.Cmp(common.Big0) == 0 {
		r, err = randomScalar()
		if err != nil {
			return nil, err
		}
	}
	return r, nil
}

func randomScalar() (*big.Int, error) {
	r, err := rand.Int(rand.Reader, q)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// TODO: should I use keygen?
func bigToSecretKey(scalar *big.Int) blst.SecretKey {
	if scalar.BitLen() > blst.BlstScalarBytes*8 {
		panic("scalar is too big")
	}
	return blst.ScalarToSecretKey(bigToScalar(scalar))
}

func toLEBytes(n *big.Int) []byte {
	be := n.Bytes()

	le := make([]byte, len(be))
	for i := range be {
		le[i] = be[len(be)-1-i]
	}
	return le
}

func bigToScalar(b *big.Int) *blstbind.Scalar {
	if b.BitLen() > blst.BlstScalarBytes*8 {
		panic("scalar is too big")
	}
	return new(blstbind.Scalar).FromLEndian(
		common.RightPadBytes(toLEBytes(b), blst.BlstScalarBytes),
	)
}

func makePoly(t int, secret *big.Int) ([]*big.Int, error) {
	coeffs := make([]*big.Int, t)
	var err error

	coeffs[0] = new(big.Int).Set(secret)
	for i := 1; i < t; i++ {
		coeffs[i], err = randomNonZeroScalar()
		if err != nil {
			return nil, err
		}
	}
	return coeffs, nil
}

// return f(x) assuming coeff[0] x^0 ... coef
func evalPoly(coeffs []*big.Int, x *big.Int) *big.Int {
	res := big.NewInt(0)
	xAcc := big.NewInt(1)
	for _, a := range coeffs {
		term := new(big.Int).Mul(a, xAcc)
		term.Mod(term, q)
		res.Add(res, term)
		res.Mod(res, q)
		xAcc.Mul(xAcc, x)
		xAcc.Mod(xAcc, q)
	}
	return res
}

// compute Lagrange coefficient at x=0
func lagrangeCoeff(i *big.Int, Js []*big.Int) *big.Int {
	num := big.NewInt(1)
	den := big.NewInt(1)

	for _, j := range Js {
		if j.Cmp(i) == 0 {
			continue
		}
		// numerator
		tmp := new(big.Int).Neg(j)
		tmp.Mod(tmp, q)
		num.Mul(num, tmp)
		num.Mod(num, q)
		diff := new(big.Int).Sub(i, j)
		diff.Mod(diff, q)
		den.Mul(den, diff)
		den.Mod(den, q)
	}
	invDen := new(big.Int).ModInverse(den, q)
	lambda := new(big.Int).Mul(num, invDen)
	lambda.Mod(lambda, q)
	return lambda
}

func thresholdSigning() {
	n := 5 // number of participants
	t := 3 // threshold
	// TODO: from boldireya paper seems like there is a constraint t < n/2
	msg := []byte("Hello Threshold BLS")

	// generate master secret
	secretScalar, err := randomNonZeroScalar()
	if err != nil {
		panic(err)
	}
	fmt.Printf("secret scalar: %s (%v)\n", secretScalar.String(), secretScalar.Bytes())
	secretKey := bigToSecretKey(secretScalar)
	fmt.Printf("secret key %s (%v)\n", secretKey.Hex(), secretKey.Marshal())
	publicKey := secretKey.PublicKey()
	fmt.Printf("public key %s (%v)\n", publicKey.Hex(), publicKey.Marshal())

	// generate shamir poly
	poly, err := makePoly(t, secretScalar)
	if err != nil {
		panic(err)
	}
	fmt.Printf("poly %v \n", poly)

	// compute shares. NOTE: secretShares[0] == secretScalar
	secretShares := make([]*big.Int, 0, n)
	for i := 0; i <= n; i++ {
		secretShares = append(secretShares, evalPoly(poly, new(big.Int).SetUint64(uint64(i))))
	}
	fmt.Printf("secret shares %v \n", secretShares)

	// generate partial sigs
	partialSigs := make([]blst.Signature, 0, n) // NOTE: partialSigs[0] == reconstructed signature
	for i := 0; i <= n; i++ {
		secretShareKey := bigToSecretKey(secretShares[i])
		partialSigs = append(partialSigs, secretShareKey.Sign(msg))
	}
	fmt.Printf("%d partial signatures generated \n", len(partialSigs))

	// global signature should verify correctly
	if !partialSigs[0].Verify(publicKey, msg, false) {
		panic("global signature verification failed")
	} else {
		fmt.Println("global signature verification succeeded")
	}

	// reconstruct global signature using t partial sigs + lagrange coefficients
	participatingShares := t
	sigSubset := make([]blst.Signature, 0, participatingShares)
	xCoords := make([]*big.Int, 0, participatingShares)
	for i := 1; i <= participatingShares; i++ {
		sigSubset = append(sigSubset, partialSigs[i])
		xCoords = append(xCoords, new(big.Int).SetUint64(uint64(i)))
	}
	fmt.Printf("selected %d partial signatures as subset\n", len(sigSubset))
	fmt.Printf("x coordinates %v \n", xCoords)

	lagrangeCoeffs := make([]*big.Int, 0, participatingShares)
	lagrangeScalars := make([]blstbind.Scalar, 0, participatingShares)
	for _, i := range xCoords {
		coeff := lagrangeCoeff(i, xCoords)
		lagrangeCoeffs = append(lagrangeCoeffs, coeff)
		lagrangeScalars = append(lagrangeScalars, *bigToScalar(coeff))
	}
	fmt.Printf("lagrange coeffs %v \n", lagrangeCoeffs)
	fmt.Printf("lagrange scalars %v \n", lagrangeScalars)

	// verify all signature shares individually
	for i, sig := range sigSubset {
		publicKeyShare := bigToSecretKey(secretShares[xCoords[i].Int64()]).PublicKey()
		if !sig.Verify(publicKeyShare, msg, false) {
			panic("partial signature verification failed")
		}
	}

	reconstructedSig := blst.MultSigs(sigSubset, lagrangeScalars)

	fmt.Printf("reconstructed sig %s \n", reconstructedSig.Hex())

	// reconstructed sig should verify correctly
	if !reconstructedSig.Verify(publicKey, msg, false) {
		panic("reconstructed signature verification failed")
	} else {
		fmt.Println("reconstructed signature verification succeeded")
	}

	if bytes.Equal(reconstructedSig.Marshal(), partialSigs[0].Marshal()) {
		fmt.Println("reconstructed signature is == to global sig")
	} else {
		panic("reconstructed signature is != to global sig")
	}
}

// -----------------------------------------------------------------------------
// Threshold Encryption / Decryption (Baek–Zheng, PKC'03 Section II.C)
// -----------------------------------------------------------------------------

// returns a copy of the source to avoid modifying it by mistake with the Mult operations
// not sure if it is necessary or not, but keep for now
func copyPoint(source *blstbind.P1) *blstbind.P1 {
	sourceBytes := source.ToAffine().Serialize()
	sourceCopyAffine := new(blstbind.P1Affine)
	sourceCopyAffine.Deserialize(sourceBytes)
	sourceCopy := new(blstbind.P1)
	sourceCopy.FromAffine(sourceCopyAffine)
	return sourceCopy
}

func bitwiseXor(a, b []byte) []byte {
	if len(a) != len(b) {
		panic("mismatch bitwise xor")
	}
	res := make([]byte, len(a))
	for i, aa := range a {
		res[i] = aa ^ b[i]
	}
	return res
}

// taken from
// https://github.com/poanetwork/threshold_crypto/blob/master/src/lib.rs#L710
func xorWithHash(g1 *blstbind.P1, b []byte) ([]byte, error) {
	digest := sha256.Sum256(g1.ToAffine().Compress())
	seed := binary.BigEndian.Uint64(digest[:]) // TODO: fine that digest is 32 bytes instead of 8?
	// TODO: fine to use mrand?
	rng := mrand.New(mrand.NewSource(int64(seed))) // TODO: conversion uint64 --> int64

	rngBytes := make([]byte, len(b))
	n, err := rng.Read(rngBytes)
	if err != nil || n != len(b) {
		return nil, fmt.Errorf("error while generating random bytes. err: %v, n: %d, len(b): %d", err, n, len(b))
	}

	return bitwiseXor(rngBytes, b), nil
}

// https://github.com/poanetwork/threshold_crypto/blob/master/src/lib.rs#L697
func reduceUV(u *blstbind.P1, v []byte) []byte {
	uBytes := u.ToAffine().Compress()
	var vBytes []byte
	if len(v) > 64 {
		sum := sha256.Sum256(v)
		vBytes = sum[:]
	} else {
		vBytes = v
	}
	b := sha256.Sum256(append(vBytes, uBytes...))
	// TODO: the way I compute W deviates from PoA network.
	// 	     needs more verification
	return b[:]
}

func bzEncrypt(pk blst.PublicKey, msg []byte) (*blstbind.P1, []byte, *blstbind.P2, error) {
	r, err := randomNonZeroScalar()
	if err != nil {
		return nil, nil, nil, err
	}
	rScalar := bigToScalar(r)
	fmt.Printf("random scalar for encryption r: %s (%v)\n", r.String(), r.Bytes())

	// U = rP
	U := copyPoint(blstbind.P1Generator())
	U.MultAssign(rScalar)

	// V = F1(rY) xor M
	rY := copyPoint(pk.ToP1())
	rY.MultAssign(rScalar)
	V, err := xorWithHash(rY, msg)
	if err != nil {
		return nil, nil, nil, err
	}

	// W = r * F2(U,V)
	// TODO: this part needs more cryptographic verification
	W := blstbind.HashToG2(reduceUV(U, V), generalDST)
	W.MultAssign(rScalar)

	// sanity check, verification of the tag should work after generating it
	if errTag := verifyTag(U, V, W); errTag != nil {
		return nil, nil, nil, fmt.Errorf("should never happen %w", errTag)
	}

	return U, V, W, nil
}

func verifyTag(U *blstbind.P1, V []byte, W *blstbind.P2) error {
	WAffine := W.ToAffine()
	if !WAffine.SigValidate(true) { // TODO: inf check or not
		return fmt.Errorf("W signature fails group check")
	}
	if !U.ToAffine().KeyValidate() {
		return fmt.Errorf("U public key fails group check")
	}
	// e(P,W) == e(U,H(U,V))
	// where W = rH(U,V)
	//       U = rP
	if !WAffine.Verify(false, U.ToAffine(), false, reduceUV(U, V), generalDST) {
		return fmt.Errorf("W signature fails verification")
	}
	return nil
}

func bzDecrypt(U *blstbind.P1, V []byte, W *blstbind.P2, sk *blstbind.Scalar) ([]byte, error) {
	// step 1. verify sig
	if err := verifyTag(U, V, W); err != nil {
		return nil, err
	}

	// decrypt ciphertext
	// M = F1(xU) xor V
	xU := copyPoint(U)
	xU.MultAssign(sk)

	m, err := xorWithHash(xU, V)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func bzPartialDecrypt(U *blstbind.P1, V []byte, W *blstbind.P2, sk *blstbind.Scalar) (*blstbind.P1, error) {
	// step 1. verify sig
	if err := verifyTag(U, V, W); err != nil {
		return nil, err
	}

	// compute decryption share
	xU := copyPoint(U)
	xU.MultAssign(sk)

	return xU, nil
}

// TODO: verifyShare is a bit different from the crypto algo of the paper because I just panic in case of invalid share
//		in production this will need to be handled as it can of course happen

func verifyShare(decryptionShare *blstbind.P1, publicKeyShare *blstbind.P1, U *blstbind.P1, V []byte, W *blstbind.P2) error {
	if err := verifyTag(U, V, W); err != nil {
		return err
	}

	// TODO: it is highly likely that the following verification might not be correct.

	// paper says to verify
	// e(P, Ui) == e(U, Yi)
	// however that is not possible because of the groups the points belong to.
	// P --> I can set it to g1 or g2
	// Ui --> g1 (forced because it is a scalar multiplication of U)
	// U --> g1
	// Yi --> g1 I guess but I could map the secret share to g2, so not really sure

	// https://github.com/poanetwork/threshold_crypto/blob/master/src/lib.rs#L182-L186 does:
	// e(Ui, H(U,V)) == e(Yi, W)
	// which seems to make sense because it is equivalent to:
	// e(xi * r * P, H(U,V)) = e(xi * P, r * H(U,V))
	// same thing seems to be done here: https://github.com/LATOKEN/lachain/blob/dev/src/Lachain.Crypto/TPKE/PublicKey.cs#L90-L94

	//pairing := blstbind.PairingCtx(false, generalDST)
	//blstbind.PairingRawAggregate(pairing, blstbind.HashToG2(reduceUV(U, V), generalDST).ToAffine(), decryptionShare.ToAffine())

	if !pairing(decryptionShare, blstbind.HashToG2(reduceUV(U, V), generalDST), publicKeyShare, W) {
		return fmt.Errorf("decryption share fails verification")
	}
	return nil
}

// TODO: the pairing function passes basic smoke-checks, but needs further scrutiny
// e(P,Q) == e(R,S)
func pairing(p *blstbind.P1, q *blstbind.P2, r *blstbind.P1, s *blstbind.P2) bool {
	pairing1 := blstbind.PairingCtx(false, generalDST)
	blstbind.PairingRawAggregate(pairing1, q.ToAffine(), p.ToAffine())
	blstbind.PairingCommit(pairing1)
	pairing2 := blstbind.PairingCtx(false, generalDST)
	blstbind.PairingRawAggregate(pairing2, s.ToAffine(), r.ToAffine())
	blstbind.PairingCommit(pairing2)

	return blstbind.Fp12FinalVerify(blstbind.PairingAsFp12(pairing1), blstbind.PairingAsFp12(pairing2))
}

// recombine decryption shares and decrypt message
func bzSharesDecrypt(U *blstbind.P1, V []byte, W *blstbind.P2, decryptionShares []*blstbind.P1) ([]byte, error) {
	// verify tag first
	if err := verifyTag(U, V, W); err != nil {
		return nil, err
	}

	// assumes indexes are 1,2,...,t
	t := len(decryptionShares)
	Js := make([]*big.Int, 0, t)
	for i := 1; i <= t; i++ {
		Js = append(Js, big.NewInt(int64(i)))
	}
	lambdas := make([]*blstbind.Scalar, 0, t)
	decryptionSharesAffine := make(blstbind.P1Affines, 0, t)
	for i := 1; i <= t; i++ {
		lambda := bigToScalar(lagrangeCoeff(big.NewInt(int64(i)), Js))
		lambdas = append(lambdas, lambda)
		decryptionSharesAffine = append(decryptionSharesAffine, *decryptionShares[i-1].ToAffine())
	}
	reconstructedKey := decryptionSharesAffine.Mult(lambdas, 255) // TODO: nbits

	return xorWithHash(reconstructedKey, V)
}

func thresholdDecryption() {
	fmt.Println("---- Threshold Encryption / Decryption (Baek–Zheng) ----")

	n := 5 // number of participants
	t := 3 // threshold
	msg := []byte("Hello Threshold BLS decryption")

	// master secret
	secretScalar, err := randomNonZeroScalar()
	if err != nil {
		panic(err)
	}
	secretKey := bigToSecretKey(secretScalar)
	publicKey := secretKey.PublicKey()

	fmt.Printf("secret scalar: %s (%v)\n", secretScalar.String(), secretScalar.Bytes())
	fmt.Printf("secret key %s (%v)\n", secretKey.Hex(), secretKey.Marshal())
	fmt.Printf("public key %s (%v)\n", publicKey.Hex(), publicKey.Marshal())

	// Shamir polynomial
	poly, err := makePoly(t, secretScalar)
	if err != nil {
		panic(err)
	}

	// compute shares. NOTE: secretShares[0] == secretScalar
	secretShares := make([]*big.Int, 0, n)
	for i := 0; i <= n; i++ {
		secretShares = append(secretShares, evalPoly(poly, new(big.Int).SetUint64(uint64(i))))
	}
	fmt.Printf("secret shares %v \n", secretShares)

	// encrypt
	u, v, w, err := bzEncrypt(publicKey, msg)
	if err != nil {
		panic(err)
	}

	fmt.Println("ciphertext encrypted")

	// master secret key should be able to decrypt
	decryptedMsg, err := bzDecrypt(u, v, w, bigToScalar(secretScalar))
	if err != nil {
		panic(err)
	}
	fmt.Printf("decryptedMsg by master secret key: %v\n", string(decryptedMsg))
	if !bytes.Equal(decryptedMsg, msg) {
		panic("decryption failed")
	}

	// pick t participants to decrypt
	// NOTE: share 0 == full decryption
	decryptionShares := make([]*blstbind.P1, 0, t)
	for i := 0; i <= t; i++ {
		decryptionShare, err := bzPartialDecrypt(u, v, w, bigToScalar(secretShares[i]))
		if err != nil {
			panic(err)
		}
		decryptionShares = append(decryptionShares, decryptionShare)
	}

	for i, share := range decryptionShares {
		if err := verifyShare(share, bigToSecretKey(secretShares[i]).PublicKey().ToP1(), u, v, w); err != nil {
			panic(err)
		}
	}

	// decrypt by recombining shares
	decryptedMsg, err = bzSharesDecrypt(u, v, w, decryptionShares[1:])
	if err != nil {
		panic(err)
	}
	fmt.Printf("decryptedMsg by share recombination: %v\n", string(decryptedMsg))
	if !bytes.Equal(decryptedMsg, msg) {
		panic("decryption by share recombination failed")
	}
}

func main() {
	thresholdSigning()
	thresholdDecryption()
}
