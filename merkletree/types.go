package merkletree

// MerkleTree represents a merkle tree
type MerkleTree struct {
	RootNode *MerkleNode
}

// MerkleNode represents a merkle tree node
type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}
