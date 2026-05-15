package callback

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lbe-io/crypto-sdk-go/client"
	"github.com/lbe-io/crypto-sdk-go/common/consts"
	"github.com/lbe-io/crypto-sdk-go/common/entity"
)

// HandlerFunc 回调处理函数类型
type HandlerFunc func(ctx context.Context, data []byte) (any, error)

// CallbackManager 回调管理器
type CallbackManager struct {
	callbackUri string
	client      *client.CryptoClient
	handlers    map[consts.CallbackEventType]HandlerFunc
}

// NewCallbackManager 创建回调管理器
// callbackUri 可选参数，指定回调的URI路径，默认为请求的URI路径。
func NewCallbackManager(client *client.CryptoClient, callbackUri string) *CallbackManager {
	return &CallbackManager{
		callbackUri: callbackUri,
		client:      client,
		handlers:    make(map[consts.CallbackEventType]HandlerFunc),
	}
}

// Register 注册回调处理器
func (m *CallbackManager) Register(eventType consts.CallbackEventType, handler HandlerFunc) {
	m.handlers[eventType] = handler
}

// HandleRequest 处理HTTP请求
func (m *CallbackManager) HandleRequest(ctx context.Context, req *http.Request) (any, error) {
	// 验证签名并获取原始数据
	body, err := m.client.HandleCallback(req, m.callbackUri)
	if err != nil {
		return nil, fmt.Errorf("回调签名验证失败: %v", err)
	}

	// 解析回调事件
	var callbackEvent entity.CallbackEvent
	if err := json.Unmarshal(body, &callbackEvent); err != nil {
		return nil, fmt.Errorf("解析回调事件失败: %v", err)
	}

	// 查找对应的处理器
	handler, exists := m.handlers[callbackEvent.EventType]
	if !exists {
		return nil, fmt.Errorf("未注册的事件类型: %s", callbackEvent.EventType)
	}

	// 调用处理器
	return handler(ctx, callbackEvent.Data)
}

// HandleRawData 处理原始数据（已验证签名的数据）
func (m *CallbackManager) HandleRawData(ctx context.Context, eventType consts.CallbackEventType, data []byte) (any, error) {
	handler, exists := m.handlers[eventType]
	if !exists {
		return nil, fmt.Errorf("未注册的事件类型: %s", eventType)
	}
	return handler(ctx, data)
}

// DefaultHandlers 默认处理器
var DefaultHandlers = struct {
	// TransactionCreated 交易创建事件处理器
	TransactionCreated HandlerFunc
	// TransactionUpdated 交易更新事件处理器
	TransactionUpdated HandlerFunc
	// WalletCreated 钱包创建事件处理器
	WalletCreated HandlerFunc
	// WalletUpdated 钱包更新事件处理器
	WalletUpdated HandlerFunc
	// AddressGenerated 地址生成事件处理器
	AddressGenerated HandlerFunc
	// TestEvent 测试事件处理器
	TestEvent HandlerFunc
}{
	TransactionCreated: func(ctx context.Context, data []byte) (any, error) {
		var eventData entity.TransactionCreatedEvent
		if err := json.Unmarshal(data, &eventData); err != nil {
			return nil, fmt.Errorf("解析交易创建事件数据失败: %v", err)
		}
		return entity.CallbackResponse{
			Code: 0,
			Msg:  "交易创建事件处理成功",
			Data: eventData,
		}, nil
	},
	TransactionUpdated: func(ctx context.Context, data []byte) (any, error) {
		var eventData entity.TransactionUpdatedEvent
		if err := json.Unmarshal(data, &eventData); err != nil {
			return nil, fmt.Errorf("解析交易更新事件数据失败: %v", err)
		}
		return entity.CallbackResponse{
			Code: 0,
			Msg:  "交易更新事件处理成功",
			Data: eventData,
		}, nil
	},
	WalletCreated: func(ctx context.Context, data []byte) (any, error) {
		var eventData entity.WalletCreatedEvent
		if err := json.Unmarshal(data, &eventData); err != nil {
			return nil, fmt.Errorf("解析钱包创建事件数据失败: %v", err)
		}
		return entity.CallbackResponse{
			Code: 0,
			Msg:  "钱包创建事件处理成功",
			Data: eventData,
		}, nil
	},
	WalletUpdated: func(ctx context.Context, data []byte) (any, error) {
		var eventData entity.WalletUpdatedEvent
		if err := json.Unmarshal(data, &eventData); err != nil {
			return nil, fmt.Errorf("解析钱包更新事件数据失败: %v", err)
		}
		return entity.CallbackResponse{
			Code: 0,
			Msg:  "钱包更新事件处理成功",
			Data: eventData,
		}, nil
	},
	AddressGenerated: func(ctx context.Context, data []byte) (any, error) {
		var eventData entity.AddressGeneratedEvent
		if err := json.Unmarshal(data, &eventData); err != nil {
			return nil, fmt.Errorf("解析地址生成事件数据失败: %v", err)
		}
		return entity.CallbackResponse{
			Code: 0,
			Msg:  "地址生成事件处理成功",
			Data: eventData,
		}, nil
	},
	TestEvent: func(ctx context.Context, data []byte) (any, error) {
		var eventData entity.TestEvent
		if err := json.Unmarshal(data, &eventData); err != nil {
			return nil, fmt.Errorf("解析测试事件数据失败: %v", err)
		}
		return entity.CallbackResponse{
			Code: 0,
			Msg:  "测试事件处理成功",
			Data: eventData,
		}, nil
	},
}
