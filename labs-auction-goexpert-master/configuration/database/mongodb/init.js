db = db.getSiblingDB("auctions"); // Nome do banco de dados

// Inserir dados na collection "users"
db.users.insertMany([
    { "name": "Admin", "_id": "a0d06633-4f0d-42ce-a653-d41a2d4aff94" },
]);

print("Dados inicializados com sucesso!");