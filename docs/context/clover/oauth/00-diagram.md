# Clover OAuth Authorization Flow

## Participants

1. Developer App (Frontend)
2. Developer Backend
3. Clover UI
4. Clover Backend

---

# Authorization Flow

## Step 1 - Request Merchant Authorization

### From
Developer App

### To
Clover UI

### Action
Redirect merchant to Clover authorization endpoint.

### Endpoint
`/oauth/v2/authorize`

---

## Step 2 - User Logs In and Authorizes App

### Actor
Merchant/User

### System
Clover UI

### Action
User logs in and grants permissions to the app.

---

## Step 3 - Clover Generates Authorization Code

### From
Clover UI

### To
Clover Backend

### Action
Generate authorization code after successful authorization.

---

## Step 4 - Redirect Back With Authorization Code

### From
Clover Backend

### To
Developer App

### Action
Redirect user back to application with authorization code.

---

## Step 5 - Send Authorization Code to Backend

### From
Developer App

### To
Developer Backend

### Action
Frontend sends authorization code securely to backend.

---

# Access Token Exchange

## Step 6 - Request Access Token

### From
Developer Backend

### To
Clover Backend

### Action
Exchange authorization code for access token.

### Endpoint
`/oauth/v2/token`

### Payload
- authorization_code
- client_secret

---

## Step 7 - Validate Authorization Code

### Actor
Clover Backend

### Action
Validate:
- authorization code
- client secret

---

## Step 8 - Return Tokens

### From
Clover Backend

### To
Developer Backend

### Action
Return:
- access token
- refresh token

---

## Step 9 - Send Access Token to Frontend

### From
Developer Backend

### To
Developer App

### Action
Backend sends access token back to frontend application.

---

# Refresh Token Flow

## Step 10 - Refresh Access Token

### From
Developer Backend

### To
Clover Backend

### Action
Request new access token using refresh token.

### Endpoint
`/oauth/v2/refresh`

### Payload
- refresh_token

---

## Step 11 - Validate Refresh Token

### Actor
Clover Backend

### Action
Validate refresh token.

---

## Step 12 - Invalidate Previous Refresh Token

### Actor
Clover Backend

### Action
Invalidate old refresh token after issuing a new one.

---

## Step 13 - Return New Tokens

### From
Clover Backend

### To
Developer Backend

### Action
Return:
- new access token
- new refresh token

---

# High-Level Summary

```text
Developer App
    -> Redirect user to Clover authorization page
    -> Receive authorization code
    -> Send code to backend

Developer Backend
    -> Exchange authorization code for tokens
    -> Store/manage tokens
    -> Refresh tokens when needed

Clover Backend
    -> Validate authorization code
    -> Issue access and refresh tokens
    -> Validate refresh tokens
    -> Rotate refresh tokens