# Coordinator

A Coordinator Node for the OpnLaaS pipeline

This project is used to run the coordination utility for OpnLaaS. This utility will handle user accounts, bookings, and management of hosts and bookings.

Please fill out a `.env` file in the root of this folder. It must look something like this:

```yaml
# Server Setup
HOST=127.0.0.1
PORT=8090

# Database setup
DB_FILE=sqlite.db
DB_SALT=RandomGoofyString
DB_QUEUE_SIZE=256

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

## Host Management

Hosts can be added and managed from the admin panel. To add a host, it must meet the following requirements:

1. Host must have a BMC with an IPMI web interface bound to a static IPv4 address
2. Host must have support for modern Redfish API standards
3. Host must have a logon for the IPMI and Redfish that grants administrative permissions

When you create a host, it will reach out via the Redfish API to query the health and specs of the host. The spec queries will be checked daily, and the health will be queried hourly. These durations can be configured in the env file.

If host issues persist, the SMTP client will email admin users about the issues. Emails will only be sent out about new issues.

> Please snsure that your BMC is running the latest firmware. If you are using a dell machine, please update your iDRAC using https://dell.com/support.