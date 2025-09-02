db = db.getSiblingDB('customer_service');

db.customers.createIndex(
    { email: 1 },
    {
        unique: true,
        partialFilterExpression: { active: true }
    }
);