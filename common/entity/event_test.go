package entity

import (
	"encoding/json"
	"testing"
)

// 原生币回调：转账信息在 native_token_transfers，token_transfers 不出现
func TestUnmarshalNativeTokenTransfer(t *testing.T) {
	body := `{
  "version": 1,
  "wallet_id": "w1",
  "event": "INCOMING_CONFIRMED",
  "item": {
    "direction": "incoming",
    "address": "TAddr",
    "blockchain": "tron",
    "network": "nile",
    "transaction_hash": "tx1",
    "trans_type": "native",
    "block_height": 67586000,
    "block_timestamp": 1779246426,
    "status": "success",
    "fee": {"amount": "1.7663", "symbol": "TRX"},
    "native_token_transfers": {
      "from_address": "TFrom",
      "to_address": "TAddr",
      "amount": "10",
      "decimals": 6,
      "symbol": "TRX",
      "token_name": "TRON"
    }
  }
}`

	var ev Event
	if err := json.Unmarshal([]byte(body), &ev); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if ev.Version != 1 || ev.WalletID != "w1" || ev.Event != EventIncomingConfirmed {
		t.Fatalf("event = %+v", ev)
	}
	if ev.Item.TransType != TransTypeNative {
		t.Fatalf("trans_type = %q", ev.Item.TransType)
	}
	if ev.Item.TokenTransfers != nil {
		t.Fatalf("原生币回调不应解析出 token_transfers: %+v", ev.Item.TokenTransfers)
	}

	nt := ev.Item.NativeTokenTransfers
	if nt == nil {
		t.Fatalf("原生币回调必须解析出 native_token_transfers: %+v", ev.Item)
	}
	if nt.FromAddress != "TFrom" || nt.ToAddress != "TAddr" {
		t.Fatalf("地址 = %s -> %s", nt.FromAddress, nt.ToAddress)
	}
	if nt.Amount.String() != "10" || nt.Decimals.String() != "6" {
		t.Fatalf("金额/精度 = %s/%s", nt.Amount, nt.Decimals)
	}
	if nt.Symbol != "TRX" || nt.TokenName != "TRON" || nt.Contract != "" {
		t.Fatalf("native_token_transfers = %+v", nt)
	}
}

// 代币回调：转账信息在 token_transfers，native_token_transfers 不出现
func TestUnmarshalTokenTransfer(t *testing.T) {
	body := `{
  "version": 1,
  "wallet_id": "w1",
  "event": "OUTGOING_CONFIRMED",
  "item": {
    "direction": "outgoing",
    "address": "TFrom",
    "blockchain": "tron",
    "network": "nile",
    "transaction_hash": "tx1",
    "trans_type": "trc20",
    "block_height": 67586000,
    "block_timestamp": 1779246426,
    "status": "success",
    "fee": {"amount": "2.1350", "symbol": "TRX"},
    "token_transfers": {
      "from_address": "TFrom",
      "to_address": "TTo",
      "amount": "5",
      "contract": "TContract",
      "decimals": 6,
      "symbol": "USDT",
      "token_name": "Tether USD"
    }
  }
}`

	var ev Event
	if err := json.Unmarshal([]byte(body), &ev); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if ev.Event != EventOutgoingConfirmed {
		t.Fatalf("event = %q", ev.Event)
	}
	if ev.Item.TransType != TransTypeTRC20 {
		t.Fatalf("trans_type = %q", ev.Item.TransType)
	}
	if ev.Item.NativeTokenTransfers != nil {
		t.Fatalf("代币回调不应解析出 native_token_transfers: %+v", ev.Item.NativeTokenTransfers)
	}
	tt := ev.Item.TokenTransfers
	if tt == nil {
		t.Fatalf("代币回调必须解析出 token_transfers: %+v", ev.Item)
	}
	if tt.Contract != "TContract" || tt.Symbol != "USDT" || tt.Amount.String() != "5" {
		t.Fatalf("token_transfers = %+v", tt)
	}
	if ev.Item.Direction != DirectionOutgoing {
		t.Fatalf("direction = %v", ev.Item.Direction)
	}
}

// 事件值与 trans_type 必须与服务端回调契约一致
func TestEventTypeValues(t *testing.T) {
	wantEvents := map[EventType]string{
		EventIncomingConfirmed: "INCOMING_CONFIRMED",
		EventOutgoingConfirmed: "OUTGOING_CONFIRMED",
		EventOutgoingFailed:    "OUTGOING_FAILED",
		EventKeepalive:         "KEEPALIVE",
	}
	for got, wantValue := range wantEvents {
		if string(got) != wantValue {
			t.Fatalf("EventType = %q, 期望 %q", got, wantValue)
		}
	}

	wantTransTypes := map[string]string{
		TransTypeNative: "native",
		TransTypeERC20:  "erc20",
		TransTypeTRC20:  "trc20",
	}
	for got, wantValue := range wantTransTypes {
		if got != wantValue {
			t.Fatalf("TransType = %q, 期望 %q", got, wantValue)
		}
	}
}

// 两个转账字段互斥且省略空值：未出现的分支不能被序列化成空对象
func TestEventJSONOmitsEmptyTransferFields(t *testing.T) {
	ev := Event{
		WalletID: "w1",
		Event:    EventIncomingConfirmed,
		Item: EventItem{
			Direction:       DirectionIncoming,
			Address:         "TAddr",
			Blockchain:      BLockchainTron,
			Network:         NetworkNile,
			TransactionHash: "tx1",
			NativeTokenTransfers: &TokenTransfer{
				FromAddress: "TFrom",
				ToAddress:   "TAddr",
				Symbol:      "TRX",
				TokenName:   "TRON",
			},
		},
	}
	b, err := json.Marshal(ev)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if _, ok := got["version"]; ok {
		t.Fatalf("version=0 时不应输出: %s", b)
	}
	item, _ := got["item"].(map[string]any)
	if _, ok := item["token_transfers"]; ok {
		t.Fatalf("token_transfers 为空时不应输出: %s", b)
	}
	if _, ok := item["native_token_transfers"]; !ok {
		t.Fatalf("native_token_transfers 必须输出: %s", b)
	}
	if _, ok := item["extra"]; ok {
		t.Fatalf("extra 为空时不应输出: %s", b)
	}
}
