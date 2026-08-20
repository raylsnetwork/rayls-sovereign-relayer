// mongo-init.js

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

async function waitForInit() {
    while (true) {
        const status = rs.initiate({
            "_id": "rs0",
            "members": [
                {
                    "_id": 0,
                    "host": "mongodb:27017"
                }
            ]
        });
        if (status.ok === 1) {
            break;
        }
        await sleep(500);
    }
}

async function waitForPrimary() {
    while (true) {
        const status = rs.status();
        const isMaster = db.isMaster()
        if (status.myState === 1 && status.ok === 1 && isMaster.ismaster === true) {
            break;
        }
        await sleep(500);
    }
}



waitForInit()
    .then(() => {
        return waitForPrimary()
    })
    .then(() => {
        const adminDB = db.getSiblingDB('admin');
        try {
            // Create admin user if it doesn't exist
            adminDB.createUser({
                user: 'admin',
                pwd: 'admin',
                roles: [
                    {
                        role: 'root',
                        db: 'admin'
                    }
                ]
            });

            console.log('MongoInit: Admin user created successfully. MongoDB is now ready for operation.')
        } catch (e) {
            var status = rs.status();
            console.error("MongoInit: error creating user:", e, status);
        }

    }).catch(err => {
        console.error('MongoInit: Error waiting for mongo readyness:', err);
    });
