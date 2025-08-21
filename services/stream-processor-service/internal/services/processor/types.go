package processor

type WhitelistEntry struct {
	Address    string `json:"address"`
	Blockchain string `json:"blockchain"`
}

type InfuraSubscriptionMessage struct {
	Params struct {
		Subscription string `json:"subscription"`
		Result       string `json:"result"`
	} `json:"params"`
}

type JSONRPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      int           `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type TransactionDetailsResponse struct {
	Result *struct {
		From string `json:"from"`
		To   string `json:"to"`
	} `json:"result"`
}
