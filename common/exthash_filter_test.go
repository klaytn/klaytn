// Modifications Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from common/types_test.go (2018/06/04).
// Modified and improved for the klaytn development.

package common

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestFiltersCase1(t *testing.T) {
	srcStr, answer := []byte("f8729e20fb89cf6444214fbeec19fa56fb55644b3d95cd60b2db025a0a9a3d0bcfb85101f84e02808005f848a302a10241648c5a079203f35708ba25630dee558dcffd6cc0a03d2d8903c595455154b8a302a102748fe1aa2f0670f02a82934c2e0fbe691455a845fcdbf053ffa14b38a8f93a9f"), []byte("f8729e20fb89cf6444214fbeec19fa56fb55644b3d95cd60b2db025a0a9a3d0bcfb85101f84e02808005f848a302a10241648c5a079203f35708ba25630dee558dcffd6cc0a03d2d8903c595455154b8a302a102748fe1aa2f0670f02a82934c2e0fbe691455a845fcdbf053ffa14b38a8f93a9f")

	dst := make([]byte, hex.DecodedLen(len(srcStr)))
	n, _ := hex.Decode(dst, srcStr)

	reAns := make([]byte, hex.DecodedLen(len(answer)))
	n2, _ := hex.Decode(reAns, answer)

	reStr, _ := RlpPaddingFilter(dst[:n])

	if !bytes.Equal(reAns[:n2], reStr) {
		t.Errorf("Expected %x got %x", reAns[:n2], reStr)
	}
}

func TestFiltersCase3(t *testing.T) {
	srcStr, answer := []byte("f8729e3aa21f22bb1d06ffd3052056f424f5a75c897909bf19054f5c5faef94ba4b85101f84e02808005f848a302a10284d9016c570bd864eaa2621c482acf9a8109a8bbd315369029b6433a6d591a6fa302a102748fe1aa2f0670f02a82934c2e0fbe691455a845fcdbf053ffa14b38a8f93a9f"), []byte("f8729e3aa21f22bb1d06ffd3052056f424f5a75c897909bf19054f5c5faef94ba4b85101f84e02808005f848a302a10284d9016c570bd864eaa2621c482acf9a8109a8bbd315369029b6433a6d591a6fa302a102748fe1aa2f0670f02a82934c2e0fbe691455a845fcdbf053ffa14b38a8f93a9f")
	dst := make([]byte, hex.DecodedLen(len(srcStr)))
	n, _ := hex.Decode(dst, srcStr)

	reAns := make([]byte, hex.DecodedLen(len(answer)))
	n2, _ := hex.Decode(reAns, answer)

	reStr, _ := RlpPaddingFilter(dst[:n])

	if !bytes.Equal(reAns[:n2], reStr) {
		t.Errorf("Expected %x got %x", reAns[:n2], reStr)
	}
}

func TestFiltersCase4(t *testing.T) {
	srcStr, answer := []byte("f8729e379923746b521cba3a25830d3c5c7c4eca96f39c4eb50d055302c6adc5bbb85101f84e02808005f848a302a1021bf80c2bb3766fe0ad90dd3838a6d37fc4e8d99ee5ca2756daf0c8a1a02ea7c7a302a102748fe1aa2f0670f02a82934c2e0fbe691455a845fcdbf053ffa14b38a8f93a9f"), []byte("f8729e379923746b521cba3a25830d3c5c7c4eca96f39c4eb50d055302c6adc5bbb85101f84e02808005f848a302a1021bf80c2bb3766fe0ad90dd3838a6d37fc4e8d99ee5ca2756daf0c8a1a02ea7c7a302a102748fe1aa2f0670f02a82934c2e0fbe691455a845fcdbf053ffa14b38a8f93a9f")
	dst := make([]byte, hex.DecodedLen(len(srcStr)))
	n, _ := hex.Decode(dst, srcStr)

	reAns := make([]byte, hex.DecodedLen(len(answer)))
	n2, _ := hex.Decode(reAns, answer)

	reStr, _ := RlpPaddingFilter(dst[:n])

	if !bytes.Equal(reAns[:n2], reStr) {
		t.Errorf("Expected %x got %x", reAns[:n2], reStr)
	}
}

func TestFiltersCase5(t *testing.T) {
	srcStr, answer := []byte("f8729e207ed723a0db1f98bfdf6ad39825aad0f8411ec118542883c9ca7b2b9f41b85101f84e02808005f848a302a103c73a7767c4294627455a4ad454fd5b0ec5cf4dbf9db0c1bd4ae3895ae754e3a8a302a102748fe1aa2f0670f02a82934c2e0fbe691455a845fcdbf053ffa14b38a8f93a9f"), []byte("f8729e207ed723a0db1f98bfdf6ad39825aad0f8411ec118542883c9ca7b2b9f41b85101f84e02808005f848a302a103c73a7767c4294627455a4ad454fd5b0ec5cf4dbf9db0c1bd4ae3895ae754e3a8a302a102748fe1aa2f0670f02a82934c2e0fbe691455a845fcdbf053ffa14b38a8f93a9f")
	dst := make([]byte, hex.DecodedLen(len(srcStr)))
	n, _ := hex.Decode(dst, srcStr)

	reAns := make([]byte, hex.DecodedLen(len(answer)))
	n2, _ := hex.Decode(reAns, answer)

	reStr, _ := RlpPaddingFilter(dst[:n])

	if !bytes.Equal(reAns[:n2], reStr) {
		t.Errorf("Expected %x got %x", reAns[:n2], reStr)
	}
}

func TestFiltersCase6(t *testing.T) {
	srcStr, answer := []byte("03c73a7767c4294627455a4ad454fd5b0ec5cf4dbf9db0c1bd4ae3895ae754e3a8"), []byte("03c73a7767c4294627455a4ad454fd5b0ec5cf4dbf9db0c1bd4ae3895ae754e3a8")
	dst := make([]byte, hex.DecodedLen(len(srcStr)))
	n, _ := hex.Decode(dst, srcStr)

	reAns := make([]byte, hex.DecodedLen(len(answer)))
	n2, _ := hex.Decode(reAns, answer)

	reStr, _ := RlpPaddingFilter(dst[:n])

	if !bytes.Equal(reAns[:n2], reStr) {
		t.Errorf("Expected %x got %x", reAns[:n2], reStr)
	}
}

func TestFiltersCase8(t *testing.T) {
	srcStr, answer := []byte("f8999e20cc9ef35b62b1335a9af34451a1c57fc07256f4df842497856bd6d808d9b87801f87501808004f86f02f86ce301a103d65fabd61151c76536b7a93f4c2731fd817f6c30f54e0d03b80631733847d121e301a103ef5bb9400ef6dfd55ec1924b3f9302064d132513af1852287aef7cfe850eb522e301a1024f450f9d909824f3a5726d586234398851d325006f0500bcc2f78dfbd92ccb72"), []byte("f8999e20cc9ef35b62b1335a9af34451a1c57fc07256f4df842497856bd6d808d9b87801f87501808004f86f02f86ce301a103d65fabd61151c76536b7a93f4c2731fd817f6c30f54e0d03b80631733847d121e301a103ef5bb9400ef6dfd55ec1924b3f9302064d132513af1852287aef7cfe850eb522e301a1024f450f9d909824f3a5726d586234398851d325006f0500bcc2f78dfbd92ccb72")
	dst := make([]byte, hex.DecodedLen(len(srcStr)))
	n, _ := hex.Decode(dst, srcStr)

	reAns := make([]byte, hex.DecodedLen(len(answer)))
	n2, _ := hex.Decode(reAns, answer)

	reStr, _ := RlpPaddingFilter(dst[:n])

	if !bytes.Equal(reAns[:n2], reStr) {
		t.Errorf("Expected %x got %x", reAns[:n2], reStr)
	}
}

func TestFiltersCase9(t *testing.T) {
	srcStr, answer := []byte("03d65fabd61151c76536b7a93f4c2731fd817f6c30f54e0d03b80631733847d121"), []byte("03d65fabd61151c76536b7a93f4c2731fd817f6c30f54e0d03b80631733847d121")
	dst := make([]byte, hex.DecodedLen(len(srcStr)))
	n, _ := hex.Decode(dst, srcStr)

	reAns := make([]byte, hex.DecodedLen(len(answer)))
	n2, _ := hex.Decode(reAns, answer)

	reStr, _ := RlpPaddingFilter(dst[:n])

	if !bytes.Equal(reAns[:n2], reStr) {
		t.Errorf("Expected %x got %x", reAns[:n2], reStr)
	}
}
