# AWS SES Mock API (Go + Gin + SQLite)

This project is a mock implementation of AWS SES (Simple Email Service) built using **Go, Gin, and SQLite**. It lets you test email sending without actually sending real emails, making it useful for development and debugging. It also includes email blocking, rate limiting, and logging without storing actual email content.

##  Features
- Mimics AWS SES API behavior
- Blocks emails based on policies and cooldown periods
- Tracks email statistics without storing content
- Uses `.env` for easy configuration
- Lightweight and simple, using SQLite

---

##  Prerequisites
Make sure you have these installed:
- **Go**
- **SQLite**
- **Git**

---


### Clone the Repository
```bash
git clone https://github.com/yourusername/aws-ses-mock.git
cd aws-ses-mock
```

### Install Dependencies
```bash
go mod tidy
```

### Configure Environment Variables
Create a `.env` file in the project root with these settings:
```ini
LOG_LEVEL=debug
LOG_FILE=server.log

EMAIL_LIMIT_PER_HOUR=5
EMAIL_COOLDOWN=1

PORT=8080
ENV=dev
```



### Run Server
```bash
go run main.go
```



Your server will be running at **http://localhost:8080**.

---

## API Endpoints

###  **Send Email**
- **Endpoint:** `POST /send-email`
- **Request Body:**
```json
{
  "from": "sender@example.com",
  "to": ["recipient@example.com"],
  "subject": "Test Email",
  "body": "This is a test email."
}
```
- **Response:**
```json
{
  "message": "Email processed successfully"
}
```

---

### **Block an Email**
- **Endpoint:** `POST /block-email`
- **Request Body:**
```json
{
  "email": "blocked@example.com",
  "reason": "Spam detected"
}
```
- **Response:**
```json
{
  "message": "Email blocked successfully"
}
```

---

### **Unblock an Email**
- **Endpoint:** `DELETE /unblock-email/{email}`
- **Response:**
```json
{
  "message": "Email unblocked successfully"
}
```

---

### **Get Blocked Emails**
- **Endpoint:** `GET /blocked-emails`
- **Response:**
```json
[
  {
    "email_address": "blocked@example.com",
    "is_blocked": true,
    "self_cooldown_date": "2025-02-10T12:00:00Z",
    "self_cooldown_days": 3
  }
]
```

---

### **View Email Stats**
- **Endpoint:** `GET /stats`
- **Response:**
```json
{
  "total_emails_sent": 10,
  "total_emails_blocked": 2,
  "total_rate_limited": 1
}
```

---

## Running with Docker

### Build Docker Image
```bash
docker build -t aws-ses-mock .
```


### Run Docker Container
```bash
docker run -p 8080:8080 --env-file .env aws-ses-mock
```


---

## Notes
- **This is for testing only, no actual emails are sent.**
- **Configure `EMAIL_LIMIT_PER_HOUR` and `EMAIL_COOLDOWN` as needed.**
- **You can inspect stored data using an SQLite browser or CLI.**


---

## Contributing
If you have ideas or improvements, feel free to fork and submit a pull request!


---

## License
MIT License © 2025


---

