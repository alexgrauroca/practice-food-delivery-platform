db = db.getSiblingDB('authentication_service');

db.staff.createIndex(
    {
        email: 1,
        restaurant_id: 1
    },
    {
        unique: true,
        partialFilterExpression: { active: true }
    }
);