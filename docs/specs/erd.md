```Mermaid
---
title: ERD
---
erDiagram
    User {
        uint64 id PK
        string email UK
        string name "? (nullable)"
    }

    Property {
        uint64 id PK
        string name
        uint32 rent_amount
        uint64 user_id FK
        uint64 tenant_id FK
    }
    User ||--o{ Property : "A user has zero or more property"
    Property ||--o| Tenant : "A property can have zero or one tenant representative"
    
    Tenant {
        uint64 id PK
        string name
        string email UK
        string phone_number "?" 
        date pay_date "default: first day of the month"
        date created_at
    }
 
    Trade {
        uint64 id PK
        uint64 tenant_id FK
        uint64 property_id FK
        uint64 user_id FK
        unint32 paid_amount
        date start_date
        date end_date
        date created_at
    }
    Tenant ||--o{ Trade : "A tenant can have zero or more trade (transactions)"
    User ||--o{ Trade : "A user can have zero or more trade"
    Property ||--o{ Trade : "A property can have zero or more trade"

```