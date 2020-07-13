package implement

import (
	"context"
	"snake/indexer"
	"snake/jieba"
	search "snake/proto"
	"snake/store"
	"strconv"
)

type SearchServer struct {
	Store store.Store
	Indexer *indexer.Indexer
}

func (ss *SearchServer) cutWords(s string)[]string {
	return jieba.Cut(s)
}

// 根据切分后的word 查找id列表
// sort 为1 代表：默认排序， 2：时间倒序
func (ss *SearchServer) searchIds(words []string, sort int32) []indexer.DocId {
	var docIds []indexer.DocId
	if sort == 2 {
		docIds = ss.Indexer.Search(words, true)
	} else {
		docIds = ss.Indexer.Search(words, false)
	}
	return docIds
}

func (ss *SearchServer) merge(docIds []indexer.DocId, tops []indexer.DocId) []DocInfo {
	var ids []DocInfo
	idsMap := make(map[indexer.DocId]interface{}, 0)
	for _, top := range tops {
		docInfo := DocInfo{
			docId:  top,
			isTop5: true,
		}
		ids = append(ids, docInfo)
		idsMap[top] = struct{}{}
	}

	for _, t := range docIds {
		id := idsMap[t]
		if id == nil {
			idsMap[t] = struct{}{}
			ids = append(ids, DocInfo{
				docId:  t,
				isTop5: false,
			})
		}
	}

	return ids
}

func (ss *SearchServer) pagination(ids []DocInfo, page int32) []DocInfo{
	if int32(len(ids)) >= page * 10 {
		ids = ids[page*10-10 : page*10 - 1]
	} else {
		ids = ids[page*10-10 :]
	}

	return ids
}

func (ss *SearchServer) detail(idInfo DocInfo, word string) *indexer.Doc {
	id := strconv.FormatUint(uint64(idInfo.docId), 10)
	doc := ss.Store.Get(id)
	doc.Star = ss.Indexer.GetStarWithId(idInfo.docId, word)
	doc.Id = id
	doc.IsTop5 = idInfo.isTop5

	return &doc
}

// 取出 搜索词条下的top5
func (ss *SearchServer) top5(word string) []indexer.DocId {
	tops :=ss.Indexer.Top5(word)
	var ids []indexer.DocId
	for _, top := range tops {
		id := indexer.DocId(int64(top.Id))
		ids = append(ids, id)
	}
	return ids
}

func(ss *SearchServer) formatResult(detail *indexer.Doc) *search.Result {

	return &search.Result{
		Id:                   detail.Id,
		Title:                detail.Title,
		Url:                  detail.Url,
		Lang: uint32(detail.Lang),
		Favicon:              detail.Favicon,
		IsTop5:               detail.IsTop5,
		Star:                 int32(detail.Star),
		TimeStamp:            detail.TimeStamp,
		Description:          detail.TitleOrDescription,
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}
}

func (ss *SearchServer) Search(ctx context.Context, sReq *search.SearchRequest) (*search.SearchResponse, error) {
	word := sReq.Word
	page := sReq.Page
	sort := sReq.Sort

	words := ss.cutWords(word)
	mergedDocInfos := ss.merge(ss.searchIds(words, sort), ss.top5(word))
	pagedDocInfos := ss.pagination(mergedDocInfos, page)

	var docs []*search.Result
	for _, info := range pagedDocInfos {
		detail := ss.detail(info, word)
		docs = append(docs, ss.formatResult(detail))
	}
	sRes := &search.SearchResponse{
		Data: docs,
		Keywords: words,
		Success: true,
		Total: uint32(len(mergedDocInfos)),
	}

	return sRes, nil
}

func (ss *SearchServer) Detail(ctx context.Context, dReq *search.DetailRequest) (*search.DetailResponse, error) {
	id := dReq.Id
	d := ss.Store.Get(id)
	r := ss.formatResult(&d)
	return &search.DetailResponse{
		Data: r,
	}, nil
}

func (ss *SearchServer) Details(ctx context.Context, dsReq *search.DetailsRequest) (*search.DetailsResponse, error) {
	ids := dsReq.Ids
	var details []*search.Result
	for _, id := range ids {
		d := ss.Store.Get(id)
		if d.DocId != 0 {
			r := ss.formatResult(&d)
			details = append(details, r)
		}
	}
	return &search.DetailsResponse{
		Data:                 details,
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}, nil
}
type  DocInfo struct{
	docId indexer.DocId `json:"docId"`
	isTop5 bool `json:"isTop5"`
}
