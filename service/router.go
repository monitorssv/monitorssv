package service

import "github.com/gin-gonic/gin"

func (ms *MonitorSSV) NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(Cors())
	r.Use(Recovery())
	r.Use(TraceLogger())
	r.Use(IpRageLimiter())

	r.GET("/api/status", ms.Status)
	r.GET("/api/dashboard", ms.Dashboard)
	r.GET("/api/operators", ms.GetOperators)
	r.GET("/api/clusters", ms.GetClusters)
	r.GET("/api/clusterDetails", ms.GetClusterDetails)
	r.GET("/api/validators", ms.GetValidators)
	r.GET("/api/events", ms.GetEvents)
	r.GET("/api/blocks", ms.GetBlocks)
	r.GET("/api/posData", ms.GetPosData)
	r.GET("/api/claim", ms.GetSSVReward)

	r.GET("/api/clusterMonitorInfo", ms.GetClusterMonitorInfo)
	r.GET("/api/clusterMonitorConfig", ms.GetClusterMonitorConfig)
	r.POST("/api/testAlarm", ms.TestAlarm)
	r.POST("/api/deleteClusterMonitorConfig", ms.DeleteMonitorConfig)
	r.POST("/api/saveClusterMonitorConfig", ms.SaveClusterMonitorConfig)

	return r
}
