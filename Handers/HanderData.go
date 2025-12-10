package Handers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"monitor-server/Metrics"
	"monitor-server/Modles"
	"os"
	"strconv"
	"strings"
	"time"
)

var projectNameDict map[string]string

// 处理 namespace，去掉 -v1 或 -v2
func cleanNamespace(namespace string) string {
	if strings.HasSuffix(namespace, "-v1") {
		return namespace[:len(namespace)-3] // 去掉 -v1
	} else if strings.HasSuffix(namespace, "-v2") {
		return namespace[:len(namespace)-3] // 去掉 -v2
	}
	return namespace
}

// 读取配置文件的函数
func LoadProjectDict(configPath string) error {
	// 打开配置文件
	file, err := os.Open(configPath)
	if err != nil {
		return fmt.Errorf("无法打开配置文件: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	// 读取文件内容
	content, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("无法读取配置文件内容: %v", err)
	}

	// 解析 JSON 数据到字典
	err = json.Unmarshal(content, &projectNameDict)
	if err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	return nil
}

// 处理 nginx 类型的数据
func HandleNginxData(data []interface{}, project string) {
	projectName := projectNameDict[project]
	if projectName == "" {
		// 如果字典中没有找到对应的中文名称，使用原值
		projectName = project
	}
	for _, item := range data {
		itemBytes, err := json.Marshal(item)
		if err != nil {
			log.Printf("无法序列化 nginx data 中的元素: %v", item)
			continue
		}

		var nginxData Modles.NginxSource
		if err := json.Unmarshal(itemBytes, &nginxData); err != nil {
			log.Printf("解析 Nginx 数据失败: %v", err)
			continue
		}

		// 更新 Nginx 指标并打印日志
		Metrics.NginxIsRunMetric.WithLabelValues(nginxData.HostName, projectName).Set(float64(nginxData.IsRun))
		Metrics.NginxReTotalMetric.WithLabelValues(nginxData.HostName, projectName).Set(float64(nginxData.ReTotal))
		Metrics.NginxLoginUserCountMetric.WithLabelValues(nginxData.HostName, projectName).Set(float64(nginxData.LoginUserCount))
		Metrics.NginxRawTotalMetric.WithLabelValues(nginxData.HostName, projectName).Set(float64(nginxData.RawTotal))
		Metrics.NginxUdptotalMetric.WithLabelValues(nginxData.HostName, projectName).Set(float64(nginxData.Udptotal))
		Metrics.NginxTcpTotalMetric.WithLabelValues(nginxData.HostName, projectName).Set(float64(nginxData.TcpTotal))
		Metrics.NginxTotalTcpMetric.WithLabelValues(nginxData.HostName, projectName).Set(float64(nginxData.TotalTcp))
		Metrics.NginxInetTotalMetric.WithLabelValues(nginxData.HostName, projectName).Set(float64(nginxData.InetTotal))
		Metrics.NginxFragTotalMetric.WithLabelValues(nginxData.HostName, projectName).Set(float64(nginxData.FragTotal))
		Metrics.NginxTcpEstabMetric.WithLabelValues(nginxData.HostName, projectName).Set(float64(nginxData.TcpEstab))
		Metrics.NginxTcpClosedMetric.WithLabelValues(nginxData.HostName, projectName).Set(float64(nginxData.TcpClosed))
		Metrics.NginxTcpOrphanedMetric.WithLabelValues(nginxData.HostName, projectName).Set(float64(nginxData.TcpOrphaned))
		Metrics.NginxTcpTimewaitMetric.WithLabelValues(nginxData.HostName, projectName).Set(float64(nginxData.TcpTimewait))
		UpdateNginxMetricWithTimestamp(fmt.Sprintf("%s_%s", nginxData.HostName, projectName))
	}
}

// 处理硬件相关的数据
func HandleHardData(data []interface{}, project string) {
	projectName := projectNameDict[project]
	if projectName == "" {
		// 如果字典中没有找到对应的中文名称，使用原值
		projectName = project
	}
	for _, item := range data {
		itemBytes, err := json.Marshal(item)
		if err != nil {
			log.Printf("无法序列化 hard data 中的元素: %v", item)
			continue
		}

		var hardData Modles.HardSource
		if err := json.Unmarshal(itemBytes, &hardData); err != nil {
			log.Printf("解析硬件数据失败: %v", err)
			continue
		}
		// 将 `hardData` 转换为 JSON
		jsonData, err := json.Marshal(hardData)
		if err != nil {
			log.Printf("将硬件数据转换为 JSON 失败: %v", err)
		} else {
			log.Println(string(jsonData)) // 打印 JSON 格式的数据
		}
		// 更新硬件相关指标
		Metrics.CpuPercentMetric.WithLabelValues(hardData.HostName, projectName, hardData.CPUModel, hardData.OSVersion, hardData.KernelVersion).Set(hardData.CPUPercent)
		Metrics.DiskTotalMetric.WithLabelValues(hardData.HostName, projectName, hardData.CPUModel, hardData.OSVersion, hardData.KernelVersion).Set(hardData.DiskTotal)
		Metrics.DiskUsedMetric.WithLabelValues(hardData.HostName, projectName, hardData.CPUModel, hardData.OSVersion, hardData.KernelVersion).Set(hardData.DiskUsed)
		Metrics.DiskFreeMetric.WithLabelValues(hardData.HostName, projectName, hardData.CPUModel, hardData.OSVersion, hardData.KernelVersion).Set(hardData.DiskFree)
		Metrics.DiskUsedPercentMetric.WithLabelValues(hardData.HostName, projectName, hardData.CPUModel, hardData.OSVersion, hardData.KernelVersion).Set(hardData.DiskUsedPercent)
		Metrics.MemoryTotalMetric.WithLabelValues(hardData.HostName, projectName, hardData.CPUModel, hardData.OSVersion, hardData.KernelVersion).Set(hardData.MemoryTotal)
		Metrics.MemoryUsedMetric.WithLabelValues(hardData.HostName, projectName, hardData.CPUModel, hardData.OSVersion, hardData.KernelVersion).Set(hardData.MemoryUsed)
		Metrics.MemoryFreeMetric.WithLabelValues(hardData.HostName, projectName, hardData.CPUModel, hardData.OSVersion, hardData.KernelVersion).Set(hardData.MemoryFree)
		Metrics.MemoryUsedPercentMetric.WithLabelValues(hardData.HostName, projectName, hardData.CPUModel, hardData.OSVersion, hardData.KernelVersion).Set(hardData.MemoryUsedPercent)
		Metrics.CpuLoad1Metric.WithLabelValues(hardData.HostName, projectName, hardData.CPUModel, hardData.OSVersion, hardData.KernelVersion).Set(hardData.CPULoad1)
		Metrics.CpuLoad5Metric.WithLabelValues(hardData.HostName, projectName, hardData.CPUModel, hardData.OSVersion, hardData.KernelVersion).Set(hardData.CPULoad5)
		Metrics.CpuLoad15Metric.WithLabelValues(hardData.HostName, projectName, hardData.CPUModel, hardData.OSVersion, hardData.KernelVersion).Set(hardData.CPULoad15)
		Metrics.CpuTotalMetric.WithLabelValues(hardData.HostName, projectName, hardData.CPUModel, hardData.OSVersion, hardData.KernelVersion).Set(hardData.CPUCount)

		UpdateHardMetricWithTimestamp(fmt.Sprintf("%s_%s_%s_%s", hardData.HostName, projectName, hardData.CPUModel, hardData.OSVersion, hardData.KernelVersion))

	}
}

// 处理SSl证书数据
func HandleSSLData(data []interface{}, project string) {
	projectName := projectNameDict[project]
	if projectName == "" {
		// 如果字典中没有找到对应的中文名称，使用原值
		projectName = project
	}
	for _, item := range data {
		itemBytes, err := json.Marshal(item)
		if err != nil {
			log.Printf("无法序列化 ssl data 中的元素: %v", item)
			continue
		}

		var sslData Modles.SslSource
		if err := json.Unmarshal(itemBytes, &sslData); err != nil {
			log.Printf("解析 SSL 数据失败: %v", err)
			continue
		}

		// 如果 comment 为空，则设置为 "未备注"
		if sslData.Comment == "" {
			sslData.Comment = "未备注"
		}

		resolve := "false" // 默认 false
		if sslData.Resolve {
			resolve = "true"
		}

		// 更新 SSL 指标并打印日志，添加 project 标签
		metricLabel := fmt.Sprintf("%s_%s_%s_%s_%s", sslData.Domain, sslData.Comment, sslData.Status, resolve, projectName)
		Metrics.SslDaysLeftMetric.WithLabelValues(sslData.Domain, sslData.Comment, sslData.Status, resolve, projectName).Set(float64(sslData.DaysLeft))

		// 存储时间戳
		sslTimestamp.Store(metricLabel, time.Now())

	}
}

// 处理容器资源数据
func HandleContainerResourceData(data []interface{}, project string) {
	// 将项目名称映射为中文名称（假设有一个字典）
	projectName := projectNameDict[project]
	if projectName == "" {
		// 如果字典中没有找到对应的中文名称，使用原值
		projectName = project
	}

	for _, item := range data {
		// 将每个资源项转化为 JSON 字节
		itemBytes, err := json.Marshal(item)
		if err != nil {
			log.Printf("无法序列化容器资源数据中的元素: %v", item)
			continue
		}

		// 将字节数据反序列化为 ContainerResource 类型
		var containerResource Modles.ContainerResource
		if err := json.Unmarshal(itemBytes, &containerResource); err != nil {
			log.Printf("解析容器资源数据失败: %v", err)
			continue
		}
		containerNamespace := cleanNamespace(containerResource.Namespace)

		Metrics.ContainerCpuUsageMetric.WithLabelValues(containerNamespace, containerResource.PodName, containerResource.Container, containerResource.ControllerName, projectName).Set(containerResource.UseCpu)
		Metrics.ContainerMemoryUsageMetric.WithLabelValues(containerNamespace, containerResource.PodName, containerResource.Container, containerResource.ControllerName, projectName).Set(float64(containerResource.UseMemory))
		Metrics.ContainerCpuLimitMetric.WithLabelValues(containerNamespace, containerResource.PodName, containerResource.Container, containerResource.ControllerName, projectName).Set(containerResource.LimitCpu)
		Metrics.ContainerMemoryLimitMetric.WithLabelValues(containerNamespace, containerResource.PodName, containerResource.Container, containerResource.ControllerName, projectName).Set(float64(containerResource.LimitMemory))
		Metrics.ContainerRestartCountMetric.WithLabelValues(containerNamespace, containerResource.PodName, containerResource.Container, containerResource.ControllerName, projectName).Set(float64(containerResource.RestartCount))
		Metrics.ContainerLastTerminationTimeMetric.WithLabelValues(containerNamespace, containerResource.PodName, containerResource.Container, containerResource.ControllerName, projectName).Set(float64(containerResource.LastTerminationTime))

		UpdateContainerMetricWithTimestamp(fmt.Sprintf("%s_%s_%s_%s_%s", containerNamespace, containerResource.PodName, containerResource.Container, containerResource.ControllerName, projectName))
	}
}
func HandleTrafficSwitchingData(data []interface{}, project string) {
	projectName := projectNameDict[project]
	if projectName == "" {
		projectName = project
	}

	for _, item := range data {
		itemBytes, err := json.Marshal(item)
		if err != nil {
			log.Printf("无法序列化 traffic switching data 中的元素: %v", item)
			continue
		}

		var ts Modles.TrafficSwitchingSource
		if err := json.Unmarshal(itemBytes, &ts); err != nil {
			log.Printf("解析 traffic switching 数据失败: %v", err)
			continue
		}

		service := ts.Service

		// 解析 success_rate（兼容字符串和数字）
		var successRate float64
		switch v := ts.TotalSuccessRate.(type) {
		case float64:
			successRate = v
		case string:
			// 去掉 % 后解析，如 "85.50%" -> 0.855
			s := strings.TrimSuffix(v, "%")
			if parsed, err := strconv.ParseFloat(s, 64); err == nil {
				successRate = parsed / 100.0
			}
		}

		// 累计统计
		Metrics.TrafficSwitchingTotalRequests.WithLabelValues(service, projectName).Set(ts.TotalRequests)
		Metrics.TrafficSwitchingTotalSuccess.WithLabelValues(service, projectName).Set(ts.TotalSuccess)
		Metrics.TrafficSwitchingTotalErrors.WithLabelValues(service, projectName).Set(ts.TotalErrors)
		Metrics.TrafficSwitchingTotalSuccessRate.WithLabelValues(service, projectName).Set(successRate)

		// 今日统计
		Metrics.TrafficSwitchingTodayRequests.WithLabelValues(service, projectName).Set(ts.TodayRequests)
		Metrics.TrafficSwitchingTodaySuccess.WithLabelValues(service, projectName).Set(ts.TodaySuccess)
		Metrics.TrafficSwitchingTodayErrors.WithLabelValues(service, projectName).Set(ts.TodayErrors)
		Metrics.TrafficSwitchingTodayCanceled.WithLabelValues(service, projectName).Set(ts.TodayCanceled)
		Metrics.TrafficSwitchingTodayStatus2xx.WithLabelValues(service, projectName).Set(ts.TodayStatus2xx)
		Metrics.TrafficSwitchingTodayStatus3xx.WithLabelValues(service, projectName).Set(ts.TodayStatus3xx)
		Metrics.TrafficSwitchingTodayStatus4xx.WithLabelValues(service, projectName).Set(ts.TodayStatus4xx)
		Metrics.TrafficSwitchingTodayStatus5xx.WithLabelValues(service, projectName).Set(ts.TodayStatus5xx)

		// 实时统计
		Metrics.TrafficSwitchingRealtimeQPS.WithLabelValues(service, projectName).Set(ts.RealtimeQPS)
		Metrics.TrafficSwitchingRealtimeSuccessQPS.WithLabelValues(service, projectName).Set(ts.RealtimeSuccessQPS)
		Metrics.TrafficSwitchingRealtimeErrorQPS.WithLabelValues(service, projectName).Set(ts.RealtimeErrorQPS)
		Metrics.TrafficSwitchingRealtimeActiveConnections.WithLabelValues(service, projectName).Set(ts.RealtimeActiveConnections)
		Metrics.TrafficSwitchingRealtimeAvgLatencyMs.WithLabelValues(service, projectName).Set(ts.RealtimeAvgLatencyMs)
		Metrics.TrafficSwitchingRealtimeMaxLatencyMs.WithLabelValues(service, projectName).Set(ts.RealtimeMaxLatencyMs)

		// 错误类型
		Metrics.TrafficSwitchingErrorBackendError.WithLabelValues(service, projectName).Set(ts.ErrorBackendError)
		Metrics.TrafficSwitchingErrorBrokenPipe.WithLabelValues(service, projectName).Set(ts.ErrorBrokenPipe)
		Metrics.TrafficSwitchingErrorConnectionRefused.WithLabelValues(service, projectName).Set(ts.ErrorConnectionRefused)
		Metrics.TrafficSwitchingErrorConnectionReset.WithLabelValues(service, projectName).Set(ts.ErrorConnectionReset)
		Metrics.TrafficSwitchingErrorDNSError.WithLabelValues(service, projectName).Set(ts.ErrorDNSError)
		Metrics.TrafficSwitchingErrorEOF.WithLabelValues(service, projectName).Set(ts.ErrorEOF)
		Metrics.TrafficSwitchingErrorTimeout.WithLabelValues(service, projectName).Set(ts.ErrorTimeout)

		// 代理缓存
		Metrics.TrafficSwitchingProxyCacheSize.WithLabelValues(service, projectName).Set(ts.ProxyCacheSize)
		Metrics.TrafficSwitchingProxyMaxCacheSize.WithLabelValues(service, projectName).Set(ts.ProxyMaxCacheSize)

		// Runtime
		Metrics.TrafficSwitchingRuntimeGoroutines.WithLabelValues(service, projectName).Set(ts.RuntimeGoroutines)
		Metrics.TrafficSwitchingRuntimeMemoryMB.WithLabelValues(service, projectName).Set(ts.RuntimeMemoryMB)
		Metrics.TrafficSwitchingRuntimeCPUCores.WithLabelValues(service, projectName).Set(ts.RuntimeCPUCores)
		Metrics.TrafficSwitchingRuntimeGomaxprocs.WithLabelValues(service, projectName).Set(ts.RuntimeGomaxprocs)
		Metrics.TrafficSwitchingRuntimeGcCycles.WithLabelValues(service, projectName).Set(ts.RuntimeGcCycles)

		// Transport 配置
		Metrics.TrafficSwitchingTransportMaxConnsPerHost.WithLabelValues(service, projectName).Set(ts.TransportMaxConnsPerHost)
		Metrics.TrafficSwitchingTransportMaxIdleConns.WithLabelValues(service, projectName).Set(ts.TransportMaxIdleConns)
		Metrics.TrafficSwitchingTransportMaxIdleConnsPerHost.WithLabelValues(service, projectName).Set(ts.TransportMaxIdleConnsPerHost)

		// 上报时间戳
		Metrics.TrafficSwitchingTimestamp.WithLabelValues(service, projectName).Set(ts.Timestamp)

		// 更新心跳时间戳
		UpdateTrafficSwitchingTimestamp(fmt.Sprintf("%s_%s", service, projectName))
	}
}

// 更新心跳数据
func HandleHeartData(data []interface{}, project string) {
	projectName := projectNameDict[project]
	if projectName == "" {
		// 如果字典中没有找到对应的中文名称，使用原值
		projectName = project
	}
	for _, item := range data {
		itemBytes, err := json.Marshal(item)
		if err != nil {
			log.Printf("无法序列化 hard data 中的元素: %v", item)
			continue
		}

		var hardData Modles.HardSource
		if err := json.Unmarshal(itemBytes, &hardData); err != nil {
			log.Printf("解析硬件数据失败: %v", err)
			continue
		}
		for _, item := range data {
			itemBytes, err := json.Marshal(item)
			if err != nil {
				log.Printf("无法序列化 ssl data 中的元素: %v", item)
				continue
			}

			var heartData Modles.HeartSource
			if err := json.Unmarshal(itemBytes, &heartData); err != nil {
				log.Printf("解析 SSL 数据失败: %v", err)
				continue
			}

			// 记录时间戳
			// 更新 心跳 指标
			Metrics.IsActiveMetric.WithLabelValues(heartData.Hostname, projectName).Set(float64(heartData.IsActive))
			Metrics.AgentVerisonMetric.WithLabelValues(heartData.Hostname, projectName).Set(float64(heartData.Version))
			metricLabel := fmt.Sprintf("%s_%s", hardData.HostName, projectName)

			agentHeartbeatTimes.Store(metricLabel, time.Now())
		}
	}
}

// 更新控制器数据
func HandleControllertResourceData(data []interface{}, project string) {
	projectName := projectNameDict[project]
	if projectName == "" {
		// 如果字典中没有找到对应的中文名称，使用原值
		projectName = project
	}
	for _, item := range data {
		itemBytes, err := json.Marshal(item)
		if err != nil {
			log.Printf("无法序列化 hard data 中的元素: %v", item)
			continue
		}

		var controllerData Modles.ControllerResource
		if err := json.Unmarshal(itemBytes, &controllerData); err != nil {
			log.Printf("解析硬件数据失败: %v", err)
			continue
		}
		for _, item := range data {
			itemBytes, err := json.Marshal(item)
			if err != nil {
				log.Printf("无法序列化 ssl data 中的元素: %v", item)
				continue
			}

			var controllerData Modles.ControllerResource
			if err := json.Unmarshal(itemBytes, &controllerData); err != nil {
				log.Printf("解析 SSL 数据失败: %v", err)
				continue
			}

			containerNamespace := cleanNamespace(controllerData.Namespace)
			// 更新 心跳 指标
			Metrics.ControllerReplicasMetric.WithLabelValues(containerNamespace, controllerData.Container, controllerData.ControllerType, projectName).Set(float64(controllerData.Replicas))
			Metrics.ControllerReplicasAvailableMetric.WithLabelValues(containerNamespace, controllerData.Container, controllerData.ControllerType, projectName).Set(float64(controllerData.ReplicasAvailable))
			Metrics.ControllerReplicasUnavailableMetric.WithLabelValues(containerNamespace, controllerData.Container, controllerData.ControllerType, projectName).Set(float64(controllerData.ReplicasUnavailable))

			UpdateControllerMetricWithTimestamp(fmt.Sprintf("%s_%s_%s_%s", containerNamespace, controllerData.Container, controllerData.ControllerType, projectName))

		}
	}
}
