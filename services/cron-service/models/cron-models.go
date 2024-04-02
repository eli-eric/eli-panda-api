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

type SyncEliBeamlinesEmployeesResponse struct {
	Rows   []SyncEliBeamlinesEmployee `json:"radky"`
	Params SyncEliBeamlinesParams     `json:"parametry"`
}

type SyncEliBeamlinesEmployee struct {
	EVID_STAV            string `json:"EVID_STAV"`
	ID_KARTY             string `json:"ID_KARTY"`
	PRACOVNI_MISTO_KOD   string `json:"PRACOVNI_MISTO_KOD"`
	EMAIL                string `json:"EMAIL"`
	NADRIZENY_OS_CISLO   string `json:"NADRIZENY_OS_CISLO"`
	PRACOVNI_MISTO_NAZEV string `json:"PRACOVNI_MISTO_NAZEV"`
	NADRIZENY_JMENO      string `json:"NADRIZENY_JMENO"`
	OJ_KOD               string `json:"OJ_KOD"`
	PRIJMENI             string `json:"PRIJMENI"`
	JMENO                string `json:"JMENO"`
	PPV_ID               string `json:"PPV_ID"`
	PHONE5               string `json:"PHONE5"`
	PHONE4               string `json:"PHONE4"`
	PHONE3               string `json:"PHONE3"`
	OSOBNI_CISLO         string `json:"OSOBNI_CISLO"`
	LDAP_USER            string `json:"LDAP_USER"`
	PRACOVISTE           string `json:"PRACOVISTE"`
}

type SyncEliBeamlinesParams struct {
	Datum        string `json:"Datum"`
	RoleIds      string `json:"roleIds"`
	OrgIds       string `json:"orgIds"`
	IdUzivatelId string `json:"id_uzivatelId"`
	Filtr        string `json:"FILTR"`
	OsobniCislo  string `json:"osobniCislo"`
	UserId       string `json:"userId"`
	Uzivatel     string `json:"UZIVATEL"`
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
