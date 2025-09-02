db = db.getSiblingDB('authentication_service');

db.refresh_tokens.createIndex(
    { token: 1 },
    { unique: true }
);