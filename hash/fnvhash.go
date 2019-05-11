package fnvhash

type FNVHash struct {
}

func (hashH *FNVHash) GenerateHash(num string) uint32 {
	var hash uint32 = FnvOffset32

	bytes := []byte(num)

	for _, val := range bytes {
		hash ^= uint32(val)
		hash *= FnvPrime32
	}

	return hash
}

const (
	FnvPrime32  uint32 = 16777619
	FnvOffset32 uint32 = 2166136261
)
