package merkletree

// make sure test files are named *_test.go
import (
	"encoding/hex"
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMerkle(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Merkle Suite")
}

var _ = Describe("Merkle Tree", func() {
	var (
		data [][]byte
	)

	BeforeEach(func() {
		data = [][]byte{
			[]byte("node1"),
			[]byte("node2"),
			[]byte("node3"),
		}
	})

	Describe("NewMerkleNode", func() {
		It("should create a Merkle node and calculate the correct hashes", func() {
			// Level 1
			n1 := NewMerkleNode(nil, nil, data[0])
			n2 := NewMerkleNode(nil, nil, data[1])
			n3 := NewMerkleNode(nil, nil, data[2])
			n4 := NewMerkleNode(nil, nil, data[2])

			// Level 2
			n5 := NewMerkleNode(n1, n2, nil)
			n6 := NewMerkleNode(n3, n4, nil)

			// Level 3
			n7 := NewMerkleNode(n5, n6, nil)

			Expect(hex.EncodeToString(n5.Data)).To(Equal("64b04b718d8b7c5b6fd17f7ec221945c034cfce3be4118da33244966150c4bd4"), "Level 1 hash 1 is correct")
			Expect(hex.EncodeToString(n6.Data)).To(Equal("08bd0d1426f87a78bfc2f0b13eccdf6f5b58dac6b37a7b9441c1a2fab415d76c"), "Level 1 hash 2 is correct")
			Expect(hex.EncodeToString(n7.Data)).To(Equal("4e3e44e55926330ab6c31892f980f8bfd1a6e910ff1ebc3f778211377f35227e"), "Root hash is correct")
		})
	})

	Describe("NewMerkleTree", func() {
		It("should create a Merkle tree and calculate the correct root hash", func() {
			// Level 1
			n1 := NewMerkleNode(nil, nil, data[0])
			n2 := NewMerkleNode(nil, nil, data[1])
			n3 := NewMerkleNode(nil, nil, data[2])
			n4 := NewMerkleNode(nil, nil, data[2])

			// Level 2
			n5 := NewMerkleNode(n1, n2, nil)
			n6 := NewMerkleNode(n3, n4, nil)

			// Level 3
			n7 := NewMerkleNode(n5, n6, nil)

			rootHash := fmt.Sprintf("%x", n7.Data)
			mTree := NewMerkleTree(data)

			Expect(fmt.Sprintf("%x", mTree.RootNode.Data)).To(Equal(rootHash), "Merkle tree root hash is correct")
		})
	})
})
