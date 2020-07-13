package indexer

import (
	"fmt"
	"strconv"
	"strings"
)

/**
评分系统：总的设计类似github上加star，如果在搜索的时候，输入 'abcdef'的时候，对某个搜索结果
star的时候，就会相应的记录在这里, 以倒排的方式索引
 */

type DocStar struct {
	Id DocId
	star int
}

func (i DocStar) Compare(other interface{}) int {
	return i.star  - other.(*DocStar).star
}

func (i DocStar) Serialize() []byte {
	return nil
}

func (i DocStar) ByteSize() uint8 {
	return 4
}

func (i DocStar) UnSerialize(bytes []byte) {

}

type StarRating struct {
	table map[string] *SkipList
	ids map[string] *SkipListNode  // key的格式为 `${word:docId}`这个用来根据id反向定位到skiplist节点
}

func GeneStarSys() *StarRating {
	return &StarRating{
		table: make(map[string] *SkipList, 0),
		ids: make(map[string] *SkipListNode, 0),
	}
}

// todo 该函数应该放在skiplist中
// 在skiplist中查找doc，先通过star找到，再进行下找
func(r *StarRating) findNode(ds DocStar, list *SkipList) *SkipListNode {
	var node *SkipListNode
	star := ds.star
	docId := ds.Id

	n := list.Find(ds)
	n = list.DiveLowest(n)
	if n != nil { // 找到该star数的doc不一定就是该doc，要迭代查询直到docId相等
		for n.next != nil && n.value.(DocStar).Id - docId != 0 {
			if n.value.(DocStar).star != star { // 遇到star不等就停止
				break
			}
		}
		if n.next != nil { // 找到
			return n
		}

		for n.pre.value != MAXELEMENT && n.value.(DocStar).Id - docId != 0 {
			if n.value.(DocStar).star != star {
				break
			}
		}
		if n.pre.value != MAXELEMENT {
			return n
		}
	}

	return node
}

func(r *StarRating) encodeKey(id DocId, word string) string {
	return fmt.Sprintf("%s:%d", word, id)
}

func(r *StarRating) decodeKey(key string) (id DocId, word string) {
	split := strings.Split(key, ":")
	if len(split) != 2 {
		return
	}

	word = split[1]
	i, err := strconv.ParseInt(split[0], 10, 64)
	if err != nil {
		return
	}
	id = DocId(i)
	return
}

func(r *StarRating) saveId(n *SkipListNode, id DocId, word string) {
	w := r.encodeKey(id, word)
	r.ids[w] = n
}

// @description :给某篇文章针对某个搜索词加星星
// @param docId: 文章id， word: 用户输入的搜索词条，preStar: 之前的star数字
// @return success: true成功, false失败; err: 错误信息
func (r *StarRating) Star (docId DocId, word string) (success bool, err error) {
	wd := fmt.Sprintf("%s:%d", word, docId)
	idn := r.ids[wd]
	list := r.table[word]
	var ds *DocStar
	if list == nil {
		list = Generate(15, func(bytes []byte) Element {
			return nil
		})
		r.table[word] = list
	}
	if idn == nil {
		ds = &DocStar{
			Id: docId,
			star:  1,
		}
	} else {
		// 删除节点
		list.DeleteNode(idn)
		ds = &DocStar{
			Id: docId,
			star:  idn.value.(*DocStar).star + 1,
		}
	}
	n := list.Insert(ds)
	r.saveId(n, docId, word)
	return true, nil
}

func (r *StarRating) Unstar(docId DocId, word string) (success bool, err error) {
	list := r.table[word]
	w := r.encodeKey(docId, word)
	idn := r.ids[w]
	var ds *DocStar
	if idn == nil {
		ds = &DocStar{
			Id: docId,
			star:  -1,
		}
	} else {
		list.DeleteNode(idn)
		ds = &DocStar{
			Id: docId,
			star:  idn.value.(*DocStar).star - 1,
		}
	}
	n := list.Insert(ds)
	r.saveId(n, docId, word)
	return true, nil
}

// 根据用户的搜索词条进行取top5 star数的docs
func (r *StarRating) Top5(word string) []*DocStar {
	var res []*DocStar
	list := r.table[word]
	if list == nil {
		return res
	}
	els := list.Iterator(5)
	for _, e := range els {
		if e.(*DocStar).star > 0 {
			res = append(res, e.(*DocStar))
		}
	}

	return res
}

func (r *StarRating) GetStarWithId(id DocId, word string) int {
	w := r.encodeKey(id, word)
	n := r.ids[w]
	fmt.Println("getStart", n)
	if n != nil {
		return n.value.(*DocStar).star
	}

	return 0
}