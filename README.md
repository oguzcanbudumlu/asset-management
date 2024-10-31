> **Project: Asset Management Service Objective**
> 
> You are required to build a small microservice-based system for managing wallets and their assets. There should be at least 2 services, but you can split these services further as you see fit.
> 
> **Requirements**
> 
> 1- Wallet Management Service: Responsible for creating, retrieving, and deleting wallets.
> Each wallet will have an address and network, and these two fields combined should be unique.
> 
> 2- Asset Management Service: Handles asset operations for existing wallets, including withdrawals, deposits, and scheduled transactions.
> It should allow withdraw and deposit operations.
> 
> It should also include a feature that allows users to schedule transactions, enabling them to send assets from one wallet to another at a specified future time.
> 
> **Example**
> 
> Wallet: 
> ```json
> {"address": "1Lbcfr7sAHTD9CgdQo3HTMTkV8LK4ZnX71", "network": "Bitcoin"}
> ```
> 
> Asset: 
> ```json
> {"name": "BTC", "amount": 10}
> ```
> **Implementation Details**
> 
> **Database:** You are free to choose any database for data storage.
> 
> **User Management:** Not required, but you may implement it if desired.
> 
> **Documentation:** Provide comprehensive documentation to set up and run the project.
> 
> **Additional Tools:** Include Docker, tests, and/or automated scripts to facilitate evaluation.
> 
> **Creativity and Questions**
> 
> We intentionally left some details unspecified to encourage creativity in your approach.


# System Architecture Diagram

```mermaid
graph TD
    User[User] -->|Create/Delete Wallet| WalletAPI[Wallet API]
    User -->|Withdraw/Deposit/Scheduled Transactions| AssetAPI[Asset API]

    subgraph WalletDB[Wallet API Database]
        wallets[["Wallets Table"]]
        wallet_deleteds[["Wallet Deleted Table"]]
    end
    WalletAPI --> WalletDB

    subgraph AssetDB[Asset API Database]
        balance[["Balance Table"]]
        scheduled_transactions[["Scheduled Transactions Table"]]
    end
    AssetAPI -->|Performs Deposit/Withdraw| balance
    AssetAPI -->|Registers Scheduled Transactions| scheduled_transactions

    AssetAPI -->|Wallet Validation| WalletAPI

    subgraph Kafka
        KafkaTopic[["Transaction Events Topic"]]
    end
    TransactionOutboxPublisher -->|Publishes Event for Scheduled Transactions| KafkaTopic
    KafkaTopic -->|Triggers Processing| TransactionConsumer
    TransactionConsumer -->|Transfers Funds| balance
    TransactionConsumer -->|Updates Status to Completed| scheduled_transactions

    TransactionOutboxPublisher -->|Polls for Due Transactions| scheduled_transactions


```


1. **Wallet Management Service (wallet-api):**
    - Creates, retrieves, and deletes wallets.
    - Each wallet has a unique combination of address and network, stored in the `wallets` table.
    - Deleted wallets are moved to the `wallet_deleteds` table for record-keeping.

2. **Asset Management Service (asset-api):**
    - Performs withdrawal and deposit operations, updating the `balance` table.
    - Supports scheduling transactions, enabling users to transfer assets from one wallet to another at a specified future time, recorded in the `scheduled_transactions` table.

3. **Transaction Outbox Publisher:**
    - Regularly polls the `scheduled_transactions` table for transactions scheduled to be executed.
    - Publishes events for due transactions to a Kafka topic, triggering the transaction processing workflow.

4. **Transaction Consumer:**
    - Listens to the Kafka topic for transaction events.
    - Processes each event by transferring funds in the `balance` table and updating the transaction status from "PENDING" to "COMPLETED" in the `scheduled_transactions` table.



```
asset-management
├── internal
├── pkg
└── services
    ├── asset-api
    │   └── main.go
    ├── transaction-consumer
    │   └── main.go
    ├── transaction-outbox-publisher
    │   └── main.go
    └── wallet-api
        └── main.go

```


> The `main.go` files above serve as the entry points for each microservice.




## Database Table Entry Examples

> **Note:** In this schema, `network` is used in place of `name` to represent the network associated with each wallet or transaction.

- Balance:
```json
{
    "wallet_address": "1Lbcfr7sAHTD9CgdQo3HTMTkV8LK4ZnX71",
    "network": "Bitcoin",
    "balance": 100.00
}
```

- Scheduled Transaction:
```json
{
    "scheduled_transaction_id": 1,
    "from_wallet_address": "1Lbcfr7sAHTD9CgdQo3HTMTkV8LK4ZnX71",
    "to_wallet_address": "1Kzo9sXeUWo12nkXnKL4WECF5DRDojps6Y",
    "network": "Bitcoin",
    "amount": 50.00,
    "scheduled_time": "2024-10-31T12:00:00Z",
    "status": "PENDING",
    "created_at": "2024-10-30T08:00:00Z"
}
```

- Wallet:

```json
{
    "id": 1,
    "network": "Bitcoin",
    "address": "1Lbcfr7sAHTD9CgdQo3HTMTkV8LK4ZnX71"
}
```

- Wallet Deleted:
```json
{
    "id": 1,
    "network": "Bitcoin",
    "address": "1Lbcfr7sAHTD9CgdQo3HTMTkV8LK4ZnX71"
}
```

