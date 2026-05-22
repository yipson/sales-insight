# Use refresh token to generate new expiring token

> **Regions:** North America | Europe | Latin America

---

## Overview

In the v2/OAuth flow, an expiring authentication token consisting of an `access_token` and `refresh_token` pair is generated. The `access_token` is **short-lived**, while the `refresh_token` lasts longer but also expires eventually.

To maintain authorization, your app must generate a new token pair before the current one expires. Send a POST request to the `/oauth/v2/refresh` endpoint with the existing `refresh_token` and `client_id` to generate a new `access_token` and `refresh_token` pair. Your app needs to handle the refreshing of access tokens to allow merchants continuous access to the app.

---

## Key features of expiring tokens

- **Limitations on the number of refresh tokens** — Clover limits the number of active refresh tokens an app can have for each merchant. If an app exceeds the limit, prior tokens become invalid.
- **Limitations on refresh token usage** — The refresh token is for **single use** and becomes invalid immediately after a new `access_token` and `refresh_token` pair is generated using the `/oauth/v2/refresh` endpoint.
- **Dynamic expiration dates and length** — You should **not** hard-code access and refresh token expirations or lengths. Handle them dynamically in your app:
  - Token expiration date is based on the time the token is generated; tokens created over a period of time have different expiration dates. The format is **Unix timestamp**.
  - Token length is **not fixed**.

---

## Prerequisite

You need an `access_token` and `refresh_token` pair along with the `client_id` for which you initiated the OAuth flow. See [Generate OAuth expiring (access and refresh) token](https://docs.clover.com/dev/docs/generate-expiring-tokens-using-v2-oauth-flow).

1. Install the app and receive authorization information from the `/oauth/v2/authorize` endpoint.
2. Send a POST request to the `/oauth/v2/token` endpoint with the following parameters: `client_id`, `client_secret`, and `code`.

The response body displays the access and refresh token pair, along with their expiration dates in Unix timestamp format.

---

## Generate new access and refresh token pair

1. Send a POST request to the `/oauth/v2/refresh` endpoint.
2. Include the `client_id` and `refresh_token` from your app's initial token request.

### Request and response example — Generate new access and refresh tokens

**Request:**

```bash
curl --request POST \
  --url 'https://apisandbox.dev.clover.com/oauth/v2/refresh' \
  --header 'content-type: application/json' \
  --data '{
      "client_id": "{APP_ID}",
      "refresh_token": "{REFRESH_TOKEN}"
  }'
```

**Response:**

```json
{
    "access_token": "{NEW_ACCESS_TOKEN}",
    "access_token_expiration": 1709498373,
    "refresh_token": "{NEW_REFRESH_TOKEN}",
    "refresh_token_expiration": 1741034373
}
```

The response body displays the new access and refresh tokens along with their expiration dates in Unix timestamp format.

---

## Bypass refresh token creation in the OAuth flow

Clover limits the number of active refresh tokens that an app can have per merchant. In some scenarios a refresh token is not needed, for example:

- **Frontend apps** that use OAuth to authenticate users to their own apps often don't need a refresh token.
- **Apps that only need the access token** to verify that the user successfully authenticated with Clover, then use that token to get details from the API.

In such cases, generating a refresh token may cause the app to reach the limit unnecessarily. Use the `no_refresh_token` parameter on the `/oauth/v2/token` endpoint to bypass refresh token creation.

To bypass generating a refresh token:

1. Install the app and receive authorization information from the `/oauth/v2/authorize` endpoint.
2. Send a POST request to `/oauth/v2/token` with `client_id`, `client_secret`, `code`, and set `no_refresh_token` to `true`.

### Request and response example — No refresh token

**Request:**

```bash
curl --request POST \
  --url 'https://apisandbox.dev.clover.com/oauth/v2/token' \
  --header 'content-type: application/json' \
  --data '{
      "code": "{AUTHORIZATION_CODE}",
      "client_id": "{APP_ID}",
      "client_secret": "{APP_SECRET}"
  }'
```

**Response:**

```json
{
    "access_token": "{ACCESS_TOKEN}",
    "access_token_expiration": 1709498373
}
```

The response returns only an `access_token` along with its expiration date in Unix timestamp format. No `refresh_token` is included.

---

## Sandbox and production environment URLs

| Request path | Sandbox URL | Production — North America | Production — Europe | Production — Latin America |
|---|---|---|---|---|
| `/oauth/v2/authorize` | `sandbox.dev.clover.com` | `www.clover.com` | `www.eu.clover.com` | `www.la.clover.com` |
| `/oauth/v2/token` | `apisandbox.dev.clover.com` | `api.clover.com` | `api.eu.clover.com` | `api.la.clover.com` |
| `/oauth/v2/refresh` | `apisandbox.dev.clover.com` | `api.clover.com` | `api.eu.clover.com` | `api.la.clover.com` |
| `/oauth/token/migrate_v2` | `apisandbox.dev.clover.com` | `api.clover.com` | `api.eu.clover.com` | `api.la.clover.com` |

---

## Related topics

- [Authenticate with v2/OAuth flow](https://docs.clover.com/dev/docs/use-oauth)
- [Generate OAuth expiring (access and refresh) token](https://docs.clover.com/dev/docs/generate-expiring-tokens-using-v2-oauth-flow)
- [Legacy token migration flow](https://docs.clover.com/dev/docs/legacy-token-migration-flow)
- [Blog: Expiring OAuth Tokens — Securing Clover Merchant Data](https://medium.com/clover-platform-blog/expiring-oauth-tokens-securing-clover-merchant-data-16243b9c00cc)
- [Blog: Fiddling Through Digital Keys — Clover Auth Tokens and Ecommerce Keys](https://medium.com/clover-platform-blog/understanding-clover-auth-tokens-and-ecommerce-keys-4e048827afa2)

---

*Source: [https://docs.clover.com/dev/docs/refresh-access-tokens](https://docs.clover.com/dev/docs/refresh-access-tokens)*
