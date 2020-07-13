package indexer

/**
时间戳索引: 包含pair: <docId, timeStamp>;根据timestamp 倒排
*/

type TimestampPair struct {
	id DocId
	timestamp uint32
}

func (t TimestampPair) Compare(other interface{}) int {
	// t.timestamp - other.(TimestampPair).timestamp 注意点 uint32 相减是得不出负数的
	return int(t.timestamp) - int(other.(TimestampPair).timestamp)
}

func (t TimestampPair) Serialize() []byte {
	byteSize := t.ByteSize()
	bytes := make([]byte, byteSize)

	return bytes
}

func (t TimestampPair) ByteSize() uint8 {
	return 8 + 4
}

// 注意这里的返回值不能写成 TimestampPair，必须要和 ElementCreate 一致
func TimePairCreate(bytes []byte) Element {
	if len(bytes) != 12 {
		return TimestampPair{}
	}

	return TimestampPair{
		id:        0,
		timestamp: 0,
	}
}

type TimestampIndex struct {
	SkipList
}

func GenerateTimestampIndex() *TimestampIndex {
	skipList :=Generate(15, TimePairCreate)
	return &TimestampIndex{*skipList}
}

func (index *TimestampIndex) Add(id DocId, timeStamp uint32) {
	p := TimestampPair{
		id:        id,
		timestamp: timeStamp,
	}

	index.Insert(p)
}

func (index *TimestampIndex) FindTimestamp(timeStamp uint32) TimestampPair {
	time := TimestampPair{
		id:        0,
		timestamp: timeStamp,
	}

	p := index.Find(time).value
	time.id = p.(TimestampPair).id
	return time
}

func (index *TimestampIndex) Size() uint32 {
	return index.total
}

func (index *TimestampIndex) IterateAll() []TimestampPair {
	var pairs []TimestampPair
	all := index.Iterator(index.Size())
	for _, e := range all {
		p, ok := e.(TimestampPair)
		if ok {
			pairs = append(pairs, p)
		}
	}

	return pairs
}

func GenerateTimeIndexFromBytes(bytes []byte) *TimestampIndex{
	skip := SkipListFromBytes(bytes, TimePairCreate, 12)
	return & TimestampIndex{*skip}
}