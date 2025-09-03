db = db.getSiblingDB('authentication_service');

db.customers.createIndex(
    { email: 1 },
    {
        unique: true,
        partialFilterExpression: { active: true }
    }
);