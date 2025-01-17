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
	SuccessRate float32
}