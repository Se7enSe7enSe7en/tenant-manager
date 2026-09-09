## Entity Relationship Diagram

```Mermaid
---
title: ERD
---
erDiagram
    User {
        uuid id PK
        string email UK
        string name "[null] nullable"
    }

    Identity {
        uuid id PK
        uuid user_id
        string provider "eg. 'local', 'google'"
        string provider_user_id "[unique] email for local, Google's sub for google"
        string password_hash "[null] [local] for provider='local', otherwise NULL"
        date created_at "default: now"
    }
    Identity }|--|| User : "A user can have multiple identities"

    Session {
        uuid id pk
        uuid user_id
        date expires_at
        date created_at
    }
    Session }o--|| User : "A user can have zero or many sessions"

    Property {
        uuid id PK
        uuid user_id FK
        string name
        float32 rent_amount
    }
    User ||--o{ Property : "A user has zero or more property"

    Tenant {
        uuid id PK
        string name
        string email UK
        string phone_number "[null]"
        date created_at
    }

    Lease {
        uuid id PK
        uuid property_id FK "[unique] 1 to 1 relationship with property"
        uuid tenant_id FK "[unique]"
        int32 expected_rent_day "default: first day of the month"
        date start_date "default: now(), time lease record is created"
        date expiry_date
        bool is_month_advance
        decimal deposit_amount
    }
    Property ||--o| Lease : "A property has 0 or 1 lease"
    Tenant ||--|| Lease : "A tenant belongs to 1 lease"

    Trade {
        uuid id PK
        uuid lease_id FK
        string type "eg. 'rent', 'deposit'"
        decimal paid_amount
        date start_date "[null] [rent] for type='rent', otherwise NULL"
        date end_date "[null] [rent]"
        string note "[null] [deposit] can be used for type='deposit', reason for reducing deposit"
        date created_at
    }
    Lease ||--o{ Trade : "a Lease can have 0 or many Trade (transactions)"

```

## To view mermaid in vs code use an extension

- https://marketplace.visualstudio.com/items?itemName=bierner.markdown-mermaid

## additional notes

- `provider_user_id` field under the Identity table is the external identifier the provider uses to identify this account.
  - For local: it's the email. That's the string the user types to log in.
  - For google: it's Google's sub claim — Google's permanent unique ID for that account.
  - Why it exists: Login lookup goes "user typed email X, do I have an identity with provider='local' AND provider_user_id=X?" So provider_user_id is what you look up against the user-supplied login key.

## postgresql notes

### `NUMERIC(10,2)`

- for money fields
- max value is: 99,999,999.99
