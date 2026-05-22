# Work with API usage and rate limits

There are four different rate limits placed on Clover REST API usage. For consistent access to the API services, ensure that your apps adhere to the rate limits and best practices.

---

## Rate limits

Rate limits are set for the number of new API requests your apps can make **per second**.

| Level | Per-second rate limit |
|-------|-----------------------|
| Per app | 50 |
| Per token | 16 |

---

## Concurrent rate limits

In addition to rate limiting, Clover maintains concurrent rate limiting to better control multiple retry requests to CPU-intensive API endpoints. Concurrent rate limits are set for the number of **in-progress** API requests your apps can make at the same time.

| Level | Concurrent rate limit |
|-------|-----------------------|
| Per app | 10 |
| Per token | 5 |

For instance, if the concurrent rate limit has been set to five, you can make up to five running API requests at the same time. While the server is responding to these five requests, a sixth request will return a `429 HTTP Too Many Requests` error message.

---

## 429 HTTP error messages

You receive a `429 Too Many Requests` error message when the number of API requests made by your app exceeds the rate limits.

### Best practices to avoid 429 HTTP error messages

Use the following best practices to reduce interruptions due to 429 HTTP error messages:

- Use [Webhooks](https://docs.clover.com/dev/docs/webhooks) to avoid constant polling for updates.
- Cache your data for storing specialized values or for rapidly reviewing large data sets.
- Query with the `modifiedTime` filters to avoid re-polling unmodified data.
- Stagger requests to multiple merchants to avoid bursts of high traffic volume.
- Backfill data only during non-peak business hours. Also, use the [Export API](https://github.com/clover/export-api-examples) to backfill orders or payments data older than two months.

---

### Types of 429 HTTP error messages

#### X-RateLimit-{type}

The following table can help you identify the type of 429 HTTP error message you are receiving in the response header.

| Level | Response header |
|-------|-----------------|
| Per app | `X-RateLimit-crossTokenLimit` |
| Per token | `X-RateLimit-tokenLimit` |
| Per app (Concurrent) | `X-RateLimit-crossTokenConcurrentLimit` |
| Per token (Concurrent) | `X-RateLimit-tokenConcurrentLimit` |

#### retry-after

Concurrent rate limit responses also include a `retry-after` header. The value of this header is the number of seconds your app must wait before making another request. If a request is made before this time, the waiting period increases.

Example with a `retry-after` interval of 21 seconds:

```
retry-after: 21
```

---

### Respond to 429 HTTP error messages

Design your app with necessary fallbacks to gracefully handle 429 HTTP error messages:

- In response to a 429 message, your app must **pause for one second** before making another API request.
- If you continue to receive 429 messages, your app must **exponentially increase** the time between API requests until you receive a successful response.

The following example demonstrates exponential backoff:

**Exponential backoff — Python sample**

```python
max_attempts = 10
attempts = 0

while attempts < max_attempts:
    # Make a request to Clover REST API
    response = requests.get(request_url, headers={"Authorization": "Bearer " + api_token})

    # If not rate limited, break out of while loop and continue with the rest of the code
    if response.status_code != 429:
        break

    # If rate limited, wait and try again
    time.sleep((2 ** attempts) + random.random())
    attempts = attempts + 1
```

---

*Source: [https://docs.clover.com/dev/docs/api-usage-rate-limits](https://docs.clover.com/dev/docs/api-usage-rate-limits)*
