# Voice Service Push System

[中文文档](README.md)

## MIT License

Copyright (c) 2025 Voice Service Push System

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

## Project Introduction

This system simulates carrier call records and status push functionality, providing standard push interfaces for third-party applications. The system includes two types of push notifications:

1. Call Detail Record (CDR) Push - Pushes complete call record data after a call ends, similar to carrier's detailed billing push
2. Call Status Push - Pushes real-time status changes during calls, including dialing, ringing, answering, and hanging up states

This system can be used for:

- Simulating carrier call record push testing
- Developing call management systems
- Integration with call center systems
- Interfacing with third-party telephony systems

## Technical Specifications

### Interface Information

- Protocol: HTTP/HTTPS
- Request Method: POST
- Data Format: JSON
- Character Encoding: UTF-8
- Authentication: No signature required
- Success Response: HTTP 2xx

### Service Types

- 100: Voice SIP Service (simulating carrier fixed/mobile services)
- 200: Privacy Number Service (simulating carrier intermediate number service)

### Retry Mechanism

The system has a built-in retry mechanism as follows:

| Retry Count | Interval   | Description                         |
| ----------- | ---------- | ----------------------------------- |
| 1           | Immediate  | Immediate retry after first failure |
| 2           | 5 seconds  | Second retry                        |
| 3           | 30 seconds | Third retry                         |
| 4           | 5 minutes  | Fourth retry                        |
| 5           | 30 minutes | Final retry                         |

## Usage Notes

1. Privacy number related fields are only valid when serviceType=200
2. Pay attention to millisecond vs. second level timestamps
3. userData field supports up to 2048 characters
4. Recommended interface response time within 3 seconds

## System Installation and Configuration

### Environment Requirements

- Go 1.16 or above
- Configuration file: config/config.yaml

### Installation Steps

1. Clone the repository
2. Enter project directory: `cd cdr`
3. Install dependencies: `go mod download`
4. Modify configuration file config/config.yaml

## Service Start and Stop

### Starting Services

1. CDR Push Service:
   ```bash
   go run cmd/cdr/main.go
   ```
2. Status Push Service:
   ```bash
   go run cmd/status/main.go
   ```

### Stopping Services

Use Ctrl+C to terminate service processes

## Interface Call Examples

### CDR Push Interface

```json
{
  "callId": "202312010001",
  "serviceType": 100,
  "callerNumber": "1380000001",
  "calleeNumber": "02188888888",
  "startTime": 1701388800000,
  "endTime": 1701389100000,
  "duration": 300,
  "userData": "test-data"
}
```

### Status Push Interface

```json
{
  "callId": "202312010001",
  "serviceType": 100,
  "status": "ANSWERED",
  "timestamp": 1701388800000
}
```

## Common Issues

1. **How to handle push failures?**

   - Check if network connection is normal
   - Confirm if target server is accessible
   - Check log files for specific errors

2. **How to modify push retry configuration?**

   - Adjust retry count and intervals in configuration file
   - Restart service to apply configuration

3. **How to view push logs?**
   - Log files are located in logs directory
   - Named by date for easy query

## Documentation

For detailed interface specifications, please refer to [Voice Service Push System Development Documentation](Voice_Service_Push_System_Development_Doc.en.md)
