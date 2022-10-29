# Go SSE Demo

A simple messaging app to transfer plain text data between peers.

### 🛠️ Setup
Start dev server using 
```bash
go run .
```

### 📃 API Endpoints
| Method | Path | Query params | Summary |
| ------------ | ------------ | ------------ | ------------ |
| GET | /ping | - | Simple ping request to check the server status |
| GET | /peers | - | Returns list of active peers |
| GET | /stream | id | Starts streaming of receiving messages by given id |
| GET | /send | to, msg | Sends given message to the receiver |
