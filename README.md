# Goldex Robot: Integration

Goldex Robot is a vending machine that evaluates gold/silver valuables, sells coins and has internal storage/safebox.

Integration with Goldex consists of two major parts: backend integration and [UI](https://github.com/goldexrobot/core.integration.ui) integration.

This document covers backend integration.

---

## TL;DR

Goldex vending machine serves UI and communicates with core Goldex backend.

This core backend exposes API to control the machine from outside.

The core backend sends optional signed callbacks to the "business" backend to notify about valuable evaluation steps and storage interaction.

"Business" backend should implement methods required by the UI and is responsible for the secure data transmission.

Check out [API, callbacks data and signing playground](https://goldexrobot.github.io/core.integration.backend/).

![Goldex environment](/docs/images/goldex_env.png)

---

## Callbacks

Goldex backend sends HTTP requests to notify the business backend about a new events or to request some information for the vending terminal in real time.

For instance: bot status changes, items evaluation, storage interaction, etc.

Exact endpoints to call by Goldex should be defined in Goldex dashboard.

The HTTP requests are: **callbacks** and custom **UI methods**.

Requests are always of method **POST** and carry **application/json; charset=utf-8** with headers:

| Header | Meaning | Example |
| --- | --- | --- |
| X-CBOT-PROJECT-ID | Origin project ID | "1" |
| X-CBOT-BOT-ID | Origin bot ID (uint64) | "42" |

All the Goldex requests are signed (see below) (check out Goldex dashboard for a verification public key) and could be optionally transferred using mutual TLS.

Business backend have to respond with successful HTTP status (200, 201, or 202) to signalize about callback consumption. Until then Goldex will continue to send callback requests.

Callback models are described in [Swagger](https://goldexrobot.github.io/core.integration.backend)

### Evaluation callbacks

| Callback | Description | Model | Inside | Expected response |
| --- | --- | --- | --- | --- |
| Started | New item evaluation is created | `EvalStarted` | Evaluation ID | Status 200 |
| Cancelled | Item evaluation is cancelled or failed | `EvalCancelled` | Evaluation ID, reason | Status 200 |
| Finished | Item evaluation is successfully finished | `EvalFinished` | Detailed evaluation data, photo file IDs | Status 200 |

### Storage callbacks

| Callback | Description | Model | Inside | Expected response |
| --- | --- | --- | --- | --- |
| Cell occupation | Cell is occupied with some item | `StorageCellEvent` | Cell address and *domain* | Status 200 |
| Cell release | Cell is freed and item is released | `StorageCellEvent` | Cell address and *domain* | Status 200 |
| Pre-occupation | Optional, called before occupation to give more flexibility | `StorageCellEvent` | Cell address and *domain* | Status 200 and `{"allowed":true}` |
| Pre-release | Optional, called before release to give more flexibility | `StorageCellEvent` | Cell address and *domain* | Status 200 and `{"allowed":true}` |

*Domain* is an origin of the cell operation:
| Domain | Whats |
| --- | --- |
| buyout | Cell is occupied during buyout flow |
| shop | Cell is freed during item purchasing / shop flow |
| pawnshop | Cell is occupied or freed during pawnshop flow |
| collection | Cell is occupied/freed by staff members (on-bot collection dashboard) |
| dashboard | Cell is occupied/freed by staff members (on-bot system dashboard) |
| other | For custom UI flows |

### UI methods

UI methods are custom named callbacks and available to call directly from the UI.

In this case Goldex backend acts as secure proxy, for instance to authenticate a machine on business backend side.

For instance, let's assume you defined a method `auth` in Goldex dashboard and bound it to `https://example.com/bot/auth`.

UI now calls `auth` with some payload. Goldex backend adds bot/project IDs to the request, signs it and and sends to `POST https://example.com/bot/auth`:

```json
{
  "project_id": 1,
  "bot_id": 42,
  "payload": {
    // bot`s request k/v goes here:
    "foo": {
      "bar": "baz"
    }
  }
}
```

Now Goldex backend waits for a response from your backend (status 200, application/json) and re-sends the response back to the UI.

---

## API

Goldex backend exposes an API to provide some extended information like photos, storage access history etc.
Moreover the API allows business backend to control a vending machine.

Calls to the HTTP API must be properly signed (see below) with per-project private key. You can get the key in the Goldex dashboard.

GRPC API is also available.

[API in Swagger](https://goldexrobot.github.io/core.integration.backend/).

---

## Signing

A request to Goldex API have to be properly signed, we are using JWT.

JWT token should be transferred in `Authorization` HTTP header with `Bearer` prefix:

Goldex callbacks are also signed with per-project key. It is not mandatory but is preferred to validate those callbacks. Developer is fully responsible for the security.

```text
Authorization: Bearer <jwt.goes.here>
```

### JWT claims

Here are default fields of JWT used during signing a request to Goldex API:

| Field | Purpose | Format | Example |
| --- | --- | --- | --- |
| aud | Recipient of the request | string(3-32): alphanumeric, `-`, `_` | `goldex` |
| iss | API login (issuer) | string(3-32): alphanumeric, `-`, `_` | `my_login` |
| jti | Unique request ID | string(6-36): alphanumeric, `-` (UUID compatible) | `16d91811-2b08-4462-8754-299eb1810963` |
| sub | The request entity or domain | string(32): alphanumeric | `request` |
| exp, nbf, iat | Expiration, not before and issuance time | According to [RFC 7519](https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.5) | - |

Additional JWT fields:

| Field | Purpose | Format | Example |
| --- | --- | --- | --- |
| bha | Request body hash algorithm | string(16): `SHA-256`, `SHA-384`, `SHA-512`, `SHA3-224`, `SHA3-256`, `SHA3-384`, `SHA3-512` | `SHA-512` |
| bhs | Request body hash | string(32-128): hexadecimal without leading `0x` | `5ea71dc6...ae04ee54` |
| mtd | Request method | string(8): GET, POST etc. | `POST` |
| url | Request URL | string(256): valid URL | `https://example.com` |

Body hash algorithm and hash fields **have to be empty** for bodiless request such as GET.

Goldex **callbacks** carries JWT with the next content:

| Field | Content |
| --- | --- |
| aud | ["project-N"] where N is project ID |
| iss | "goldex" |
| sub | "notification" or "request" depending on context |

Allowed JWT __signing__ algorithms: `HS256` (HMAC SHA-256), `HS384` (HMAC SHA-384), `HS512` (HMAC SHA-512).

Try out signing [here](https://goldexrobot.github.io/core.integration.backend/).
