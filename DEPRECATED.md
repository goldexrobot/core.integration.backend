# Deprecated

Callbacks and proxy endpoints are deprecated.

## Callbacks

Goldex backend sends optional HTTP requests to notify the business backend about new events or to request some information from the machine in a real time.

See `POST /callbacks` in [Swagger](https://goldexrobot.github.io/core.integration.backend/#api-v1).

The HTTP requests are: **callbacks** and **named proxy endpoints** (see below).

Requests are always of method **POST** and carry **application/json; charset=utf-8** with the headers:

| Header | Meaning | Example |
| --- | --- | --- |
| X-CBOT-PROJECT-ID | Origin project ID | "1" |
| X-CBOT-BOT-ID | Origin bot ID (uint64) | "42" |

All the Goldex requests are signed (see below) and **should be** validated at the business side.

Business backend has to respond with the successful HTTP status (200, 201, or 202) to signalize about callback consumption. Until then, Goldex may continue to resend requests.

### Named proxy endpoints

Proxy endpoints are optional named callbacks (name-to-endpoint bindings) and are available at the UI side.

The purpose of endpoints is to give the UI a secure way to call business backend. For instance, to get an access token.

See `POST /proxy` in [Swagger](https://goldexrobot.github.io/core.integration.backend/#api-v1).

In this case Goldex backend acts as a secure proxy between the machine and business backend.

For instance, let's assume business backend have defined an endpoint `auth` pointing to `https://example.com/bot/auth`.

UI now calls proxy method `auth` using UI API. Goldex signs the request and sends `POST https://example.com/bot/auth`:

```json
{
  "project_id": 1,
  "bot_id": 42,
  "payload": {
    // optional k/v from the machine side
  }
}
```

Then Goldex backend waits for a response from the business backend (status 200, application/json) and returns the response back to the UI.

---

### Signing

Goldex signs callbacks with a JWT. Token is signed with a per-project secret and is transferred in `Authorization` HTTP header (bearer).

Business backend **should** validate those callbacks. Developer is fully responsible for the security.

```text
Authorization: Bearer <jwt.goes.here>
```

#### JWT claims

Here are default fields of JWT used during signing:

| Field | Purpose | Format | Example |
| --- | --- | --- | --- |
| aud | Recipient name | string(3-32): alphanumeric, `-`, `_` | `project-1` |
| iss | Issuer name | string(3-32): alphanumeric, `-`, `_` | `goldex` |
| jti | Unique request ID | string(6-36): alphanumeric, `-` (UUID compatible) | `16d918112b0844628754299eb1810963` |
| sub | Subject (request or notification) | string(32): alphanumeric | `request` |
| exp, nbf, iat | Expiration, not before and issuance time | According to [RFC 7519](https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.5) | - |

Additional JWT fields:

| Field | Purpose | Format | Example |
| --- | --- | --- | --- |
| bha | Request body hash algorithm | string(16): `SHA-256`, `SHA-384`, `SHA-512`, `SHA3-224`, `SHA3-256`, `SHA3-384`, `SHA3-512` | `SHA-512` |
| bhs | Request body hash | string(32-128): hexadecimal without leading `0x` | `5ea71dc6...ae04ee54` |
| mtd | Request method | string(8): GET, POST etc. | `POST` |
| url | Request URL | string(256): valid URL | `https://example.com` |

Body hash algorithm and hash fields are empty for bodiless request such as GET.

JWT **signing** algorithms are: `HS256` (HMAC SHA-256), `HS384` (HMAC SHA-384), `HS512` (HMAC SHA-512).

Try out signing [here](https://goldexrobot.github.io/core.integration.backend/signature/).
