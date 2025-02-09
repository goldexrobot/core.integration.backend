# GR: backend integration

Goldex Robot is a vending machine that evaluates gold/silver valuables, sells coins and has an internal items storage.

This document covers backend integration.

[<img src="/docs/images/swagger-button.png" alt="Swagger" width="120"/>](https://goldexrobot.github.io/core.integration.backend/#api-v1)

---

## TL;DR

GR machine serves UI and communicates with the core Goldex backend.

This core backend exposes an API to control the machine from a business side.

![Goldex environment](/docs/images/goldex_env.png)

---

## API

Goldex core backend exposes an API to provide some extended information like machines available for some specific project, photos available, evaluation history etc.

Moreover the API allows business backend to control a vending machine (switch operational mode).

Calls to the API must be supplied with a basic HTTP auth header. Contact Goldex support to get the credentials.

