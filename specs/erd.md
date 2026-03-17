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

    Property {
        uuid id PK
        uuid user_id FK
        string name
        float32 rent_amount
    }
    User ||--o{ Property : "A user has zero or more property"

    Tenant {
        uuid id PK
        uuid property_id FK
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

## postgresql notes

### `NUMERIC(10,2)`

- for money fields
- max value is: 99,999,999.99
