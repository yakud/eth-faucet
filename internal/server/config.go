package server

type Config struct {
	network     string
	httpPort    int
	interval    int
	payout      float64
	proxyCount  int
	queueCap    int
	tokenSymbol string
}

func NewConfig(network string, httpPort, interval int, payout float64, proxyCount, queueCap int, tokenSymbol string) *Config {
	return &Config{
		network:     network,
		httpPort:    httpPort,
		interval:    interval,
		payout:      payout,
		proxyCount:  proxyCount,
		queueCap:    queueCap,
		tokenSymbol: tokenSymbol,
	}
}
