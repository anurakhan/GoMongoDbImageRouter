package hashring

import (
	"fmt"

	fnvhash "github.com/anurakhan/go-mongo-lb-driver/hash"
	"github.com/anurakhan/go-mongo-lb-driver/server"
)

type IConsistentHashRing interface {
	AddServer(server *server.Server)
	GetServerForKey(key string) *server.Server
	InitRing()
}

type ConsistentHashRing struct {
	ring map[uint32]*server.Server
}

func (consistentHashRing *ConsistentHashRing) InitRing() {
	consistentHashRing.ring = make(map[uint32]*server.Server)
}

func (consistentHashRing *ConsistentHashRing) AddServer(server *server.Server) {
	hashFunc := new(fnvhash.FNVHash)
	var hash uint32 = hashFunc.GenerateHash(string(server.KeyForCh))
	fmt.Println("Added server| hash: " + fmt.Sprint(hash) + " | address: " + string(server.KeyForCh))
	consistentHashRing.ring[hash] = server
}

func (consistentHashRing *ConsistentHashRing) GetServerForKey(key string) *server.Server {
	if len(consistentHashRing.ring) == 0 {
		return nil
	}

	hashFunc := new(fnvhash.FNVHash)
	hash := hashFunc.GenerateHash(key)
	if val, ok := consistentHashRing.ring[hash]; ok {
		return val
	}

	var minServer *server.Server
	var resServer *server.Server
	var dif uint32 = ^uint32(0)
	var min uint32 = ^uint32(0)

	for key, value := range consistentHashRing.ring {
		if key > hash && key-hash < dif {
			resServer = value
			dif = key - hash
		}

		if key < min {
			min = key
			minServer = value
		}
	}

	//when all keys are smaller take first on the wheel, which is the min.
	if resServer == nil {
		resServer = minServer
	}

	fmt.Println("key from request| hash: " + fmt.Sprint(hash) + " | key: " + key + " | instance: " + resServer.Name)
	return resServer
}
