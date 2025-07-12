# MailerSend Setup Guide

## Step 1: MailerSend Account Setup

1. **Sign up for MailerSend**:
   - Go to [MailerSend](https://www.mailersend.com/)
   - Create a free account (includes 12,000 emails/month)

2. **Get your API Key**:
   - Go to your MailerSend dashboard
   - Navigate to "API Tokens" in the left sidebar
   - Create a new token with "Full Access" permissions
   - Copy your API key (starts with `mlsn.`)

3. **Add and verify your domain**:
   - Go to "Domains" in the dashboard
   - Add your domain (or use a subdomain like `mail.yourdomain.com`)
   - Follow DNS verification steps
   - **For testing**: You can use MailerSend's sandbox domain initially

## Step 2: Update .env File

```bash
# Email Configuration (MailerSend)
MAILERSEND_API_KEY=mlsn.b5c409592e2a906c4d8e0a1c3564f9f66bd9d0a460df4bb2fd82da3d8edebbfa
MAILERSEND_FROM_EMAIL=noreply@yourdomain.com
MAILERSEND_FROM_NAME=LuxSUV Support
```

**Important Notes**:
- Use your actual API key from MailerSend dashboard
- The FROM_EMAIL must be from a verified domain
- For testing, you can use MailerSend's sandbox domain

## Step 3: Test Email Functionality

1. **Restart your server**:
```bash
go run cmd/server/main.go
```

2. **Check the logs** - you should see:
```
INFO: Email service initialized
INFO: MailerSend From Email: noreply@yourdomain.com
INFO: MailerSend From Name: LuxSUV Support
INFO: Email service enabled
```

3. **Test password reset**:
```bash
curl -X POST http://localhost:8080/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{"email": "demirkansafa@gmail.com"}'
```

4. **Check your email inbox** for the beautifully designed password reset email!

## Step 4: Test User Registration (Welcome Email)

```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser2",
    "email": "demirkansafa@gmail.com",
    "password": "password123",
    "role": "rider"
  }'
```

## What You'll Get

### 🎨 **Beautiful Email Templates**:
- **Password Reset**: Professional design with gradient headers, security notices, and clear call-to-action
- **Welcome Email**: Engaging onboarding email with feature highlights and modern styling

### 📊 **MailerSend Benefits**:
- **Reliable Delivery**: Better inbox placement than SMTP
- **Analytics**: Track opens, clicks, bounces, and more
- **Templates**: Visual template editor (optional)
- **Webhooks**: Real-time delivery notifications
- **Suppression Lists**: Automatic bounce/complaint handling

## Troubleshooting

### If you see "Email service disabled":
- Make sure `MAILERSEND_API_KEY` and `MAILERSEND_FROM_EMAIL` are set in .env
- Restart the server after updating .env

### If emails aren't sending:
- Verify your API key is correct and has proper permissions
- Check that your FROM_EMAIL domain is verified in MailerSend
- Look at server logs for detailed error messages
- Check MailerSend dashboard for delivery status

### For Testing Without Domain:
- Use MailerSend's sandbox mode
- Or set up a free subdomain and verify it

### Rate Limits:
- Free plan: 12,000 emails/month
- Rate limit: 1 email per second
- Upgrade plans available for higher volumes

## Production Checklist

- [ ] Verify your sending domain in MailerSend
- [ ] Set up proper DNS records (SPF, DKIM, DMARC)
- [ ] Configure webhooks for delivery tracking
- [ ] Set up suppression lists
- [ ] Monitor your sender reputation
- [ ] Update frontend URLs in email templates

Your emails will now be delivered reliably with professional styling and full tracking capabilities!