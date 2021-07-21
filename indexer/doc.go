package indexer

import (
	"bytes"
	"snake/util"
)

func (d DocId) Op_Subtract(o DocId) uint64 {
	return uint64(d) - uint64(o)
}

type Doc struct {
	Id string `json:"id"`  // 其实就是docId的值，为了防止int64精度丢失
	DocId DocId  `json:"docId"`// DocId为文档的唯一标识，用ID生成器生成，保证递增
	Url string	`json:"url"`
	Lang uint8 `json:"lang"`
	Title string `json:"title"`
	TimeStamp uint32 `json:"time_stamp"` // 文章的创建时间
	Favicon string `json:"favicon"`
	TitleOrDescription string `json:"title_or_description"` // 文章的title或者description
	//Score int8  // 存储score
	Star int	`json:"star"`
	IsTop5 bool `json:"isTop5"`
}

type SnakeId struct {
	Id int `json:"id"`
}

type DocIndex struct {
	DocId DocId
	Keywords []*KeywordIndex
}

type KeywordIndex struct {
	Text string
	Frequency float32
	Starts []int
}

// 每个key 对应的后面的索引项包含docid自增数组，timeStamp索引
type KeywordIndices struct {
	docIds []DocId
	//timestampIndex *TimestampIndex  时间排序先去掉
}

type SnakeIdIndices struct {
	snakeIds []SnakeId
}

func GenerateIdIndices() * SnakeIdIndices {
	return &SnakeIdIndices{
		snakeIds: make([]SnakeId, 0),
	}
}


func (indices *SnakeIdIndices) Add(id SnakeId) {
	indices.snakeIds = append(indices.snakeIds, id)
	//indices.timestampIndex.Add(doc.DocId, doc.TimeStamp)  // 时间排序先去掉
}

func GenerateIndices() *KeywordIndices {
	return &KeywordIndices{
		docIds:         make([]DocId, 0),
		//timestampIndex: GenerateTimestampIndex(), // 时间排序先去掉
	}
}

func (indices *KeywordIndices) ByteSize() uint32 { // 返回total 字节长度
	//indexBytes := indices.timestampIndex.SeriliazedByteSize()  // 时间排序先去掉
	docIdBytes := uint32(len(indices.docIds) * 8)

	//return indexBytes + docIdBytes // 时间排序先去掉
	return docIdBytes
}

func (indices *KeywordIndices) Add(doc Doc) {
	indices.docIds = append(indices.docIds, doc.DocId)
	//indices.timestampIndex.Add(doc.DocId, doc.TimeStamp)  // 时间排序先去掉
}

func (indices *KeywordIndices) AddDynamic(doc Doc) {
	len := len(indices.docIds)
	if len == 0 {
		indices.docIds = append(indices.docIds, doc.DocId)
		return
	}
	i := indices.FindIndex(doc.DocId)
	var ids []DocId
	pre, next := indices.docIds[0: i], indices.docIds[i: len]
	ids = append(ids, pre...)
	ids = append(ids, doc.DocId)
	ids = append(ids, next...)
	indices.docIds = ids
	//indices.timestampIndex.Add(doc.DocId, doc.TimeStamp)  // 时间排序先去掉
}

// 将keywordIndice 序列化，总体的data有两部分组成 docIds 和 timestampIndex
// 写入的[]byte的格式:
// 1. 总[]byte的长度 4字节 数字为 （3）的长度 + (4)的长度
// 2. 总 docIds的长度 4字节 数字为 (3) 的长度
// 3. docids的数据 id1(8字节) + id2(8字节) + id3(8字节)......
// 4. timestampIndex的数据

func (indices *KeywordIndices) Flush() []byte{
	buf := bytes.NewBuffer([]byte{})
	//serialized := indices.timestampIndex.SerializeToBytes() // 时间排序先去掉
	total := indices.ByteSize()

	buf.Write(util.Uint32bytes(total)) // 1.写入总长度

	//buf.Write(util.Uint32bytes(total - uint32(len(serialized)))) // 2. 写入docid数据长度   时间排序先去掉
	buf.Write(util.Uint32bytes(total))
	for _, doc := range indices.docIds { // 3.写入docid 数据
		buf.Write(util.Uint64bytes(uint64(doc)))
	}
	//buf.Write(serialized) // 4.写入timestampIndex的数据    时间排序先去掉

	return buf.Bytes()
}

// bytes 不包括总数字这个数据
func GenerateIndiceFromBytes(bytes []byte) *KeywordIndices{
	i := &KeywordIndices{}
	offset := 0
	total := util.Uint32(bytes[offset: offset + 4]) // 读取 docId byte长度
	offset = offset + 4
	var docIds []DocId
	for i := uint32(0); i < total; i ++ {  // 读取docId数据
		id := util.Uint64(bytes[offset : offset + 8])
		docIds = append(docIds, DocId(id))
		offset = offset + 8
	}
	i.docIds = docIds

	//i.timestampIndex = GenerateTimeIndexFromBytes(bytes[offset: ]) // 读入timestamp index的数据  时间排序先去掉

	return i
}

// timeOrder 代表是否要按照时间排序返回 ids
func (indices *KeywordIndices) IterateDocIds(timeOrder bool) []DocId {
	var ids []DocId
	if timeOrder {
		//for _, pair := range indices.timestampIndex.IterateAll() {
		//	ids = append(ids, pair.id)
		//}
	} else {
		ids = indices.docIds
	}

	return ids
}
// 单词项对应的docId数组是否含有某id, 二分查找
func (indices *KeywordIndices) Contains(id DocId) bool {
	ids := indices.docIds
	if len(ids) == 0 {
		return false
	}

	if len(ids) == 1 && ids[0] == id {
		return true
	}

	if len(ids) == 2 && (ids[0] == id || ids[1] == id) {
		return true
	}

	start := 0
	end := len(ids) - 1
	if ids[start] > id || ids[end] < id {
		return false
	}

	for (end - start) > 1 {
		if ids[start] == id || ids[end] == id {
			return true
		}

		mid := (start + end) / 2
		if ids[mid] == id {
			return true
		}
		if ids[mid] > id {
			end = mid
		} else {
			start = mid
		}
	}
	return false
}

// 查找一个新DocId在docIds中的位置 i : indices.docIds[i] > id && indices.docIds[i-1] < id
func (indices *KeywordIndices) FindIndex(nid DocId) (i int) {
	if len(indices.docIds) == 0 {
		return 0
	}
	if indices.docIds[0] > nid {
		return 0
	}
	len := len(indices.docIds)
	if indices.docIds[len-1] < nid {
		return len
	}
	start := 0
	end := len
	for (end - start) > 1 {
		mid := (start + end) / 2
		if indices.docIds[mid] < nid {
			if indices.docIds[mid + 1] > nid {
				i = mid + 1
				break
			} else {
				start = mid
			}
		} else {
			if indices.docIds[mid - 1] < nid {
				i = mid
				break
			} else {
				end = mid
			}
		}
	}

	return i
}
