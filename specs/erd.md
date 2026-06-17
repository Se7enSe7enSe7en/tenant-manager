## Entity Relationship Diagram

```Mermaid
---
title: ERD
---
erDiagram
    User {
        uuid id PK
        string email UK
        string name "? (nullable)"
    }

    Identity {
        uuid id PK
        uuid user_id
        string provider "[unique] eg. 'local', 'google'"
        string provider_user_id "[unique] email for local, Google's sub for google"
        string password_hash "NULL unless provider='local'"
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
        uuid property_id FK "[unique] 1 to 1 relationship with property"
        string name
        string email UK
        string phone_number "?"
        int8 expected_rent_day "default: first day of the month"
        date created_at
    }
    Property ||--o| Tenant : "A property can have zero or one tenant representative"

    Trade {
        uuid id PK
        uuid tenant_id FK
        uuid property_id FK
        uuid user_id FK
        float32 paid_amount
        date start_date
        date end_date
        date created_at
    }
    Tenant ||--o{ Trade : "A tenant can have zero or more trade (transactions)"
    User ||--o{ Trade : "A user can have zero or more trade"
    Property ||--o{ Trade : "A property can have zero or more trade"

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
