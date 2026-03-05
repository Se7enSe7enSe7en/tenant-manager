## Core Features
- the system should be able to provide the details of the receipt after the tenant has paid
    - for example:
        1. tenant # 1 just paid rent
        2. the user clicks a button or issues a command that tenant # 1 has paid
        3. system shows:
            1. current date range for the rent (eg. September 29 to October 29)
            2. amount of rent (eg. 15 000 php)



- UI (User Interface) to see if a tenant is already paid for the current date
- UI to see previous records of tenants
- CRUD (Create Read Update Delete) a tenant
    1. (Create) Register a new tenant 
    2. (Read) Show the tenant details in one page
    3. (Update) Edit tenant details
    4. (Delete) Archive the tenant to keep its records then optionally delete a tenant in the archives

- Auth, multiple users can log in
    - can be through their gmail
    - auth can be difficult (based on people in the internet), if it is hard for us as well then we could use a third party to handle this (eg. Auth0) 



## Extra Features
- System for tenants (can be a separate app but uses the same DB)
    1. check status (if paid or not)
    2. invoice or notification system for rent if almost due
    3. automated sending of receipts to google drive or any cloud storage for the tenant's past receipts (idea credits to Allen)
        1. make sure the storage is secure and only the tenant can access


- Digital receipts for tenants
    - after payment, an email of the receipt as a PDF will be sent to the tenant

- Digital payment for tenants
    - payment through the system by entering credit/debit card details (using Stripe perhaps)

- 2 different tenant payment modes "+30 days" mode or "same day of the month" mode
    - +30 days mode, lets say tenant paid jan 30, their next payment will be

- Family or Group of users, with different roles can share the same group of properties

