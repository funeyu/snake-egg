package indexer

import (
	"fmt"
	"testing"
)

func TestGeneDocId2(t *testing.T) {
	id, _ := GeneDocId(12, 200, 899)
	fmt.Println("id", id)

	fmt.Println("rank_id", id.RankId())
	fmt.Println("sub_rank", id.SubRankId())
	fmt.Println("i", id.Index())
}
