package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/lbe-io/crypto-sdk-go/callback"
	"github.com/lbe-io/crypto-sdk-go/client"
	"github.com/lbe-io/crypto-sdk-go/common/entity"
)

func main() {
	// 1. 初始化SDK客户端
	config := client.SDKConfig{
		ServerHost: "https://api.example.com",
		MerchantID: "your_merchant_id",
		APIKey:     "your_api_key",
		APISecret:  "your_api_secret",
	}

	sdkClient, err := client.NewClient(config)
	if err != nil {
		log.Fatal("初始化SDK失败:", err)
	}

	// 2. 创建回调管理器
	callbackManager := callback.NewCallbackManager(sdkClient)

	// 3. 注册交易创建事件处理器
	callbackManager.Register(entity.EventTypeTransactionCreated, handleTransactionCreated)

	// 4. 注册钱包创建事件处理器
	callbackManager.Register(entity.EventTypeWalletCreated, handleWalletCreated)

	// 5. 注册测试事件处理器
	callbackManager.Register(entity.EventTypeTest, handleTestEvent)

	// 6. 创建并启动HTTP服务器
	httpHandler := callback.NewHTTPHandler(callbackManager)

	// 可以添加其他业务路由
	mux := http.NewServeMux()
	mux.Handle("/callback", httpHandler.HandlerFunc())
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	fmt.Println("回调服务器启动在 http://localhost:8080")
	fmt.Println("回调地址: POST http://localhost:8080/callback")

	log.Fatal(http.ListenAndServe(":8080", mux))
}

// handleTransactionCreated 处理交易创建事件
func handleTransactionCreated(ctx context.Context, data []byte) (any, error) {
	var eventData entity.TransactionCreatedEvent
	if err := json.Unmarshal(data, &eventData); err != nil {
		return nil, fmt.Errorf("解析交易数据失败: %v", err)
	}

	// 在这里添加你的业务逻辑
	// 例如：保存到数据库、发送通知、更新缓存等
	fmt.Printf("收到交易创建事件:\n")
	fmt.Printf("  交易ID: %s\n", eventData.TransactionID)
	fmt.Printf("  金额: %s %s\n", eventData.Amount, eventData.Currency)
	fmt.Printf("  状态: %s\n", eventData.Status)

	// 返回处理结果
	return entity.CallbackResponse{
		Code: 0,
		Msg:  "交易处理成功",
		Data: map[string]interface{}{
			"transaction_id": eventData.TransactionID,
			"processed_at":   "2024-01-01T12:00:00Z",
		},
	}, nil
}

// handleWalletCreated 处理钱包创建事件
func handleWalletCreated(ctx context.Context, data []byte) (any, error) {
	var eventData entity.WalletCreatedEvent
	if err := json.Unmarshal(data, &eventData); err != nil {
		return nil, fmt.Errorf("解析钱包数据失败: %v", err)
	}

	fmt.Printf("收到钱包创建事件:\n")
	fmt.Printf("  钱包ID: %s\n", eventData.WalletID)
	fmt.Printf("  钱包名称: %s\n", eventData.WalletName)
	fmt.Printf("  币种: %s\n", eventData.Currency)

	// 返回处理结果
	return entity.CallbackResponse{
		Code: 0,
		Msg:  "钱包创建成功",
		Data: eventData,
	}, nil
}

// handleTestEvent 处理测试事件
func handleTestEvent(ctx context.Context, data []byte) (any, error) {
	var eventData entity.TestEvent
	if err := json.Unmarshal(data, &eventData); err != nil {
		return nil, fmt.Errorf("解析测试数据失败: %v", err)
	}

	fmt.Printf("收到测试事件: %s\n", eventData.Message)

	// 返回处理结果
	return entity.CallbackResponse{
		Code: 0,
		Msg:  "测试成功",
		Data: map[string]interface{}{
			"message": eventData.Message,
			"time":    eventData.Time,
			"status":  "processed",
		},
	}, nil
}
