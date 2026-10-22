package entity

import "encoding/json"

// EventType 回调事件类型：
// 主链币/代币转入成功：INCOMING_CONFIRMED
// 主链币/代币转出成功：OUTGOING_CONFIRMED
// 主链币/代币转出失败：OUTGOING_FAILED
type EventType string

const (
	EventIncomingConfirmed EventType = "INCOMING_CONFIRMED"
	EventOutgoingConfirmed EventType = "OUTGOING_CONFIRMED"
	EventOutgoingFailed    EventType = "OUTGOING_FAILED"
	EventKeepalive         EventType = "KEEPALIVE"
)

// TransType 交易类型：主链币为 native，代币为对应代币标准。
// 当前暂时固定 Ethereum 代币为 erc20、Tron 代币为 trc20。
const (
	TransTypeNative = "native"
	TransTypeERC20  = "erc20"
	TransTypeTRC20  = "trc20"
)

type Event struct {
	Version  int       `json:"version,omitempty"` // 对外回调契约版本，供商户按版本分支解析
	WalletID string    `json:"wallet_id"`         // 钱包id
	Event    EventType `json:"event"`             // 事件类型 INCOMING_CONFIRMED, OUTGOING_CONFIRMED, OUTGOING_FAILED
	Item     EventItem `json:"item"`              // 事件
}

// EventItem
// 当一笔交易的转入转出地址均为主钱包相关地址时，需要发送2次事件：转入地址的incoming事件, 转出地址的outgoing事件
// 当转入或转出地址，包含手续费地址（主钱包零地址时），不发送该事件
type EventItem struct {
	Direction            Direction         `json:"direction"`                        // 交易方向，只支持incoming, outgoing
	Address              string            `json:"address"`                          // 产生回调事件的地址
	Blockchain           BlockchainName    `json:"blockchain"`                       // 区块链名称, 如：Ethereum, Solana, Binance Smart Chain, Tron
	Network              NetworkName       `json:"network"`                          // 网络名称, 如：mainnet, testnet, sepolia, nile
	TransactionHash      string            `json:"transaction_hash"`                 // 交易哈希
	TransType            string            `json:"trans_type"`                       // 交易类型：native；代币标准 erc20 / trc20
	BlockHeight          int64             `json:"block_height"`                     // 区块高度
	BlockTimestamp       int64             `json:"block_timestamp"`                  // 区块时间戳（秒）
	Status               TransactionStatus `json:"status"`                           // 交易状态, 如：success, failed
	Fee                  TransactionFee    `json:"fee"`                              // 交易Gas费用
	TokenTransfers       *TokenTransfer    `json:"token_transfers,omitempty"`        // 代币转账数量；仅代币交易出现
	NativeTokenTransfers *TokenTransfer    `json:"native_token_transfers,omitempty"` // 主链币转账信息；仅主链币交易出现，无合约地址

	// Extra 链特有字段扩展位。新增链需要暴露通用字段无法表达的信息时，走这里加法式扩展，
	// 既有字段与解析方不受影响。内容来自 asset_ledger.chain_extra。
	Extra map[string]any `json:"extra,omitempty"`
}

func (e *Event) Json() string {
	bytes, _ := json.Marshal(e)
	return string(bytes)
}
