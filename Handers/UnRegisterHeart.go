package Handers

import (
	"log"
	"monitor-server/Metrics"
	"strings"
	"time"
)

var agentHeartbeatTimes = make(map[string]time.Time)

// 反解析 label 字符串并获取 hostname 和 project
func parseHeartbeatLabel(metricLabel string) (string, string) {
	parts := strings.Split(metricLabel, "_")
	if len(parts) < 2 {
		log.Printf("标签 %s 无法解析，格式不正确", metricLabel)
		return "", ""
	}
	return parts[0], parts[1]
}

// 定期检查超时的心跳数据
func CheckHeartbeats() {

	currentTime := time.Now()

	// 遍历 agentHeartbeatTimes 以检查每个 agent 的心跳状态
	for metricLabel, lastHeartbeat := range agentHeartbeatTimes {
		// 反解析 metricLabel 获取 hostname 和 project
		hostname, project := parseHeartbeatLabel(metricLabel)

		// 如果超过 10 秒没有接收到心跳
		if currentTime.Sub(lastHeartbeat) > 10*time.Second {
			// 设置 IsActive 为 0，表示该 agent 已经不活跃
			Metrics.IsActiveMetric.WithLabelValues(hostname, project).Set(0)
		}
	}

}
