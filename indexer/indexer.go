package indexer

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
)

type Indexer struct {
	StarRating
	table map[string] *KeywordIndices
}

func Init() *Indexer {
	return &Indexer{
		StarRating: *GeneStarSys(),
		table: make(map[string]*KeywordIndices, 10000000),
	}
}

func(indexer *Indexer) addDoc(doc *Doc, keyword string, isDynamic bool) {
	docIndice, ok := indexer.table[keyword]
	if ok {
		docIndice.Add(*doc)
	} else {
		indexer.table[keyword] = GenerateIndices()
		if isDynamic {
			indexer.table[keyword].AddDynamic(*doc)
		} else {
			indexer.table[keyword].Add(*doc)
		}
	}
}

// 顺序添加文章
func (indexer *Indexer) AddDocOrderly(doc *Doc, keywords []string) {
	for _ , keyword := range keywords {
		if keyword != "" {
			indexer.addDoc(doc, keyword, false)
		}
	}
}

// 动态插入一个文章，在增量merge 新文章时调用
func (indexer *Indexer) AddDocDynamic(doc *Doc, keywords []string) {
	for _, keyword := range keywords {
		indexer.addDoc(doc, keyword, true)
	}
}

type FoundIndices [] *KeywordIndices

func (f FoundIndices) Len() int {
	return len(f)
}

func (f FoundIndices) Less(i, j int) bool {
	return len(f[i].docIds) < len(f[j].docIds)
}

func (f FoundIndices) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

// todo 添加频次 相关度等指标
// 先查找star数为top5的放在前列
func (indexer *Indexer) Search(keywords []string, timeSort bool) []DocId {
	var res []DocId
	var found []*KeywordIndices
	for _, keyword := range keywords {
		indices, ok := indexer.table[keyword]
		if ok {
			found = append(found, indices)
		}
	}

	if len(found) > 0 {
		sort.Sort(FoundIndices(found))

		// 进行merge操作，以最少那个为标准merge操作
		docIds := found[0].IterateDocIds(timeSort)
		for _, id := range docIds {
			foundCount := 1
			for i:=1; i < len(found); i ++ {
				if found[i].Contains(id) {
					foundCount = foundCount + 1
				}
			}
			if foundCount == len(found) {
				res = append(res, id)
			}
		}
	}

	return res
}

func (indexer *Indexer) uint64bytes(u uint64) []byte {
	var bytes []byte
	for i:=uint64(0); i < 8; i ++ {
		bytes = append(bytes, uint8(u >> (8 * i)))
	}
	return bytes
}

func (indexer *Indexer) uint16bytes(u uint16) []byte {
	var bytes []byte
	for i := uint16(0); i<2; i ++ {
		bytes = append(bytes, uint8(u >> (8 * i)))
	}
	return bytes
}

func (indexer *Indexer) uint24bytes(u uint32) []byte {
	var bytes []byte
	for i:= uint32(0); i < 3; i ++ {
		bytes = append(bytes, uint8(u >> (8 * i)))
	}
	return bytes
}

func (indexer *Indexer) uint32bytes(u uint32) []byte {
	var bytes []byte
	for i:= uint32(0); i < 4; i ++ {
		bytes = append(bytes, uint8(u >> (8 * i)))
	}
	return bytes
}

func (indexer *Indexer) uint16(bytes []byte) uint16 {
	return uint16(bytes[0]) | uint16(bytes[1]) << 8
}

func (indexer *Indexer) uint24(bytes []byte) uint32 {
	return uint32(bytes[0]) | uint32(bytes[1]) << 8 | uint32(bytes[2]) << 16
}

func (indexer *Indexer) uint32(bytes []byte) uint32 {
	return uint32(bytes[0]) | uint32(bytes[1]) << 8 | uint32(bytes[2]) << 16 | uint32(bytes[3]) << 24
}

func (indexer *Indexer) uint64(bytes []byte) uint64 {
	return uint64(bytes[0]) | uint64(bytes[1]) << 8 | uint64(bytes[2]) << 16 | uint64(bytes[3]) << 24 |
		uint64(bytes[4]) << 32 | uint64(bytes[5]) << 40 | uint64(bytes[6]) << 48 | uint64(bytes[7]) << 56
}

// 将倒排整理成可以写入磁盘的形式
// format: word1 + d1...dn + word2 + d21 ...d2n
func (indexer *Indexer) arrangeMemory() []byte {
	buf := bytes.NewBuffer([]byte{})
	for k, v := range indexer.table {
		kbyte := uint8(len(k))
		if kbyte == 0 {
			fmt.Println("kbyte", k, v)
		}
		buf.WriteByte(kbyte) // 写入k的长度
		buf.WriteString(k)   // 写入k的string值
		vbytes := indexer.uint32bytes(uint32(len(v.docIds))) // 写入value的长度
		buf.Write(vbytes)
		for _, id := range v.docIds {
			idbytes := indexer.uint64bytes(uint64(id)) // 写入docid
			buf.Write(idbytes)
		}
	}

	return buf.Bytes()
}

func (indexer *Indexer) FlushDisk() {
	bytes := indexer.arrangeMemory()
	file, _ := os.Create("./index.data")
	defer file.Close()
	file.Write(bytes)
}

func (indexer *Indexer) LoadFromDisk() {
	buf, err := ioutil.ReadFile("./index.data")
	if err != nil {
		log.Fatal("loadFromDisk,error:" , err)
	}

	for len(buf) > 0 {
		size := buf[0]
		buf = buf[1:]
		if size == 0 {
			fmt.Println("aaa", size)
			continue
		}
		word := string(buf[0: size])
		buf = buf[size: ]
		dlen := indexer.uint32(buf[0:4])
		buf = buf[4: ]
		var docIds []DocId
		fmt.Println("dlen", dlen, word, len(word))
		for i:= 0; i < int(dlen); i++ {
			id := indexer.uint64(buf[i*8: i*8 + 8])
			docIds = append(docIds, DocId(id))
		}
		buf = buf[dlen * 8 : ]
		indexer.table[word] = &KeywordIndices{
			docIds: docIds,
		}
	}
}
