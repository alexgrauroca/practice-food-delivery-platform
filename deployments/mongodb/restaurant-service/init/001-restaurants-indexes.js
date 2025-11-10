db = db.getSiblingDB('restaurant_service');

db.customers.createIndex(
    { vat_code: 1 },
    {
        unique: true,
        partialFilterExpression: { active: true }
    }
);