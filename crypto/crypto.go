package crypto

import (
	"context"
	"net/http"

	"github.com/lbe-io/crypto-sdk-go/client"
	"github.com/lbe-io/crypto-sdk-go/crypto/entity"
)

type ICrypto interface {
	TransactionDetail(ctx context.Context, req *entity.TransactionDetailReq) (*entity.TransactionDetailResp, error)
	WalletList(ctx context.Context, req *entity.WalletListReq) (*entity.WalletListResp, error)
	WalletAddressGenerate(ctx context.Context, req *entity.WalletAddressGenerateReq) (*entity.WalletAddressGenerateResp, error)
	WalletAddressList(ctx context.Context, req *entity.WalletAddressListReq) (*entity.WalletAddressListResp, error)
	WalletAddressAssets(ctx context.Context, req *entity.WalletAddressAssetsReq) (*entity.WalletAddressAssetsResp, error)
	Callback(ctx context.Context, req *http.Request,) (any, error)
}

type crypto struct {
	client *client.CryptoClient
}

func NewCrypto(client *client.CryptoClient) ICrypto {
	return &crypto{
		client: client,
	}
}
