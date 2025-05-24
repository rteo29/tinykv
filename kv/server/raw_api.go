package server

import (
	"context"

	"github.com/pingcap-incubator/tinykv/kv/storage"
	"github.com/pingcap-incubator/tinykv/proto/pkg/kvrpcpb"
)

// The functions below are Server's Raw API. (implements TinyKvServer).
// Some helper methods can be found in sever.go in the current directory

// RawGet return the corresponding Get response based on RawGetRequest's CF and Key fields
func (server *Server) RawGet(_ context.Context, req *kvrpcpb.RawGetRequest) (*kvrpcpb.RawGetResponse, error) {
	reader, err := server.storage.Reader(req.Context)
	if err != nil {
		return nil, err
	}

	val, err := reader.GetCF(req.Cf, req.Key)
	resp := kvrpcpb.RawGetResponse{
		Value: val,
	}
	if err != nil {
		resp.Error = err.Error()
	} else {
		resp.NotFound = val == nil
	}
	return &resp, nil
}

// RawPut puts the target data into storage and returns the corresponding response
func (server *Server) RawPut(_ context.Context, req *kvrpcpb.RawPutRequest) (*kvrpcpb.RawPutResponse, error) {
	put := storage.Put{
		Key:   req.Key,
		Value: req.Value,
		Cf:    req.Cf,
	}
	modify := storage.Modify{
		Data: put,
	}
	batch := []storage.Modify{modify}
	err := server.storage.Write(req.Context, batch)
	resp := &kvrpcpb.RawPutResponse{}
	if err != nil {
		resp.Error = err.Error()
	}
	return resp, err
}

// RawDelete delete the target data from storage and returns the corresponding response
func (server *Server) RawDelete(_ context.Context, req *kvrpcpb.RawDeleteRequest) (*kvrpcpb.RawDeleteResponse, error) {
	del := storage.Delete{
		Key: req.Key,
		Cf:  req.Cf,
	}
	modify := storage.Modify{
		Data: del,
	}
	batch := make([]storage.Modify, 0)
	batch = append(batch, modify)
	err := server.storage.Write(req.Context, batch)
	return &kvrpcpb.RawDeleteResponse{}, err
}

// RawScan scan the data starting from the start key up to limit. and return the corresponding result
func (server *Server) RawScan(_ context.Context, req *kvrpcpb.RawScanRequest) (*kvrpcpb.RawScanResponse, error) {
	reader, err := server.storage.Reader(req.Context)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	itr := reader.IterCF(req.Cf)
	defer itr.Close()

	itr.Seek(req.StartKey)

	var kvs []*kvrpcpb.KvPair
	limit := req.Limit
	for itr.Valid() && limit > 0 {
		limit = limit - 1
		item := itr.Item()
		key := item.KeyCopy(nil)
		value, err := item.ValueCopy(nil)
		if err != nil {
			kvs = append(kvs, &kvrpcpb.KvPair{
				Key:   key,
				Error: &kvrpcpb.KeyError{},
			})
		} else {
			kvs = append(kvs, &kvrpcpb.KvPair{
				Key:   key,
				Value: value,
			})
		}
		itr.Next()
	}
	return &kvrpcpb.RawScanResponse{
		Kvs: kvs,
	}, nil
}
