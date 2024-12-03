package Metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	//ssl
	IsActiveMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "is_active", // SSL 证书剩余天数
			Help: "agnet状态是否存活",
		},
		[]string{"hostName", "project"}, // 标签
	)
)