package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/crypto/blst"
	blstbind "github.com/supranational/blst/bindings/go"
)

// TODO: verify
// BLS12-381 group order
var q, _ = new(big.Int).SetString(
	"52435875175126190479447740508185965837690552500527637822603658699938581184513", 10,
)

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

// TODO: correct to compute in x=0?
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

func main() {
	thresholdSigning()
}
