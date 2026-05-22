# Use Webhooks

> **Regions:** North America | Europe | Latin America

Webhooks send you notifications when merchants who have installed your app perform certain events, such as creating an inventory or updating an order. For example, if you have Read Inventory permission and the merchant updates their inventory, you can set up webhooks to receive a POST request to an endpoint you specified.

---

## Prerequisites

- Install your app on the test merchant.
- Provide the Read permission required to receive webhook notifications for an event. For example, the Read Inventory permission is needed for events related to the inventory item webhooks. See [Event type keys](#event-type-keys).
- Set up a publicly accessible HTTPS endpoint on your app as a webhook receiver. **Localhost does not work for webhooks.**

---

## Set up a test URL for webhooks

You must have a simple web server running on your computer that can receive a webhook setup request and verification code. You can:

- Create a web server and test the webhook notifications. You can use a tool such as [Postman](https://www.postman.com/) for testing.
- Use one of the following tools to make your localhost available for testing after your web server is running:
  - **[ngrok](https://ngrok.com/):** Connects your local server to a public endpoint (`http://subdomain.ngrok.com`)
  - **[Pagekite](https://pagekite.net/):** Syncs your local server to a public endpoint (`yourname.pagekite.me`) (requires Python)

---

### Step 1: Configure a callback link (URL)

You need to configure a callback link (URL) to receive notifications. Clover supports only **HTTPS-enabled** callbacks. The response from your URL needs to be a `200 OK` code.

1. Log in to your sandbox [Developer Dashboard](https://sandbox.dev.clover.com/developer-home/dashboard).
2. From the left navigation menu, click **Your Apps** > *App name* > **App Settings**. The App Settings page appears.
3. Click **Webhooks**. The Edit Webhooks page appears.
4. In the **Webhook URL** field, enter a callback link (URL).
   - Example using Postman mock server: `https://021f8db8-195b-45b2-b7ae-0c4173efcdfb.mock.pstmn.io/nwtr`
5. Click **Send Verification Code**.
   - The Verification Code field displays on the Edit Webhooks page.
   - A POST request with a verification code is sent to your callback URL.

**Sample payload:**

```json
{"verificationCode":"5220ecf5-7dea-4396-b0ba-a1659c182887"}
```

6. From the POST request, copy the `verificationCode` value, paste it in the **Verification Code** field, and then click **Verify**.
7. Click **Save**. The Events Subscriptions section displays on the Edit Webhooks page.

---

### Step 2: Subscribe to event types

You need to subscribe to event types to send webhook notifications based on event triggers in your app.

1. After completing the steps to configure your callback URL, subscribe to categories of event types in the **Edit Webhooks—Events Subscriptions** section.
2. Select a checkbox to subscribe to an event type.
3. Click **Save**. The webhook callback URL and verification code display in the Webhooks section on the App Settings page.

> ⚠️ **IMPORTANT**
>
> Each event type subscription requires the corresponding read permission from the merchant. If you change permissions after a merchant installs your app, the permissions won't update for that merchant until the merchant uninstalls and reinstalls the app.

---

### Step 3: Validate and test notifications success rate

After the webhook setup for an event is complete, you receive webhook requests at the callback URL each time the event occurs.

#### Verify the Clover auth code header

Use the Clover Auth Code key to verify that your webhook messages originate from Clover.

1. Log in to the [Developer Dashboard](https://sandbox.dev.clover.com/developer-home/dashboard).
2. From the left navigation menu, click **Your Apps** > *App name* > **App Settings**. The Webhooks section displays the **Clover Auth Code** key. This key displays in every message header after the webhook callback URL is validated.

**Sample webhook message header:**

```
"X-Clover-Auth":"0307e264-b33a-4913-88aa-37b10014a0b8"
```

#### Test webhook notifications

Clover recommends testing the success rate of notifications received through the webhook service.

1. Log in to the [Merchant Dashboard](https://sandbox.dev.clover.com/dashboard).
2. From the left menu, click **Inventory** > **Item**. The available items appear.
3. Select an item checkbox. The Item Basics page appears.
4. Edit the item name and click **Save**. Clover sends out a notification to the webhook endpoint.

**Sample notification payload:**

```json
{
    "appId": "DRKVJT2ZRRRSC",
    "merchants": {
        "XYZVJT2ZRRRSC": [
            {
                "objectId": "O:GHIVJT2ABCRSC",
                "type": "CREATE",
                "ts": 1537970958000
            },
            {
                "objectId": "O:ABCVJTABCRRSC",
                "type": "UPDATE",
                "ts": 1536156558000
            }
        ],
        "MNOVJT2ZRRRSC": [
            ...
        ]
    }
}
```

---

## Understand the message format and event objects

### Webhook message format

The messages that Clover sends to your webhook URL contain the following information:

| Field | Description |
|-------|-------------|
| `appId` | Identifier (ID) of the application that sends the data updates. |
| `merchants` | One or more merchant arrays indicated by the merchant identifier (`mID`). |
| `update` | One or more `update` objects, each containing an `objectId`, `type`, and `ts`. See [Update object](#update-object). |

### Update object

Each update is indicated by a JSON or XML object for the webhook event containing the following information:

| Field | Description |
|-------|-------------|
| `objectId` | `<key of event type>:<event object ID>` — See [Event type keys](#event-type-keys). |
| `type` | Operation type: `CREATE`, `UPDATE`, or `DELETE`. |
| `ts` | Timestamp in Unix time (milliseconds). |

### Event type keys

The `objectId` value begins with a key to indicate the event type:

| Key | Event type | Description | Permissions |
|-----|------------|-------------|-------------|
| `A` | Apps | App is installed, uninstalled, or the subscription is changed. | Read merchant |
| `C` | Customers | Customer record is created, updated, or deleted. | Read customers |
| `CA` | Cash adjustments | Cash log event occurs. | Read merchant |
| `E` | Employees | Employee record is created, updated, or deleted. | Read employees |
| `I` | Inventory | Inventory item is created, updated, or deleted. | Read inventory |
| `IC` | Inventory category | Inventory category is created, updated, or deleted. | Read inventory |
| `IG` | Inventory modifier group | Inventory modifier group is created, updated, or deleted. | Read inventory |
| `IM` | Inventory modifier | Inventory modifier is created, updated, or deleted. | Read inventory |
| `O` | Orders | Order is created, updated, or deleted. | Read orders |
| `M` | Merchants | Merchant property is changed, or a new merchant is added. | Read merchant |
| `P` | Payments | Payment is created or updated. | Read payments |
| `SH` | Service hour | Service hour is created, updated, or deleted. | Read merchant |

---

## Troubleshoot webhook notifications

If you are not receiving webhook notifications, use the following checklist:

| | Description |
|-|-------------|
| ✅ | Verify that your webhook URL is correct. |
| ✅ | Verify that the webhook URL is set up to accept POST requests. |
| ✅ | Make sure your endpoint is an HTTPS webhook address with a valid SSL certificate that can correctly process event notifications. |
| ✅ | Verify if your app is already subscribed to a webhook if you set up a new webhook server. If already subscribed, clear the webhook, save, then reselect and save again to start receiving new webhook messages. |
| ✅ | Make sure you are subscribed to the event type for which you want to receive a notification. |
| ✅ | Verify that your app has required permissions for that event type. For example, the Read inventory permission is needed for inventory item webhooks. |
| ✅ | Check if the merchant has uninstalled and then reinstalled the app if you have changed app permissions after the merchant installed your app. |

---

## Appendix: Use Postman to test your webhook

If you have a [Postman](https://www.postman.com/) account, you can use it to create a mock server for the callback URL, get a verification code, and test your webhook.

### Step 1: Create a Webhook collection in Postman

1. Log in to [Postman](https://www.postman.com/) and access your Workspace.
2. Click **Create New Collection** > **Blank Collection**. The New Collection page appears.
3. Rename the new collection to indicate your webhook collection, for example, **NewWebhookTest**.
4. In the left navigation menu, under the new collection, click **Add a request**.
5. Replace the default request name with something descriptive, for example, **NWTR** (New Webhook Test Request).
6. From the HTTP method drop-down list, select **POST** and enter the request name.
   > **Note:** POST webhook requests respond to the request body and also contain properties like authentication tokens.
7. Click **Save**.
8. In the left navigation menu, click the ellipsis icon (⋯) next to the POST request and click **Add example**.
9. In the POST request—Body section, enter a sample response, for example:
   ```json
   {"testResponse": "pass"}
   ```
10. Click **Save**.

### Step 2: Create a mock server

Mock servers let you simulate API endpoints by returning example responses linked to each request.

1. After you create a POST request for your webhook collection, in the left navigation menu, click the ellipsis icon (⋯) next to the webhook collection name and click **Mock collection**. The Create a mock server page appears.
2. In the **Mock Server Name** field, enter the POST request name, for example, **NWTR**.
3. Click **Create Mock Server**. After the server is running, when a request is sent, the logs display on the Mock server calls page.
4. From the top right corner, click **Copy URL** to copy the mock server URL. You need this URL to set up webhooks in your [Developer Dashboard](https://sandbox.dev.clover.com/developer-home/dashboard).
5. In your Postman Workspace, click **Refresh Logs** and copy the verification code that displays.
6. In the Developer Dashboard—Edit Webhooks page, in the **Verification Code** field, enter the copied verification code and click **Verify**.
7. Complete the steps to [test webhook notifications](#test-webhook-notifications). In your Postman Workspace, click **Refresh Logs** and review the notification payload that displays.

---

*Source: [https://docs.clover.com/dev/docs/webhooks](https://docs.clover.com/dev/docs/webhooks)*
