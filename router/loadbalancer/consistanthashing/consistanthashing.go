// NOTE: will not be used in research analysis due to nature of bench simulation 
package consistanthashing

import (
	"hash/fnv"
	"llm-routing-bench/router/backend"
	"log"
	"net/http"
	"slices"
	"sort"
)

type CHashing struct {
	sortedKeys []uint32
	Ring       map[uint32]*backend.Backend
}

func NewConsistantHash(backends []backend.Backend) *CHashing {
	cHashInstance := CHashing{
		Ring: make(map[uint32]*backend.Backend),
	}
	for i := range backends {
		hashVal := cHashInstance.addNode(&backends[i])
		cHashInstance.sortedKeys = append(cHashInstance.sortedKeys, hashVal)
	}
	slices.Sort(cHashInstance.sortedKeys)
	return &cHashInstance
}

func (ch *CHashing) Route(r *http.Request) *backend.Backend {
	if len(ch.sortedKeys) == 0 {
		return nil
	}

	key := r.URL.String()
	nodeHash := ch.getNode(key)

	return ch.Ring[nodeHash]
}

func (ch *CHashing) hashFunc(str string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(str))
	return h.Sum32()
}

func (ch *CHashing) addNode(b *backend.Backend) uint32 {
	hash := ch.hashFunc(b.BackendURI)

	if _, exists := ch.Ring[hash]; exists {
		log.Fatalln("Collision Detected")
	}

	ch.Ring[hash] = b
	return hash
}

func (ch *CHashing) getNode(str string) uint32 {
	hash := ch.hashFunc(str)

	i := sort.Search(len(ch.sortedKeys), func(i int) bool {
		return ch.sortedKeys[i] >= hash
	})

	if i == len(ch.sortedKeys) {
		return ch.sortedKeys[0]
	}

	return ch.sortedKeys[i]
}
