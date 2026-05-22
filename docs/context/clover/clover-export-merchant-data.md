# Export merchant data

> **Regions:** North America | Europe | Latin America

The Export API allows you to obtain bulk payment and order data from merchants who have installed your app. This service is not real-time, so your chance of running into API limit restrictions is lower. Consider using this service for data collection if you need merchant orders and payment histories.

> 📘 **NOTE**
>
> If you have not used this endpoint before, email [developer-relations@clover.com](mailto:developer-relations@clover.com) during your app review to inquire about permissions.

---

## Prerequisites

Install and set up the following before using the Export API:

- Python 3.7.3 (backward compatible with 2.7)
- pip Python package manager
- Clover sandbox [developer account](https://docs.clover.com/dev/docs/developer-accounts)

---

## Usage considerations

Consider the following points before using the Export API:

- **Availability time:**
  - United States (US) UTC: MON–FRI 10:00–15:00, SAT 10:00–15:00, SUN 10:00–15:00
  - Europe UTC: MON–FRI 01:00–05:00, SAT 02:00–05:00, SUN 02:00–06:00
- Sandbox is **not** limited to specific times.
- Returns **1,000 objects per file**. If your response contains more than 1,000 objects, returns an array of files to capture.
- Exported files are **deleted after 24 hours**.
- **Maximum time range is 30 days.** If your data spans are longer, break the request into multiple requests.
- The server handles a few exports concurrently. Wait for each export to finish before starting another one. The server returns `503` errors whenever the threshold is reached.

---

## Use the endpoint

### Step 1: Create the request

Make a **POST** call to `/v3/merchants/{merchant_Id}/exports` with the appropriate payload:

| Field | Description |
|-------|-------------|
| `export_type` | Type of data to export: `PAYMENTS` or `ORDERS` |
| `start_time` | Start of the time range (in UTC) |
| `end_time` | End of the time range (in UTC) |

**Request body example:**

```json
{
    "type": "ORDERS",
    "startTime": "<start_time>",
    "endTime": "<end_time>"
}
```

The exports object response includes the **export ID** for use in the next step.

---

### Step 2: Check the status and percent complete periodically

Using the export ID from the previous step, make a **GET** call to `/v3/merchants/{merchant_Id}/exports/{export_Id}`.

The following statuses may appear:

| Status | Description |
|--------|-------------|
| `PENDING` | In the queue to be processed. |
| `IN_PROGRESS` | Currently being processed. |
| `DONE` | Process complete. |

**Sample response:**

```json
{
    "status": "PENDING",
    "modifiedTime": "<modified_time>",
    "merchantRef": {
        "id": "<merchant_id>"
    },
    "percentComplete": 0,
    "startTime": "<start_time>",
    "createdTime": "<created_time>",
    "endTime": "<end_time>",
    "type": "ORDERS",
    "id": "<export_id>"
}
```

---

### Step 3: Download the exported files

After the request completes (status `DONE`), your export object includes URLs to the exported files for download.

**Export object example:**

```json
{
    "id": "<export_id>",
    "type": "ORDERS",
    "status": "DONE",
    "percentComplete": 100,
    "availableUntil": "<available_until_time>",
    "startTime": "<start_time>",
    "endTime": "<end_time>",
    "createdTime": "<created_time>",
    "modifiedTime": "<modified_time>",
    "exportUrls": {
        "elements": [
            {
                "url": "<download_url>",
                "export": {
                    "id": "<export_id>"
                }
            }
        ]
    },
    "merchantRef": {
        "id": "<merchant_id>"
    }
}
```

> ⚠️ Files are only available for **24 hours** after export. Download them promptly.

---

*Source: [https://docs.clover.com/dev/docs/exporting-merchant-data](https://docs.clover.com/dev/docs/exporting-merchant-data)*
