# Notification Service

Multi-channel notification system with event-driven architecture.

## Features

- ✅ Multi-channel notifications (Email, Push, In-App)
- ✅ Event-driven with RabbitMQ
- ✅ PostgreSQL persistence
- ✅ Mock senders (perfect for demo)
- ✅ Multiple concurrent workers
- ✅ Order & User event consumers

## Architecture

```
RabbitMQ Events → Consumer → Notification Service → Channels
                                                     ├─ Email (Mock)
                                                     ├─ Push (Mock)
                                                     └─ Database
```

## Quick Start

```bash
# Start dependencies
docker compose up -d notification-db tokohobby-rabbitmq

# Start worker
docker compose up -d notification-worker

# Check logs
docker logs -f notification-worker
```

## Environment Variables

```bash
DB_HOST=notification-db
DB_PORT=5432
DB_USER=user
DB_PASSWORD=supersecret123
DB_NAME=notifications
RABBITMQ_URL=amqp://admin:admin123@tokohobby-rabbitmq:5672/tokohobby
MOCK_MODE=true
LOG_LEVEL=info
```

## Event Consumers

### Order Events
- Queue: `notifications.order.events`
- Workers: 5
- Events: OrderCreated, OrderStatusChanged, OrderShipped

### User Events  
- Queue: `notifications.user.events`
- Workers: 3
- Events: UserRegistered (Welcome email)

## Database Schema

```sql
CREATE TABLE notifications (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    type VARCHAR(50),
    category VARCHAR(50),
    title VARCHAR(255),
    message TEXT,
    channels TEXT[],
    status VARCHAR(20),
    created_at TIMESTAMP
);
```

## Testing

```bash
# Watch logs
docker logs -f notification-worker

# Check database
docker exec -it notification-db psql -U user -d notifications \
  -c "SELECT * FROM notifications ORDER BY created_at DESC LIMIT 10;"
```

## Documentation

- [Implementation Guide](./IMPLEMENTATION_SUMMARY.md)
- [Database Setup](./DATABASE_IMPLEMENTATION.md)
- [Deployment Walkthrough](./WALKTHROUGH.md)
