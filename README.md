# Coordinator

A Coordinator Node for the OpnLaaS pipeline

This project is used to run the coordination utility for OpnLaaS. This utility will handle user accounts, bookings, and management of hosts and bookings.

Please fill out a `.env` file in the root of this folder. It must look something like this:

```yaml
# Server Setup
HOST=127.0.0.1
PORT=8090

# Database setup
DB_SALT=GoofyRandomString

# SMTP Email setup
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-service-account@gmail.com
SMTP_PASSWORD=YOUR_SERVICE_PASSWORD

# Configuration
LAB_NAME=Local Lab
LAB_ORG=Local Domain
LAB_CONTACT=example@example.com
EMAIL_DOMAIN_WHITELIST=gmail.com,university.edu
```

If this is configured wrong, an error will be thrown and the server will not start.