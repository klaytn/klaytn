package validator

import (
	"math/rand"
	"testing"
	"time"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/istanbul"
	"github.com/klaytn/klaytn/crypto"
	"github.com/stretchr/testify/assert"
)

func genAddress() common.Address {
	key, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(key.PublicKey)
	return addr
}

func genValidator() istanbul.Validator {
	rand.Seed(time.Now().UnixNano())
	testAddr := genAddress()
	testReward := genAddress()
	testVotingPower := rand.Uint64()
	testWeight := rand.Uint64()
	return newWeightedValidator(testAddr, testReward, testVotingPower, testWeight)
}

func makePoorValidator(val istanbul.Validator) istanbul.Validator {
	wval, _ := val.(*weightedValidator)
	wval.hasMinStaking = false
	return wval
}

func genValidatorSet(size int, poorValList ...int) []istanbul.Validator {
	var valSet []istanbul.Validator
	for i := 0; i < size; i++ {
		valSet = append(valSet, genValidator())
	}
	for _, i := range poorValList {
		valSet[i] = makePoorValidator(valSet[i])
	}
	return valSet
}

func isInValidatorSet(v istanbul.Validator, valSet []istanbul.Validator) bool {
	for _, member := range valSet {
		if v == member {
			return true
		}
	}
	return false
}

func TestSelectRandomCommittee(t *testing.T) {
	type testcase struct {
		name                  string
		validators            []istanbul.Validator
		committeeSize         uint64
		expectedCommitteeSize int
		seed                  int64
		proposerIdx           int
		nextProposerIdx       int
		isValid               bool
	}

	testcases := []testcase{
		{
			"valid-4-4",
			genValidatorSet(4),
			4, 4, 1, 1, 2, true,
		},
		{
			"valid-4-3",
			genValidatorSet(4),
			3, 3, 3, 2, 3, true,
		},
		{
			"valid-4-1",
			genValidatorSet(4),
			1, 1, 3, 2, 3, true,
		},
		{
			"valid-4-4-1poor",
			genValidatorSet(4, 3),
			4, 3, 3, 0, 1, true,
		},
		{
			"invalid-10-4-10poors",
			genValidatorSet(10, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9),
			4, 0, 4, 0, 1, false,
		},
		{ // when committeeSize is larger than validator size
			"invalid-4-5-larger-committee",
			genValidatorSet(4),
			5, 0, 5, 2, 3, false,
		},
		{ // when the proposer index is invalid
			"invalid-4-5-negative-proposerIdx",
			genValidatorSet(4),
			4, 0, 6, -1, 3, false,
		},
		{ // when the committee size is zero
			"invalid-4-5-negative-nextProposerIdx",
			genValidatorSet(4),
			0, 0, 7, 0, -1, false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			committee := SelectRandomCommittee(tc.validators, tc.committeeSize, tc.seed, tc.proposerIdx, tc.nextProposerIdx)
			if !tc.isValid {
				// to check invalid params which returns nil value
				assert.Nil(t, committee)
				return
			}
			assert.Equal(t, tc.expectedCommitteeSize, len(committee))
			assert.Equal(t, committee[0], tc.validators[tc.proposerIdx]) // check proposer
			if tc.committeeSize > 1 {                                    // if committee size is larger than 1, no next proposer exists
				assert.Equal(t, committee[1], tc.validators[tc.nextProposerIdx]) // check next proposer
			}
			for _, member := range committee {
				assert.True(t, true, member.HasMinStaking())
				if !isInValidatorSet(member, tc.validators) {
					t.Errorf("the committee member is wrong. expected: %v", member)
				}
			}
		})
	}
}

func TestCalcSeed(t *testing.T) {
	type testCase struct {
		hash        common.Hash
		expected    int64
		expectedErr error
	}
	testCases := []testCase{
		{
			hash:        common.HexToHash("0x1111"),
			expected:    0,
			expectedErr: nil,
		},
		{
			hash:        common.HexToHash("0x1234000000000000000000000000000000000000000000000000000000000000"),
			expected:    81979586966978560,
			expectedErr: nil,
		},
		{
			hash:        common.HexToHash("0x1234123412341230000000000000000000000000000000000000000000000000"),
			expected:    81980837895291171,
			expectedErr: nil,
		},
		{
			hash:        common.HexToHash("0x1234123412341234000000000000000000000000000000000000000000000000"),
			expected:    81980837895291171,
			expectedErr: nil,
		},
		{
			hash:        common.HexToHash("0x1234123412341234123412341234123412341234123412341234123412341234"),
			expected:    81980837895291171,
			expectedErr: nil,
		},
		{
			hash:        common.HexToHash("0x0034123412341234123412341234123412341234123412341234123412341234"),
			expected:    916044602622243,
			expectedErr: nil,
		},
		{
			hash:        common.HexToHash("0xabcdef3412341234123412341234123412341234123412341234123412341234"),
			expected:    773738372352131363,
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		actual, err := ConvertHashToSeed(tc.hash)
		assert.Equal(t, err, tc.expectedErr)
		assert.Equal(t, actual, tc.expected)
	}
}
