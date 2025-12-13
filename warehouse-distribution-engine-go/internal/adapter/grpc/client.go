package grpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/antonchaban/warehouse-distribution-engine-go/gen/go/distribution"
	"github.com/antonchaban/warehouse-distribution-engine-go/internal/algorithm"
)

type Client struct {
	conn   *grpc.ClientConn
	remote pb.DistributionResultReceiverClient
}

func NewClient(javaServiceAddr string) (*Client, error) {
	conn, err := grpc.NewClient(
		javaServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to java service: %w", err)
	}

	client := pb.NewDistributionResultReceiverClient(conn)

	return &Client{
		conn:   conn,
		remote: client,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) SendPlan(ctx context.Context, plan algorithm.DistributionPlan, sourceID int64) error {
	type moveKey struct {
		WarehouseID string
		ProductID   string
	}

	aggregated := make(map[moveKey]struct {
		Count    int32
		VolumeM3 float64
	})

	for _, m := range plan.Moves {
		key := moveKey{
			WarehouseID: m.WarehouseID,
			ProductID:   m.ProductID,
		}

		current := aggregated[key]
		current.Count++
		current.VolumeM3 = m.VolumeM3
		aggregated[key] = current
	}

	protoMoves := make([]*pb.Move, 0, len(aggregated))
	for key, data := range aggregated {
		protoMoves = append(protoMoves, &pb.Move{
			ProductId:   key.ProductID,
			WarehouseId: key.WarehouseID,
			VolumeM3:    data.VolumeM3,
			Quantity:    data.Count,
		})
	}

	protoUnallocated := make([]*pb.UnallocatedItem, 0, len(plan.UnallocatedItems))
	for _, u := range plan.UnallocatedItems {
		protoUnallocated = append(protoUnallocated, &pb.UnallocatedItem{
			ProductId: u.ID,
			VolumeM3:  u.VolumeM3,
			Reason:    "insufficient_capacity",
		})
	}

	req := &pb.DistributionPlan{
		RequestId:        plan.RequestID,
		SourceId:         sourceID,
		Moves:            protoMoves,
		UnallocatedItems: protoUnallocated,
		GeneratedAt:      time.Now().Unix(),
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := c.remote.ProcessPlan(ctx, req)
	if err != nil {
		return fmt.Errorf("rpc call failed: %w", err)
	}

	return nil
}
