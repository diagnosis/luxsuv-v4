# Email Setup Testing Guide

## Step 1: Configure Gmail App Password

1. **Enable 2-Factor Authentication** on your Google account:
   - Go to [Google Account Settings](https://myaccount.google.com/)
   - Click "Security" in the left sidebar
   - Under "Signing in to Google", click "2-Step Verification"
   - Follow the setup process

2. **Generate App Password**:
   - Still in Security settings, click "2-Step Verification"
   - Scroll down to "App passwords"
   - Click "App passwords"
   - Select "Mail" from the dropdown
   - Click "Generate"
   - **Copy the 16-character password** (it will look like: `abcd efgh ijkl mnop`)

## Step 2: Update .env File

Edit your `.env` file and add your email configuration:

```bash
# Email Configuration
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=demirkansafa@gmail.com
SMTP_PASSWORD=abcd efgh ijkl mnop
SMTP_FROM=demirkansafa@gmail.com
```

**Important**: 
- Use the 16-character app password, NOT your regular Gmail password
- Remove any spaces from the app password when copying

## Step 3: Test Email Functionality

1. **Restart your server** after updating .env:
```bash
go run cmd/server/main.go
```

2. **Check the logs** - you should see:
```
INFO: Email service initialized
INFO: SMTP Host: smtp.gmail.com
INFO: SMTP Port: 587
INFO: SMTP From: demirkansafa@gmail.com
INFO: Email service enabled
```

3. **Test password reset**:
```bash
curl -X POST http://localhost:8080/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{"email": "demirkansafa@gmail.com"}'
```

4. **Check your email inbox** for the password reset email!

## Troubleshooting

### If you see "Email service disabled":
- Make sure all SMTP fields are filled in .env
- Restart the server after updating .env

### If you get SMTP authentication errors:
- Double-check the app password (16 characters)
- Make sure 2FA is enabled on your Google account
- Try generating a new app password

### If emails aren't arriving:
- Check your spam/junk folder
- Verify the email address is correct
- Check the server logs for detailed error messages

### Alternative Email Providers:

**Outlook/Hotmail:**
```bash
SMTP_HOST=smtp-mail.outlook.com
SMTP_PORT=587
SMTP_USERNAME=your-email@outlook.com
SMTP_PASSWORD=your-regular-password
SMTP_FROM=your-email@outlook.com
```

**Yahoo:**
```bash
SMTP_HOST=smtp.mail.yahoo.com
SMTP_PORT=587
SMTP_USERNAME=your-email@yahoo.com
SMTP_PASSWORD=your-app-password
SMTP_FROM=your-email@yahoo.com
```