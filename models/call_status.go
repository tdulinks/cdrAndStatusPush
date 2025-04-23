package models

// CallStatus 呼叫状态数据结构
type CallStatus struct {
	AccountID      string `json:"accountId"`
	CallID         string `json:"callId"`
	ServiceType    int    `json:"serviceType"`
	Caller         string `json:"caller"`
	Callee         string `json:"callee"`
	EventTime      string  `json:"eventTime"`
	EventType      int    `json:"eventType"`
	AllEventType   []int  `json:"allEventType"`
	MessageType    int    `json:"messageType"`
	PhoneNoX       string `json:"phoneNoX,omitempty"`
	PhoneNoA       string `json:"phoneNoA,omitempty"`
	PhoneNoB       string `json:"phoneNoB,omitempty"`
	Party          int    `json:"party"`
	SubscriptionID string `json:"subscriptionId"`
	UserData       string `json:"userData"`
}

// CallEventType 呼叫事件类型
const (
	EventTypeCalling  = 1 // 呼叫中
	EventTypeRinging  = 2 // 振铃中
	EventTypeAnswered = 3 // 已接听
	EventTypeEnded    = 4 // 已结束
)
