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

Check out [Swagger](https://goldexrobot.github.io/core.integration.backend/).

![Goldex environment](/docs/images/goldex_env.png)

---

## API

Goldex backend exposes an API to provide some extended information like photos, storage access history etc.
Moreover the API allows business backend to control a vending machine.

Calls to the API must be supplied with basic HTTP auth header. You can get the login/key in the Goldex dashboard.

See [Swagger](https://goldexrobot.github.io/core.integration.backend/#goldex-api-v1).

---

## Callbacks

Goldex backend sends HTTP requests to notify the business backend about a new events or to request some information for the vending terminal in real time.

For instance: items evaluation, storage interaction, etc.

Exact endpoints to call by Goldex should be defined in Goldex dashboard.

The HTTP requests are: **callbacks** and named **proxy methods** (see below).

Requests are always of method **POST** and carry **application/json; charset=utf-8** with the headers:

| Header | Meaning | Example |
| --- | --- | --- |
| X-CBOT-PROJECT-ID | Origin project ID | "1" |
| X-CBOT-BOT-ID | Origin bot ID (uint64) | "42" |

All the Goldex requests are signed (see below) and should be validated at business backend side.

Business backend has to respond with successful HTTP status (200, 201, or 202) to signalize about callback consumption. Until then, Goldex will continue to send callback requests.

### Evaluation callbacks

Evaluation related callbacks are **optional** but give extra control over UI from backend side.

See in [Swagger](https://goldexrobot.github.io/core.integration.backend/#business-callbacks).

### Storage callbacks

Storage callbacks are sent during the access to the storage and are **mandatory**.

See in [Swagger](https://goldexrobot.github.io/core.integration.backend/#business-callbacks)

#### Domain/flow

A storage cell could be occupied or released in multiple scenarios. To control the access to the cells and keep a history for audit there are **domains**.

For example, an item is bought from a customer during buyout flow. It's now under *buyout* domain. Then the item is collected from the robot (flow is outside of UI), so it's released under *collection* domain.

Another example. To sell an item it should be first loaded into the robot. In this particular scenario item is loaded using internal robot dashboard (outside of UI), so the loading action (cell occupation) has been made under *dashboard* domain/flow.

Here are predefined domains/flows:

*Domain* is an origin of the cell operation, or in other words a *business flow* in terms of UI. Here how domains and cell events are correlate:
| Domain/flow | Cell occupation | Cell release | Comment |
| --- | --- | --- | --- |
| buyout | YES | | Buyout business flow: a cell can only be **occupied** |
| shop | | YES | Shop flow: a cell can only be **released** |
| pawnshop | YES | YES | Pawhshop flow: a cell can be both **occupied/released** |
| collection/dashboard | YES | YES | Not a UI flow, but a storage management (using internal bot dashboard): a cell can be both **occupied/released** |

Note about a shop flow. First of all, an item should be loaded into the storage of the bot. Then it could appear in the UI as a product we're selling. \
So the storage is loaded through the internal dashboard of the bot, therefore in this case the cell is loaded under the collection/dashboard domain.

### Named proxy methods (UI methods)

Proxy methods are custom named callbacks (name-to-endpoint bindings) and available to call directly from the UI.

In this case Goldex backend acts as secure proxy between the robot and business backend (initially UI doesn't have an identity and can't be verified by the business backend side).

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

Now Goldex backend waits for a response from business backend (status 200, application/json) and returns the response back to the UI. Voila! The robot is authenticated on business backend side.

---

### Callbacks signature

Goldex signs callbacks with JWT. Token is signed with a per-project key (see Goldex dashboard) and is transferred in `Authorization` HTTP header (bearer).

Business backend **SHOULD** validate those callbacks. Developer is fully responsible for the security.

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
