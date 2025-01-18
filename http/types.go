package http

type Chain struct {
	ChainID string
	SolverAddress string
}

type ApiResponse struct {
	SrcChain string
	DstChain string
	TotalOrders int
	TotalFilled int
	TotalRevenueUSDC float32
	TotalOsmosisTxFeesUSD float32
	TotalArbitrumTxFeesUSD float32
	TotalArbitrumTxFeesETH float32
	TotalOsmosisTxFeesOSMO float32
	SuccessRate float32
}