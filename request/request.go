package request

import (
	"context"

	"github.com/lbe-io/crypto-sdk-go/client"
	"github.com/lbe-io/crypto-sdk-go/common/entity"
)

// ICrypto 主接口
type IRequest interface {
	TransactionDetail(ctx context.Context, req *entity.TransactionDetailReq) (*entity.TransactionDetailResp, error)
	WalletList(ctx context.Context, req *entity.WalletListReq) (*entity.WalletListResp, error)
	WalletAddressGenerate(ctx context.Context, req *entity.WalletAddressGenerateReq) (*entity.WalletAddressGenerateResp, error)
	WalletAddressList(ctx context.Context, req *entity.WalletAddressListReq) (*entity.WalletAddressListResp, error)
	WalletAddressAssets(ctx context.Context, req *entity.WalletAddressAssetsReq) (*entity.WalletAddressAssetsResp, error)
}

type crypto struct {
	client *client.CryptoClient
}

func NewRequest(client *client.CryptoClient) IRequest {
	return &crypto{
		client: client,
	}
}

// 交易详情查询
func (c *crypto) TransactionDetail(ctx context.Context, req *entity.TransactionDetailReq) (*entity.TransactionDetailResp, error) {

	resp := &entity.TransactionDetailResp{}

	err := c.client.Post(ctx, "/path", req, resp)
	if err != nil {
		c.client.Logger().Errorf(ctx, "TransactionDetail err:%v", err)
		return nil, err
	}

	return resp, nil
}

// 钱包列表查询
func (c *crypto) WalletList(ctx context.Context, req *entity.WalletListReq) (*entity.WalletListResp, error) {

	return &entity.WalletListResp{}, nil
}

// 钱包地址生成
func (c *crypto) WalletAddressGenerate(ctx context.Context, req *entity.WalletAddressGenerateReq) (*entity.WalletAddressGenerateResp, error) {

	return &entity.WalletAddressGenerateResp{}, nil
}

// 钱包地址列表查询
func (c *crypto) WalletAddressList(ctx context.Context, req *entity.WalletAddressListReq) (*entity.WalletAddressListResp, error) {

	return &entity.WalletAddressListResp{}, nil
}

// 钱包地址资产查询
func (c *crypto) WalletAddressAssets(ctx context.Context, req *entity.WalletAddressAssetsReq) (*entity.WalletAddressAssetsResp, error) {

	return &entity.WalletAddressAssetsResp{}, nil
}
