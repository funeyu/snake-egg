package indexer

import "errors"

type DocId uint64

// 生成一个唯一的文章id, 组成为 48 位的域名信息 + 文章在该域名的index信息 16位信息，假定一个独立博客的总有文章数不能超过2<<16
//  文章的域名不是一成不变的，后续会继续添加，为了保证之前的域名id不变，将48位的域名信息分成两部分
// rank_id: 文章域名的排序id，范围为：0 ~ 2^20
// sub_id: 范围为： 0 ~2^28
// i : 范围为： 0 ~ 2^16, 这里的值从65536依次递减，最旧的文章的i就是65536,最新的文章i为65536 -n(该博客的总文章数)
func GeneDocId(rank_id uint32, sub_rank_id uint32, i uint16) (DocId, error) {
	if rank_id > 2 << 20 {
		return 0, errors.New("rank_id 不能超出最大值：2 ^ 20")
	}
	if sub_rank_id > 2 << 28 {
		return 0, errors.New("sub_rank_id 不能超出最大值： 2 ^ 28")
	}

	num := uint64(rank_id) << 44 + uint64(sub_rank_id) << 16 + uint64(i)
	return DocId(num), nil
}

func (id DocId) RankId() uint32 {
	i := uint64(id)
	return uint32(i >> 44)
}

func (id DocId) SubRankId() uint32 {
	i := uint64(id)
	rankId := id.RankId()
	r := uint64(uint64(rankId) << 44)
	subId := i - r
	return uint32(subId >> 16)
}

func (id DocId) Index() uint16 {
	i := uint64(id)
	return uint16(i & (2 << 16 - 1))
}
