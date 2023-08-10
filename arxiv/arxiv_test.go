package arxiv

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"arxiv_server/utils"

	"github.com/stretchr/testify/assert"
)

func newArxivManager() ArxivMetadataManager {
	return ArxivMetadataManager{Root_path: utils.GetEnv("DUMP_PATH", "E:/data/arxiv_dump")}
}

func TestArxivMetadataSearch(t *testing.T) {
	mgr := newArxivManager()
	res, n, err := mgr.SearchArxivMetadata(1e4,
		func(elm *ArxivMetadata) bool { return strings.Contains(elm.Authors, "Louis") })
	if err != nil {
		t.Error(err)
	}
	assert.Greater(t, len(res), 0)
	fmt.Printf("%v elements found amongst %v\n", len(res), n)
	// for _, elm := range res {
	// 	fmt.Println(elm.Authors)
	// }
}

func BenchmarkReadAllArxivMetadata(b *testing.B) {
	mgr := newArxivManager()
	_, n, err := mgr.SearchArxivMetadata(-1, nil)
	if err != nil {
		b.Error(err)
	}
	fmt.Printf("%v elements found\n", n)
}

func BenchmarkGenerateArxivMetadataIndex(b *testing.B) {
	mgr := newArxivManager()
	err := mgr.generateArxivMetadataIndex()
	if err != nil {
		b.Error(err)
	}
}

func TestGetIndexedArxivMetadata(t *testing.T) {
	mgr := newArxivManager()
	err := mgr.InitializeManager()
	if err != nil {
		t.Error(err)
	}
	elm, err := mgr.GetIndexedArxivMetadata("0704.0426")
	if err != nil {
		t.Error(err)
	}
	val, err := json.Marshal(elm)
	if err != nil {
		t.Error(err)
	}
	assert.Contains(t, elm.Authors, "Kuhlen")
	fmt.Printf("%v\n", string(val))
}

func BenchmarkGetIndexedArxivMetadata(b *testing.B) {
	mgr := newArxivManager()
	err := mgr.InitializeManager()
	if err != nil {
		b.Error(err)
	}

	randIds, err := mgr.getRandomArxivIds(100000)
	if err != nil {
		b.Error(err)
	}
	b.Logf("Generated %d random ids\n", len(*randIds))

	for _, id := range *randIds {
		elm, err := mgr.GetIndexedArxivMetadata(id)
		if err != nil {
			b.Error(err)
		}
		assert.Equal(b, elm.Id, id)
		//b.Logf("%s\n", elm.Id)
	}
}

func TestGetRandomArxivIds(t *testing.T) {
	mgr := newArxivManager()
	err := mgr.InitializeManager()
	if err != nil {
		t.Error(err)
	}
	ids, err := mgr.getRandomArxivIds(1000)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%v\n", ids)
}
