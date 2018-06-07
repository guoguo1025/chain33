package relay

import (
	"fmt"

	"gitlab.33.cn/chain33/chain33/types"
)

const (
	relayOrderSCAIH    = "relay-sellorder-scaih:"
	relayOrderCSAIH    = "relay-sellorder-csaih:"
	relayOrderASCIH    = "relay-sellorder-ascih:"
	relayOrderACSIH    = "relay-sellorder-acsih:"
	relayBuyOrderACSIH = "relay-buyorder-acsih:"
	orderIDPrefix      = "mavl-relay-orderid-"

	relayBTCHeaderHash       = "relay-btcheader-hash"
	relayBTCHeaderHeight     = "relay-btcheader-height"
	relayBTCHeaderHeightList = "relay-btcheader-height-list"
	relayBTCHeaderLastHeight = "relay-btcheader-last-height"
	relayBTCHeaderBaseHeight = "relay-btcheader-base-height"
	relayRcvBTCHighestHead   = "relay-rcv-btcheader-ht"
)

func calcBtcHeaderKeyHash(hash string) []byte {
	key := fmt.Sprintf(relayBTCHeaderHash+"%s", hash)
	return []byte(key)
}

func calcBtcHeaderKeyHeight(height int64) []byte {
	key := fmt.Sprintf(relayBTCHeaderHeight+"%d", height)
	return []byte(key)
}

func calcBtcHeaderKeyHeightList(height int64) []byte {
	key := fmt.Sprintf(relayBTCHeaderHeightList+"%d", height)
	return []byte(key)
}

func calcBtcHeaderKeyLastHeight() []byte {
	return []byte(relayBTCHeaderLastHeight)
}

func calcBtcHeaderKeyBaseHeight() []byte {
	return []byte(relayBTCHeaderBaseHeight)
}

func calcBtcHeightListKey() []byte {
	return []byte(relayBTCHeaderHeightList)
}

func calcOrderKeyStatus(order *types.RelayOrder, status int32) []byte {
	key := fmt.Sprintf(relayOrderSCAIH+"%d:%s:%s:%s:%d",
		status, order.Coin, order.CreaterAddr, order.Id, order.Height)
	return []byte(key)
}

func calcOrderKeyCoin(order *types.RelayOrder, status int32) []byte {
	key := fmt.Sprintf(relayOrderCSAIH+"%s:%d:%s:%s:%d",
		order.Coin, status, order.CreaterAddr, order.Id, order.Height)
	return []byte(key)
}

func calcOrderKeyAddrStatus(order *types.RelayOrder, status int32) []byte {
	key := fmt.Sprintf(relayOrderASCIH+"%s:%d:%s:%s:%d",
		order.CreaterAddr, status, order.Coin, order.Id, order.Height)
	return []byte(key)
}

func calcOrderKeyAddrCoin(order *types.RelayOrder, status int32) []byte {
	key := fmt.Sprintf(relayOrderACSIH+"%s:%s:%d:%s:%d",
		order.CreaterAddr, order.Coin, status, order.Id, order.Height)
	return []byte(key)
}

func calcOrderPrefixStatus(status int32) []byte {
	prefix := fmt.Sprintf(relayOrderSCAIH+"%d:", status)
	return []byte(prefix)
}

func calcOrderPrefixCoinStatus(coin string, status int32) []byte {
	prefix := fmt.Sprintf(relayOrderCSAIH+"%s:%d:", coin, status)
	return []byte(prefix)
}

func calcOrderPrefixAddrCoin(addr string, coin string) []byte {
	key := fmt.Sprintf(relayOrderACSIH+"%s:%s", addr, coin)
	return []byte(key)
}

func calcOrderPrefixAddr(addr string) []byte {
	return []byte(fmt.Sprintf(relayOrderACSIH+"%s", addr))
}

func calcAcceptKeyAddr(order *types.RelayOrder, status int32) []byte {
	if order.AcceptAddr != "" {
		return []byte(fmt.Sprintf(relayBuyOrderACSIH+"%s:%s:%d:%s:%d",
			order.AcceptAddr, order.Coin, status, order.Id, order.Height))
	}
	return nil

}

func calcAcceptPrefixAddr(addr string) []byte {
	return []byte(fmt.Sprintf(relayBuyOrderACSIH+"%s", addr))
}

func calcAcceptPrefixAddrCoin(addr, coin string) []byte {
	return []byte(fmt.Sprintf(relayBuyOrderACSIH+"%s:%s", addr, coin))
}

func calcRelayOrderID(hash string) string {
	return orderIDPrefix + hash
}
