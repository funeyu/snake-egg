package indexer

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/bits"
	"math/rand"
)

type Element interface {
	// 相等返回0， 大于other 返回 1 小于other 返回 -1
	Compare(other interface{}) int
	// 每个节点value可序列化
	Serialize() []byte
	ByteSize() uint8 	//约定跳跃表中存储的数据长度要小于256字节
}

type  MaxElement struct {}

func (m MaxElement) Compare(other interface{}) int{
	return 1
}

func (m MaxElement) Serialize()[]byte  {
	size := m.ByteSize()
	out := make([]byte, size)
	return out
}

func (m MaxElement) ByteSize() uint8 {
	return 1
}

var MAXELEMENT  = MaxElement{}


// 该跳跃表是从大到小的插入数据，为了方便取top数据，比如要取star数大的，要取timeStamp的
type SkipListNode struct {
	value Element
	pre *SkipListNode
	next *SkipListNode
	above *SkipListNode
	below *SkipListNode
	height uint8  //这个参数只在bottom node上才持有
}

func (node *SkipListNode) updateBottomNodeHeight(h uint8) {
	current := node
	for current.below != nil { // 先dive into the bottom node
		current = current.below
	}

	current.height = h
}

func (node *SkipListNode) increaseNodeHeight() {
	current := node
	for current.below != nil { // 先dive into the bottom node
		current = current.below
	}

	current.height = current.height + 1
}

type ElementCreate func(bytes []byte) Element

type  SkipList struct {
	CreateElement ElementCreate  // 每个skiplist都要有一个可以创建Element的函数
	nodeByteSize uint8   // skiplist value的byte size
	total uint32
	maxLevel uint8
	height uint8
	head []*SkipListNode
}

func Generate(maxLevel uint8, create ElementCreate) *SkipList {
	n := &SkipListNode{
		value: MAXELEMENT,
		height: 1,
	}
	var h []*SkipListNode
	h = append(h, n)
	return &SkipList {
		CreateElement: create,
		total: 0,
		maxLevel: maxLevel,
		height: 1,
		head: h,
	}
}

func (s *SkipList) level() uint8 {
	maxLevel := s.maxLevel
	level := maxLevel - 1
	x := rand.Uint64() & ((1 << uint(maxLevel-1)) - 1)
	zeroes := bits.TrailingZeros64(x)
	if uint8(zeroes) <= maxLevel {
		level = uint8(zeroes)
	}

	return level + 1
}

func (s *SkipList) isEmpty() bool {
	return s.head[0].next == nil
}

// 根据初始值新建一level，返回该level的 head节点指针
func (s *SkipList) newLevel(newNode *SkipListNode, belowHead *SkipListNode) *SkipListNode{
	head := &SkipListNode {
		value: MAXELEMENT, next: newNode, below: belowHead, pre: nil, above: nil,
	}
	head.increaseNodeHeight()

	belowHead.above = head
	newNode.pre = head
	return head
}

func (s *SkipList) increaseNumber(step uint32) {
	s.total = s.total + step
}

func (s *SkipList) decreaseNumber(step uint32) {
	s.total = s.total - step
}


func (s *SkipList) Insert(n Element) *SkipListNode {
	s.increaseNumber(1)

	node := &SkipListNode{
		value: n,
	}
	if s.isEmpty() {
		s.head[0].next = node
		node.pre = s.head[0]
	} else {
		walked := make([]*SkipListNode, s.height)  // 记录下探到最底层经历的节点
		c := s.head[s.height-1]
		for i := s.height; i > 0; i -- {
			for c != nil { // todo 换方式
				if c.value == n || c.next == nil {
					walked[i-1] = c
					c = c.below
					break
				}

				if c.next.value.Compare(n) >= 0 {
					c = c.next
				} else {
					walked[i-1] = c
					c = c.below
					break
				}
			}
		}
		l := s.level()
		height := uint8(1)
		for ; height <= l && height <= s.height + 1; height ++ {
			if height == s.height + 1 { // 代表新增的跳跃表已经到了突破之前的层数，需要新建一层
				head :=s.newLevel(node, s.head[s.height-1])
				s.head = append(s.head, head)
				s.height = s.height + 1
				break
			} else { // 双向链的insert操作
				next := walked[height-1].next
				walked[height-1].next = node
				node.next = next
				node.pre = walked[height-1]
				if next != nil {
					next.pre = node
				}
				nn := &SkipListNode{
					value: node.value,
					below: node,
				}
				node.above = nn
				node = nn
			}
		}
		node.updateBottomNodeHeight(height - 1)
	}

	return node
}

func (s *SkipList) Find(n Element) *SkipListNode{
	start := s.head[s.height -1]
	for i:= s.height; i > 0; i -- {
		for ; start.next != nil && start.next.value.Compare(n) >= 0; start = start.next {}
		if start.value == n {
			return start
		}
		start = start.below
	}
	if start !=nil && start.value == n {
		return start
	}
	return nil
}

func (s *SkipList) Delete(n Element) (success bool, err error) {
	found := s.Find(n)
	if found == nil {
		return false, errors.New(fmt.Sprintf("Can not find element"))
	}
	if found != nil {
		s.total = s.total - 1
	}

	for found != nil {
		found.pre.next = found.next
		if found.next != nil {
			found.next.pre = found.pre
		}
		found = found.below
	}
	return true, nil
}

// 删除某个节点
func(s *SkipList) DeleteNode(n *SkipListNode) {
	s.decreaseNumber(1)
	s.DeleteLowerNode(n)
	s.DeleteHigherNode(n)
}

// 从某个节点往低level开始删除
func(s *SkipList) DeleteLowerNode(n *SkipListNode) {
	for n != nil {
		if n.pre != nil {
			n.pre.next = n.next
		}
		if n.next != nil {
			n.next.pre = n.pre
		}
		tem := n
		n = n.below
		if n != nil {
			n.above = nil
		}
		tem.below = nil
	}
}

// 从节点往高处level删除
func (s *SkipList) DeleteHigherNode(n *SkipListNode) {
	for n != nil {
		if n.pre != nil {
			n.pre.next = n.next
		}
		if n.next != nil {
			n.next.pre = n.pre
		}
		if n.above != nil {
			n.above.below = nil
		}
		tem := n
		tem.above = nil
		n = n.above
	}
}

// 从一个节点递归到最低节点
func (s *SkipList) DiveLowest(n *SkipListNode) *SkipListNode {
	current := n
	if current.below != nil {
		current = n.below
	}
	return current
}

// 遍历取出前n个element, 如果超出total，取total个
func (s *SkipList) Iterator(n uint32) []Element {
	var res []Element
	if s.total < n {
		n = s.total
	}

	node := s.head[0].next // head[0]为max 最大值，为空， 所以要取next
	for i:=uint32(0); i < n ; i ++ {
		res = append(res, node.value)
		node = node.next
	}

	return res
}

func (s *SkipList) NodeByteSize() uint8 {
	return s.head[0].next.value.ByteSize()
}

func (s *SkipList) MaxElementSize() uint32 {
	return uint32(MAXELEMENT.ByteSize())
}

// 跳跃表序列化的格式为
// [total(4字节)]  标识跳跃表的元素个数
// [maxLevel(1字节)] 标识跳跃表的允许的最大高度
// [height(1字节)] 标识跳跃表的高度
// [height(1字节) + serialize(MinElement)] 存储最左节点信息
// [height(1字节) + serialize(n0)]--[height(1字节) + serialize(n1)]......  存储数据节点的信息
func (s *SkipList) SerializeToBytes() []byte {
	size := uint32(s.NodeByteSize())
	minElementSize := s.MaxElementSize()
	out := make([]byte, 7 + minElementSize + s.total * (size + 1))

	// 写入total数据
	binary.BigEndian.PutUint32(out[0:4], s.total)

	// 写入maxLevel
	copy(out[4:5], []byte{s.maxLevel})

	// 写入height
	copy(out[5:6], []byte{s.height})

	// 写入最左节点数据
	current := s.head[0]
	copy(out[6: 7], []byte{current.height})
	copy(out[7: 7 + minElementSize], MAXELEMENT.Serialize())

	// 写入数据节点的数据
	current = current.next
	offset :=  7 + minElementSize
	for current.next != nil {
		copy(out[offset : offset + 1 ], []byte{current.height})
		copy(out[offset + 1 : offset + 1 + size], current.value.Serialize())
		current = current.next
		offset = offset + size + 1
	}

	return out
}

// todo 一会32位数字，一会int，好乱
func (s *SkipList) totalBytes() uint32 {
	return 4
}

func (s *SkipList) maxLevelBytes() uint32 {
	return 1
}

func(s *SkipList) heightBytes() uint32 {
	return 1
}

func (s *SkipList) maxElementBytes() uint32 {
	return s.heightBytes() + uint32(MAXELEMENT.ByteSize())
}

func (s *SkipList) valueBytes() uint32 {
	return s.total * (uint32(s.nodeByteSize) + s.heightBytes())
}

func (s *SkipList) SeriliazedByteSize() uint32 { // 返回序列化后存储在disk的字节长度
	return s.totalBytes() + s.maxLevelBytes() + s.heightBytes() +  s.maxLevelBytes() + s.valueBytes()
}

// 形成跳跃表信息, 并将其返回
func (s *SkipList) formHead (bytes []byte) []*SkipListNode {
	var nodes [] *SkipListNode
	height := bytes[0]
	head := SkipListNode{
		value:  MAXELEMENT,
		pre:    nil,
		next:   nil,
		above:  nil,
		below:  nil,
		height: height,
	}
	nodes = append(nodes, &head)

	for i := uint8(1); i < height; i ++ {
		n := SkipListNode{
			value:  MAXELEMENT,
			pre:    nil,
			next:   nil,
			above:  nil,
			below:  nodes[i-1],
			height: 0,
		}
		nodes[i-1].above = &n
		nodes = append(nodes, &n)
	}

	s.head = nodes
	return nodes
}

// 返回 跳跃表从从该节点的右视图 节点， 一定会有最左节点存在，因为 最左节点是最高的
func (s *SkipList) formNode(bytes []byte, preNodes []*SkipListNode) []*SkipListNode {
	height := bytes[0]
	ele := s.CreateElement(bytes[1:])

	var res []*SkipListNode
	for _, node := range preNodes {
		res = append(res, node)
	}

	var below *SkipListNode
	for i := uint8(0); i < height; i ++ {
		node := &SkipListNode{
			value:  ele,
			pre:    res[i],
			next:   nil,
			above:  nil,
			below:  nil,
			height: 0,
		}
		res[i].next = node
		res[i] = node
		if i == 0 {
			node.height = height
		} else {
			node.below = below
			below.above = node
		}
		below = node
	}
	return res
}

func SkipListFromBytes(bytes []byte,  create ElementCreate, nodeByteSize uint8) *SkipList {
	s := &SkipList{
		nodeByteSize: nodeByteSize,
		CreateElement: create,
	}
	total := binary.BigEndian.Uint32(bytes[0:4])
	s.total = total
	s.maxLevel = bytes[4]
	s.height = bytes[5]

	preNodes := s.formHead(bytes[6: 7 + s.MaxElementSize()])
	offset := 7 + s.MaxElementSize()
	for i := uint32(0); i < s.total; i ++ {
		preNodes = s.formNode(bytes[offset: offset + uint32(nodeByteSize) + 1], preNodes)
		offset = offset + uint32(nodeByteSize) + 1
	}

	return s
}
