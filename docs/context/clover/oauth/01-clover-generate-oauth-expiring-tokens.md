# Generate OAuth expiring (access and refresh) token

> **Regions:** North America | Europe | Latin America | Asia Pacific

All REST API endpoints require an OAuth-generated `access_token` with specific permissions. Use the v2/OAuth flow to create an expiring authentication token, which includes an `access_token` and a `refresh_token` pair.

---

## Prerequisites

1. [Create a global developer account](https://docs.clover.com/docs/gdp-create-global-developer-account).
2. [Manage test merchant accounts and information](https://docs.clover.com/docs/gdp-manage-test-merchants-accounts).
3. [Create your app](https://docs.clover.com/docs/gdp-create-new-app) in the sandbox environment.
4. Configure [settings](https://docs.clover.com/dev/docs/gdp-manage-app-settings) and [permissions](https://docs.clover.com/dev/docs/gdp-set-app-permissions) that your app requires to access Clover merchant data.
5. Set the **Alternate Launch Path** — Required when the app OAuth is initiated from the left navigation menu on the Merchant Dashboard or directly from the Clover App Market. See [Set app link (URL) and CORS domain](https://docs.clover.com/dev/docs/using-cors).

---

## Steps

The Clover OAuth flow starts when the merchant selects your app directly from the [Clover App Market](https://www.clover.com/appmarket/) or from the left navigation on the Merchant Dashboard (**More Tools > Clover App Market**). Clover redirects the merchant to your app with the `merchantId` included in the Redirect URI as a query parameter. From there, your app must call the `/oauth/v2/authorize` endpoint to initiate the v2/OAuth flow and get an `access_token` and `refresh_token` pair.

> If a merchant accesses the app from your website instead of installing it from the Clover App Market, your app needs to redirect to the `/oauth/v2/authorize` endpoint.

To generate an expiring access and refresh token pair:

1. Log in to the [Global Developer Dashboard](https://www.clover.com/global-developer-home).
2. Navigate to the Merchant Dashboard for your test merchant.
3. From the left navigation menu, click **More**, and then select your app on the Clover App Market page.
4. Click **Connect** to install your app for the test merchant.

From here, the flow proceeds in two sub-steps:

**a) Merchant authorization**

Clover redirects the merchant to the location specified in the **Alternate Launch Path** field, and the app calls `/oauth/v2/authorize` with the authorization `code` as a query parameter.

OAuth callback format:

```
https://www.example.com/oauth_callback?code={AUTHORIZATION_CODE}&merchant_id={MERCHANT_ID}
```

**b) Token exchange**

Your app makes a POST request with the `client_id`, `client_secret`, and `code` to `/oauth/v2/token`. The response provides an `access_token` and `refresh_token` pair.

---

## Request and response examples

### Expiring OAuth token — High-trust app

**Request:**

```bash
curl --request POST \
  --url 'https://apisandbox.dev.clover.com/oauth/v2/token' \
  --header 'content-type: application/json' \
  --data '{
      "client_id": "{APP_ID}",
      "client_secret": "{APP_SECRET}",
      "code": "{AUTHORIZATION_CODE}"
  }'
```

**Response:**

```json
{
    "access_token": "{ACCESS_TOKEN}",
    "access_token_expiration": 1677875430,
    "refresh_token": "{REFRESH_TOKEN}",
    "refresh_token_expiration": 1709497830
}
```

---

### Expiring OAuth token — Low-trust app (PKCE)

**Request:**

```bash
curl --request POST \
  --url 'https://apisandbox.dev.clover.com/oauth/v2/token' \
  --header 'content-type: application/json' \
  --data '{
      "client_id": "{APP_ID}",
      "code": "{AUTHORIZATION_CODE}",
      "code_verifier": "{CODE_VERIFIER}"
  }'
```

**Response:**

```json
{
    "access_token": "{ACCESS_TOKEN}",
    "access_token_expiration": 1677875430,
    "refresh_token": "{REFRESH_TOKEN}",
    "refresh_token_expiration": 1709497830
}
```

For more information, see:
- [High-trust apps — Auth code flow](https://docs.clover.com/dev/docs/high-trust-app-auth-flow)
- [Low-trust apps — Auth code flow with PKCE](https://docs.clover.com/dev/docs/oauth-flow-for-low-trust-apps-pkce)

---

## Generate a new OAuth expiring token with a refresh token

In the v2/OAuth flow, an expiring authentication token consisting of an `access_token` and `refresh_token` pair is generated. The `access_token` is **short-lived**, while the `refresh_token` lasts longer but also expires eventually.

To maintain authorization, your app must generate a new token pair before the current one expires. Send a POST request to the `/oauth/v2/refresh` endpoint with the existing `refresh_token` and `client_id` to generate a new `access_token` and `refresh_token` pair.

For more information, see [Use refresh token to generate new expiring token](https://docs.clover.com/dev/docs/refresh-access-tokens).

---

## Sandbox and production environment URLs

Clover sandbox and production environments use different URLs. The following table lists which URL to use for OAuth requests in each environment.

| Request path | Sandbox URL | Production — North America | Production — Europe | Production — Latin America |
|---|---|---|---|---|
| `/oauth/v2/authorize` | `sandbox.dev.clover.com` | `www.clover.com` | `www.eu.clover.com` | `www.la.clover.com` |
| `/oauth/v2/token` | `apisandbox.dev.clover.com` | `api.clover.com` | `api.eu.clover.com` | `api.la.clover.com` |
| `/oauth/v2/refresh` | `apisandbox.dev.clover.com` | `api.clover.com` | `api.eu.clover.com` | `api.la.clover.com` |
| `/oauth/token/migrate_v2` | `apisandbox.dev.clover.com` | `api.clover.com` | `api.eu.clover.com` | `api.la.clover.com` |

---

## Related topics

- [Authenticate with v2/OAuth flow](https://docs.clover.com/dev/docs/use-oauth)
- [Use refresh token to generate new expiring token](https://docs.clover.com/dev/docs/refresh-access-tokens)
- [Blog: Expiring OAuth Tokens — Securing Clover Merchant Data](https://medium.com/clover-platform-blog/expiring-oauth-tokens-securing-clover-merchant-data-16243b9c00cc)
- [Blog: Fiddling Through Digital Keys — Clover Auth Tokens and Ecommerce Keys](https://medium.com/clover-platform-blog/understanding-clover-auth-tokens-and-ecommerce-keys-4e048827afa2)

---

*Source: [https://docs.clover.com/dev/docs/generate-expiring-tokens-using-v2-oauth-flow](https://docs.clover.com/dev/docs/generate-expiring-tokens-using-v2-oauth-flow)*
