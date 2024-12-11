package Metrics

import "github.com/prometheus/client_golang/prometheus"

var CustomRegistry = prometheus.NewRegistry()

// 预先注册静态指标
func init() {
	// 注册所有静态指标
	CustomRegistry.MustRegister(CpuPercentMetric)
	CustomRegistry.MustRegister(DiskTotalMetric)
	CustomRegistry.MustRegister(DiskUsedMetric)
	CustomRegistry.MustRegister(DiskFreeMetric)
	CustomRegistry.MustRegister(DiskUsedPercentMetric)
	CustomRegistry.MustRegister(MemoryTotalMetric)
	CustomRegistry.MustRegister(MemoryUsedMetric)
	CustomRegistry.MustRegister(MemoryFreeMetric)
	CustomRegistry.MustRegister(MemoryUsedPercentMetric)
	CustomRegistry.MustRegister(CpuLoad1Metric)
	CustomRegistry.MustRegister(CpuLoad5Metric)
	CustomRegistry.MustRegister(CpuLoad15Metric)
	CustomRegistry.MustRegister(NginxIsRunMetric)
	CustomRegistry.MustRegister(NginxReTotalMetric)
	CustomRegistry.MustRegister(NginxLoginUserCountMetric)
	CustomRegistry.MustRegister(NginxRawTotalMetric)
	CustomRegistry.MustRegister(NginxUdptotalMetric)
	CustomRegistry.MustRegister(NginxTcpTotalMetric)
	CustomRegistry.MustRegister(NginxTotalTcpMetric)
	CustomRegistry.MustRegister(NginxInetTotalMetric)
	CustomRegistry.MustRegister(NginxFragTotalMetric)
	CustomRegistry.MustRegister(NginxTcpEstabMetric)
	CustomRegistry.MustRegister(NginxTcpClosedMetric)
	CustomRegistry.MustRegister(NginxTcpOrphanedMetric)
	CustomRegistry.MustRegister(NginxTcpTimewaitMetric)
	CustomRegistry.MustRegister(SslDaysLeftMetric)
	CustomRegistry.MustRegister(ContainerCpuUsageMetric)
	CustomRegistry.MustRegister(ContainerMemoryUsageMetric)
	CustomRegistry.MustRegister(ContainerCpuLimitMetric)
	CustomRegistry.MustRegister(ContainerMemoryLimitMetric)
	CustomRegistry.MustRegister(ContainerRestartCountMetric)
	CustomRegistry.MustRegister(ContainerLastTerminationTimeMetric)
	CustomRegistry.MustRegister(IsActiveMetric)
	CustomRegistry.MustRegister(AgentVerisonMetric)
	CustomRegistry.MustRegister(ControllerReplicasMetric)
	CustomRegistry.MustRegister(ControllerReplicasAvailableMetric)
	CustomRegistry.MustRegister(ControllerReplicasUnavailableMetric)
}
