package merkle

import (
	"math/big"
	"testing"
)

func TestMerkleTree(t *testing.T) {
	t.Run("TestMerkleThreeRoot", func(t *testing.T) {
		mt, err := initMerkleTree(8, 0)
		if err != nil {
			t.Fatalf("initMerkleTree: %v", err)
		}
		comm1, _ := big.NewInt(0).
			SetString("18770007740858590904834323875249445187602972996157638065511776196192874703345", 10)
		comm2, _ := big.NewInt(0).
			SetString("8220608021600102054607175486139306658950270982542023113115274387441557849105", 10)

		mt.Tree[0] = []*big.Int{comm1, comm2}
		// Recompute intermediate levels.
		if err := mt.RebuildSparseTree(); err != nil {
			t.Fatalf("RebuildSparseTree: %v", err)
		}
		root := mt.Root()
		want, _ := big.NewInt(0).
			SetString("1847418543371880308816668005262896997760433961565326404501878257505837164890", 10)
		if root.Cmp(want) != 0 {
			t.Errorf("Root is not correct, got %s, want %s", root.String(), want.String())
		}
	})

	t.Run("TestMerkleThreeRoot_2", func(t *testing.T) {
		mt, err := initMerkleTree(8, 0)
		if err != nil {
			t.Fatalf("initMerkleTree: %v", err)
		}
		comm1, _ := big.NewInt(0).
			SetString("6069980975591022844080630325257076595188855687539168746753896913806183535207", 10)

		mt.Tree[0] = []*big.Int{comm1}
		// Recompute intermediate levels.
		if err := mt.RebuildSparseTree(); err != nil {
			t.Fatalf("RebuildSparseTree: %v", err)
		}
		root := mt.Root()
		want, _ := big.NewInt(0).
			SetString("3401459129837206331229789859975562117587875955762767199191635934962960578494", 10)
		if root.Cmp(want) != 0 {
			t.Errorf("Root is not correct, got %s, want %s", root.String(), want.String())
		}
	})
}
