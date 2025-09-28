db = db.getSiblingDB('authentication_service');

db.staff.createIndex(
    { email: 1 },
    {
        unique: true,
        partialFilterExpression: { active: true }
    }
);