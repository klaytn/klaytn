package blst

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"io"
	"os"
	"path"
	"testing"

	"github.com/klaytn/klaytn/common"
	"github.com/stretchr/testify/assert"
)

// Use the file downloaded according to https://github.com/ethereum/bls12-381-tests
// TESTS_VERSION=v0.1.2
// wget https://github.com/ethereum/bls12-381-tests/releases/download/${TESTS_VERSION}/bls_tests_json.tar.gz

type testVectorChecker func(t *testing.T, name string, vectorJson []byte)

func TestVectors(t *testing.T) {
	checkers := map[string]testVectorChecker{
		"aggregate/":             checkAggregate,
		"batch_verify/":          checkBatchVerify,
		"deserialization_G1/":    checkDeserializationG1,
		"deserialization_G2/":    checkDeserializationG2,
		"fast_aggregate_verify/": checkFastAggregateVerify,
		"sign/":                  checkSign,
		"verify/":                checkVerify,
	}

	// If some test vectors cannot pass, list below with a reason.
	exemptions := map[string]bool{
		// PublicKeyFromBytes rejects inifinity point via KeyValidate()
		"deserialization_G1/deserialization_succeeds_infinity_with_true_b_flag.json": true,
		// AggregateVerify(asig, msgs, pks) not implemented
		"aggregate_verify/aggregate_verify_infinity_pubkey.json":                   true,
		"aggregate_verify/aggregate_verify_na_pubkeys_and_infinity_signature.json": true,
		"aggregate_verify/aggregate_verify_na_pubkeys_and_na_signature.json":       true,
		"aggregate_verify/aggregate_verify_tampered_signature.json":                true,
		"aggregate_verify/aggregate_verify_valid.json":                             true,
		// Low-level HashToG2 not implemented
		"hash_to_G2/hash_to_G2__2782afaa8406d038.json": true,
		"hash_to_G2/hash_to_G2__7590bd067999bbfb.json": true,
		"hash_to_G2/hash_to_G2__a54942c8e365f378.json": true,
		"hash_to_G2/hash_to_G2__c938b486cf69e8f7.json": true,
	}

	f, err := os.Open("./bls_tests_json.tar.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	stream, err := gzip.NewReader(f)
	reader := tar.NewReader(stream)

	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			t.Fatal(err)
		}

		switch header.Typeflag {
		case tar.TypeReg:
			name := path.Clean(header.Name)
			dir, _ := path.Split(name)
			if exemptions[name] {
				t.Logf("skip  %s", name)
				continue
			}

			checker := checkers[dir]
			if checker == nil {
				t.Fatalf("unrecognized vector file: %s", header.Name)
				// continue
			}

			data, err := io.ReadAll(reader)
			if err != nil {
				t.Fatal(err)
			}

			t.Logf("check %s", name)
			checker(t, name, data)
		}
	}
}

func checkAggregate(t *testing.T, name string, vectorJson []byte) {
	vector := struct {
		Input  []string
		Output *string
	}{}
	if err := json.Unmarshal(vectorJson, &vector); err != nil {
		t.Fatal(err)
	}

	sigbs := fromHexBatch(vector.Input)
	asig, err := AggregateSignaturesFromBytes(sigbs)
	if vector.Output != nil {
		asigb := common.FromHex(*vector.Output)
		assert.Nil(t, err)
		assert.Equal(t, asigb, asig.Marshal())
	} else {
		assert.NotNil(t, err)
	}
}

func checkBatchVerify(t *testing.T, name string, vectorJson []byte) {
	vector := struct {
		Input struct {
			Pubkeys    []string
			Messages   []string
			Signatures []string
		}
		Output bool
	}{}
	if err := json.Unmarshal(vectorJson, &vector); err != nil {
		t.Fatal(err)
	}

	var (
		pkbs      = fromHexBatch(vector.Input.Pubkeys)
		msgs      = fromHexBatch32(vector.Input.Messages)
		sigbs     = fromHexBatch(vector.Input.Signatures)
		pks, err1 = MultiplePublicKeysFromBytes(pkbs)
		ok, err2  = VerifyMultipleSignatures(sigbs, msgs, pks)
	)
	if vector.Output == true {
		assert.Nil(t, err1, name)
		assert.Nil(t, err2, name)
		assert.True(t, ok, name)
	} else {
		// At least an error or verify false.
		assert.True(t, (err1 != nil || err2 != nil || !ok), name)
	}
}

func checkFastAggregateVerify(t *testing.T, name string, vectorJson []byte) {
	vector := struct {
		Input struct {
			Pubkeys   []string
			Message   string
			Signature string
		}
		Output bool
	}{}
	if err := json.Unmarshal(vectorJson, &vector); err != nil {
		t.Fatal(err)
	}

	var (
		pkbs      = fromHexBatch(vector.Input.Pubkeys)
		msg       = common.FromHex(vector.Input.Message)
		sigb      = common.FromHex(vector.Input.Signature)
		apk, err1 = AggregatePublicKeysFromBytes(pkbs)
		sig, err2 = SignatureFromBytes(sigb)
		ok        = false
	)
	if err1 == nil && err2 == nil {
		ok = Verify(sig, msg, apk)
	}

	if vector.Output == true {
		assert.Nil(t, err1, name)
		assert.Nil(t, err2, name)
		assert.True(t, ok, name)
	} else {
		// At least an error or verify false.
		assert.True(t, (err1 != nil || err2 != nil || !ok), name)
	}
}

func checkDeserializationG1(t *testing.T, name string, vectorJson []byte) {
	vector := struct {
		Input  struct{ Pubkey string }
		Output bool
	}{}
	if err := json.Unmarshal(vectorJson, &vector); err != nil {
		t.Fatal(err)
	}

	b := common.FromHex(vector.Input.Pubkey)
	pk, err := PublicKeyFromBytes(b)
	if vector.Output == true {
		assert.Nil(t, err, name)
		assert.Equal(t, b, pk.Marshal(), name)
	} else {
		assert.NotNil(t, err, name)
	}
}

func checkDeserializationG2(t *testing.T, name string, vectorJson []byte) {
	vector := struct {
		Input  struct{ Signature string }
		Output bool
	}{}
	if err := json.Unmarshal(vectorJson, &vector); err != nil {
		t.Fatal(err)
	}

	b := common.FromHex(vector.Input.Signature)
	sig, err := SignatureFromBytes(b)
	if vector.Output == true {
		assert.Nil(t, err, name)
		assert.Equal(t, b, sig.Marshal(), name)
	} else {
		assert.NotNil(t, err, name)
	}
}

func checkSign(t *testing.T, name string, vectorJson []byte) {
	vector := struct {
		Input struct {
			Privkey string
			Message string
		}
		Output *string
	}{}
	if err := json.Unmarshal(vectorJson, &vector); err != nil {
		t.Fatal(err)
	}

	skb := common.FromHex(vector.Input.Privkey)
	msg := common.FromHex(vector.Input.Message)

	sk, err := SecretKeyFromBytes(skb)
	if vector.Output != nil {
		sigb := common.FromHex(*vector.Output)
		assert.Nil(t, err, name)
		assert.Equal(t, sigb, Sign(sk, msg).Marshal(), name)
	} else {
		assert.NotNil(t, err, name)
	}
}

func checkVerify(t *testing.T, name string, vectorJson []byte) {
	vector := struct {
		Input struct {
			Pubkey    string
			Message   string
			Signature string
		}
		Output bool
	}{}
	if err := json.Unmarshal(vectorJson, &vector); err != nil {
		t.Fatal(err)
	}

	var (
		pkb       = common.FromHex(vector.Input.Pubkey)
		msg       = common.FromHex(vector.Input.Message)
		sigb      = common.FromHex(vector.Input.Signature)
		pk, err1  = PublicKeyFromBytes(pkb)
		sig, err2 = SignatureFromBytes(sigb)
		ok        = false
	)
	if err1 == nil && err2 == nil {
		ok = Verify(sig, msg, pk)
	}

	if vector.Output == true {
		assert.Nil(t, err1, name)
		assert.Nil(t, err2, name)
		assert.True(t, ok, name)
	} else {
		// At least an error or verify false.
		assert.True(t, (err1 != nil || err2 != nil || !ok), name)
	}
}

func fromHexBatch(sArr []string) [][]byte {
	bArr := make([][]byte, len(sArr))
	for i, s := range sArr {
		bArr[i] = common.FromHex(s)
	}
	return bArr
}

func fromHexBatch32(sArr []string) [][32]byte {
	bArr := make([][32]byte, len(sArr))
	for i, s := range sArr {
		copy(bArr[i][:], common.FromHex(s))
	}
	return bArr
}
