package partition

import (
	pb "github.com/anujga/dstk/pkg/api/proto"
	"github.com/anujga/dstk/pkg/ss/partition"
)

func proxyToRetired(actor partition.Actor, partIdMap map[int64]partition.Actor, part *pb.Partition) interface{} {
	return &partition.Retire{}
}

