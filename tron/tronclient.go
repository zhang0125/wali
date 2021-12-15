package tron

import (
	"context"
	"fmt"

	"github.com/zhang0125/wali/tron/pb"
	"google.golang.org/grpc"
)

// Client defines typed wrappers for the Tron RPC API.
type Client struct {
	client pb.WalletClient
}

// NewWalletClient creates a client that uses the given RPC client.
func NewWalletClient(url string) pb.WalletClient {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	return pb.NewWalletClient(conn)
}

// NewClient creates a client that uses the given RPC client.
func NewClient(url string) *Client {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	return &Client{
		client: pb.NewWalletClient(conn),
	}
}

func (tc *Client) GetNowBlock() (*pb.BlockExtention, error) {
	return tc.client.GetNowBlock2(context.Background(), &pb.EmptyMessage{})
}

func (tc *Client) TriggerContract(ownerAddress, contractAddress, data []byte) (*pb.Transaction, error) {
	response, err := tc.client.TriggerContract(context.Background(),
		&pb.TriggerSmartContract{
			OwnerAddress:    ownerAddress,
			ContractAddress: contractAddress,
			CallValue:       100,
			Data:            data,
			CallTokenValue:  0,
			TokenId:         0,
		})
	if err != nil {
		return nil, err
	}
	if response.Result.Code != pb.Return_SUCCESS {
		return nil, fmt.Errorf("code:%v message:%v", response.Result.Code, string(response.Result.Message))
	}
	return response.Transaction, nil
}

func (tc *Client) TriggerConstantContract(ownerAddress, contractAddress, data []byte) ([]byte, error) {
	response, err := tc.client.TriggerConstantContract(context.Background(),
		&pb.TriggerSmartContract{
			OwnerAddress:    ownerAddress,
			ContractAddress: contractAddress,
			CallValue:       100,
			Data:            data,
			CallTokenValue:  0,
			TokenId:         0,
		})
	if err != nil {
		return nil, err
	}
	if response.Result.Code != pb.Return_SUCCESS || response.Transaction.GetRet()[0].Ret == pb.Transaction_Result_FAILED {
		return nil, fmt.Errorf("code:%v message:%v", response.Result.Code, string(response.Result.Message))
	}
	fmt.Println("result:", response.Result.Result, "code:", response.Result.Code, " message:", string(response.Result.Message))
	return response.ConstantResult[0], nil
}

func (tc *Client) BroadcastTransaction(ctx context.Context, trx *pb.Transaction) error {
	result, err := tc.client.BroadcastTransaction(ctx, trx)
	if err != nil {
		return err
	}
	if result.Code != pb.Return_SUCCESS {
		return fmt.Errorf("code:%v message:%v", result.Code, string(result.Message))
	}
	return nil
}
