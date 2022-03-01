db.createUser(
        {
            user: "dockerMongoUser",
            pwd: "dockerMongoPassword",
            roles: [
                {
                    role: "readWrite",
                    db: "dockerdb"
                }
            ]
        }
);