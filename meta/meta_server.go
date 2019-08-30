package meta

import (
	"context"

	"github.com/climbcomp/climbcomp-go/climbcomp"
	meta_pb "github.com/climbcomp/climbcomp-go/climbcomp/meta/v1"
)

// NewMetaServer creates a new server
func NewMetaServer() *MetaServer {
	return &MetaServer{}
}

type MetaServer struct {
}

func (s *MetaServer) GetVersion(ctx context.Context, req *meta_pb.GetVersionRequest) (*meta_pb.GetVersionResponse, error) {
	resp := &meta_pb.GetVersionResponse{
		Version: climbcomp.VERSION,
	}
	return resp, nil
}
