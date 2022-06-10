/*eslint no-undef: "error"*/
// noinspection JSUnresolvedVariable,JSUnresolvedFunction
db.createUser(
        {
            user: "dockerMongoUser",
            pwd: "dockerMongoPassword",
            roles: [
                {
                    role: "readWrite",
                    db: "dockerdb",
                },
            ],
        },
);
