package models

type CronJobHistory struct {
	UID         string `json:"uid"`
	JobUID      string `json:"jobUid"`
	JobName     string `json:"jobName"`
	JobType     string `json:"jobType"`
	JobStatus   string `json:"jobStatus"`
	JobMessage  string `json:"jobMessage"`
	StartTime   int64  `json:"startTime"`
	EndTime     *int64 `json:"endTime"`
	TotalTimeMs int64  `json:"totalTimeMs"`
}

type CronJob struct {
	JobUID  string `json:"jobUid"`
	JobName string `json:"jobName"`
	JobType string `json:"jobType"`
}

const (
	CRON_JOB_STATUS_SUCCESS = "SUCCESS"
	CRON_JOB_STATUS_ERROR   = "ERROR"
	CRON_JOB_STATUS_RUNNING = "RUNNING"
	CRON_JOB_STATUS_PENDING = "PENDING"

	CRON_JOB_TYPE_ORDERS    = "ORDERS"
	CRON_JOB_TYPE_EMPLOYEES = "EMPLOYEES"
)

var CRON_JOB_SYNC_EXCHANGE_RATES_AND_PRICES_ELI_BM = CronJob{
	JobUID:  "orders-sync-exchange-rates-and-prices-eli-bm",
	JobName: "Sync exchange rates and prices for ELI - Beamlines",
	JobType: CRON_JOB_TYPE_ORDERS,
}

var CRON_JOB_SYNC_EMPLOYEES_FROM_HR_ELI_BM = CronJob{
	JobUID:  "sync-employees-from-hr-eli-bm",
	JobName: "Sync employees from HR - ELI - Beamlines",
	JobType: CRON_JOB_TYPE_EMPLOYEES,
}
