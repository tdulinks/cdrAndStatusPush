# Voice Service Push System Development Documentation

[中文文档](语音服务推送系统开发文档.md)

## Table of Contents

1. [Document Overview](#1-document-overview)
2. [Interface Basic Information](#2-interface-basic-information)
   - [General Specifications](#21-general-specifications)
   - [Retry Mechanism](#22-retry-mechanism)
3. [Call Record Push](#3-call-record-push)
   - [JSON Structure](#31-json-structure)
   - [Field Descriptions](#32-field-descriptions)
4. [Call Status Push](#4-call-status-push)
   - [JSON Structure](#41-json-structure)
   - [Field Descriptions](#42-field-descriptions)
5. [Enumeration Value Reference](#5-enumeration-value-reference)
6. [Notes](#6-notes)

## 1. Document Overview

This document describes two types of push interface specifications for the voice service system:

1. **Call Detail Record Push** (CDR): Pushes complete call record data after call ends
2. **Call Status Push**: Pushes real-time status changes during the call

## 2. Interface Basic Information

### 2.1 General Specifications

- **Protocol**: HTTP/HTTPS
- **Method**: POST
- **Data Format**: JSON
- **Encoding**: UTF-8
- **Authentication**: No signature verification
- **Success Determination**: HTTP status code 2xx

### 2.2 Retry Mechanism

| Retry Count | Interval   | Description                         |
| ----------- | ---------- | ----------------------------------- |
| 1           | Immediate  | Immediate retry after first failure |
| 2           | 5 seconds  | Second retry                        |
| 3           | 30 seconds | Third retry                         |
| 4           | 5 minutes  | Fourth retry                        |
| 5           | 30 minutes | Final retry                         |

## 3. Call Record Push

### 3.1 JSON Structure

```json
{
  "accountId": "1",
  "callId": "NM20010115474511150001012500007a8c",
  "serviceType": 100,
  "subServiceType": 101,
  "numberPoolNo": "NP160102116501911493",
  "caller": "13100001111",
  "callerCountryIsoCode": "CN",
  "callerProvinceCode": "GD",
  "callerCityCode": "SZ",
  "callee": "13100002222",
  "calleeCountryIsoCode": "CN",
  "calleeProvinceCode": "BJ",
  "calleeCityCode": "BJ",
  "beginCallTime": 1607588903000,
  "ringTime": 1607588903000,
  "startTime": 1607588904000,
  "endTime": 1607588923000,
  "releaseType": 1,
  "callDuration": 20,
  "callResult": 1,
  "audioRecordFlag": 1,
  "cdrCreateTime": 1607588924000,
  "subscriptionId": "sub123456",
  "phoneNoX": "13100003333",
  "phoneNoA": "13100001111",
  "phoneNoB": "13100002222",
  "secretCallType": 10,
  "messageType": 1,
  "callDisplayType": 1,
  "cdrType": 1,
  "userData": "{\"orderId\":\"123456\"}"
}
```

### 3.2 Field Descriptions

| Field         | Type    | Required | Description               |
| ------------- | ------- | -------- | ------------------------- |
| accountId     | String  | Yes      | Account ID                |
| callId        | String  | Yes      | Unique call ID            |
| serviceType   | Integer | Yes      | Service type              |
| beginCallTime | Long    | Yes      | Start time (milliseconds) |
| endTime       | Long    | Yes      | End time                  |
| callDuration  | Integer | Yes      | Call duration (seconds)   |
| ...           | ...     | ...      | ...                       |

## 4. Call Status Push

### 4.1 JSON Structure

```json
{
  "accountId": "10001",
  "callId": "NM20010115474511150001012500007a8c",
  "serviceType": 100,
  "caller": "13100001111",
  "callee": "13100002222",
  "eventTime": "1607588903",
  "allEventType": [1, 2, 3, 4],
  "eventType": 2,
  "messageType": 1,
  "phoneNoX": "13100003333",
  "phoneNoA": "13100001111",
  "phoneNoB": "13100002222",
  "party": 1,
  "subscriptionId": "sub123456",
  "userData": "{\"sessionId\":\"abcd1234\"}"
}
```

### 4.2 Field Descriptions

| Field     | Type    | Required | Description          |
| --------- | ------- | -------- | -------------------- |
| accountId | String  | Yes      | Account ID           |
| callId    | String  | Yes      | Unique call ID       |
| eventType | Integer | Yes      | Event type           |
| eventTime | String  | Yes      | Event time (seconds) |
| ...       | ...     | ...      | ...                  |

## 5. Enumeration Value Reference

### 5.1 Service Type (serviceType)

| Value | Description            |
| ----- | ---------------------- |
| 100   | Voice SIP Service      |
| 200   | Privacy Number Service |

### 5.2 Call Result (callResult)

| Value | Description       |
| ----- | ----------------- |
| 1     | Normal connection |
| 2     | Power off         |
| 3     | Suspended         |
| ...   | ...               |

## 6. Notes

1. Privacy number fields are only valid when serviceType=200
2. Pay attention to the difference between millisecond and second timestamps
3. userData maximum length is 2048 characters
4. Recommended interface processing time within 3 seconds
