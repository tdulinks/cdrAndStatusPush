package models

// CDR 通话记录数据结构
type CDR struct {
	AccountID          string `json:"accountId"`
	CallID             string `json:"callId"`
	ServiceType        int    `json:"serviceType"`
	SubServiceType     int    `json:"subServiceType"`
	NumberPoolNo       string `json:"numberPoolNo"`
	Caller             string `json:"caller"`
	CallerCountryISO   string `json:"callerCountryIsoCode"`
	CallerProvinceCode string `json:"callerProvinceCode"`
	CallerCityCode     string `json:"callerCityCode"`
	Callee             string `json:"callee"`
	CalleeCountryISO   string `json:"calleeCountryIsoCode"`
	CalleeProvinceCode string `json:"calleeProvinceCode"`
	CalleeCityCode     string `json:"calleeCityCode"`
	BeginCallTime      int64  `json:"beginCallTime"`
	RingTime           int64  `json:"ringTime"`
	StartTime          int64  `json:"startTime"`
	EndTime            int64  `json:"endTime"`
	ReleaseType        int    `json:"releaseType"`
	CallDuration       int    `json:"callDuration"`
	CallResult         int    `json:"callResult"`
	AudioRecordFlag    int    `json:"audioRecordFlag"`
	CDRCreateTime      int64  `json:"cdrCreateTime"`
	SubscriptionID     string `json:"subscriptionId"`
	PhoneNoX           string `json:"phoneNoX"`
	PhoneNoA           string `json:"phoneNoA"`
	PhoneNoB           string `json:"phoneNoB"`
	SecretCallType     int    `json:"secretCallType"`
	MessageType        int    `json:"messageType"`
	CallDisplayType    int    `json:"callDisplayType"`
	CDRType            int    `json:"cdrType"`
	UserData           string `json:"userData"`
}
