package implement

import (
	"context"
	"snake/indexer"
	"snake/proto"
	"strconv"
)

type StarServer struct {
	StarRating *indexer.StarRating
}

func InitStarServer(sr *indexer.StarRating) *StarServer {
	return &StarServer{
		sr,
	}
}

func (ss *StarServer) Star(ctx context.Context, ssReq *proto.StarRequest)(*proto.StarResponse, error) {
	w := ssReq.Word
	d := ssReq.DocId
	t := ssReq.Type
	var success bool
	id, err := strconv.ParseInt(d, 10, 64)
	if err != nil {
		return nil, err
	}
	if t == 1 { //1 为加星操作
		success, _ = ss.StarRating.Star(indexer.DocId(id), w)
	} else {
		success, _ = ss.StarRating.Unstar(indexer.DocId(id), w)
	}
	sp := &proto.StarResponse{
		IsSuccess:            success,
	}

	return sp, nil
}


