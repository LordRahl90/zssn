## Zombie Survival Social Network

* [Postman Docs](https://www.getpostman.com/collections/1fa844327cfc6b9b3894)
* [BaseURL](https://zssn-a336s4xzcq-ew.a.run.app)

## Local Setup:
Make sure docker is installed on the target machine
* Clone Repository

## Steps
* RUN `make bi` to build the docker image
* RUN `docker-compose up` to start the service and the accompanying mysql database

## OR
 * RUN `make`

## Assumptions
NB: Items are given constants:
1: Water <br />
2: Food
3: Medication
4: Ammunition


## API Documentation:
* POST `/users` -> Creates a new survivor record with the following payload and returns a JWT Token for future authentication:
```json
{
    "email": "tolaabbey009@carroll.net",
    "name": "Bart Beatty",
    "age": 20,
    "gender": "Male",
    "latitude": -78.18533654085428,
    "longitude": -123.65306829619516,
    "inventories": [
        {
            "item": 1,
            "quantity": 200
        },
        {
            "item": 2,
            "quantity": 300
        },
        {
            "item": 3,
            "quantity": 350
        },
        {
            "item": 4,
            "quantity": 2000
        }
    ]
}
```
* POST `/users/new-token` -> Since there's no authentication as such, we allow users to get a new token with their emails. Expected Payload is:
```json
{
    "email":"tolaabbey009@carroll.net"
}
```

* GET `/users/me` -> Returns the user's information with balances and a new token that can be used to make future requests
* POST `/users/flag` -> Creates a new flag for the given `infectedUserID`. The expected payload is:
```json
{
    "infected_user_id":"user_Id"
}
```
* PATCH `/users/location` -> Allows users to updates their location. Users are detected with their auth token. Payload is:
```json
{
    "latitude": 6.5,
    "longitude": 3.5
}
```

* POST `/trades/initiate` -> Initiates the trade. The originating user is detected via the auth token. Payload:
```json
{
    "originator": {
        "items": [
            {
                "item": 1,
                "quantity": 1
            },
            {
                "item": 3,
                "quantity": 1
            }
        ]
    },
    "second_party": {
        "userID": "fcc9da93-46ce-4230-8b51-c5b8fc7e04fc",
        "reference": "",
        "items": [
            {
                "item": 4,
                "quantity": 6
            }
        ]
    }
}
```
`second_party` contains the details of the receiving party on the other side of the trade. This returns a reference ID and the inventory balance for the user.

* GET `/reports/survivor` -> returns the total number of survivors (`total_survivors`), total currently clean (`clean`) and percentage of clean survivors (`percentage_clean`)
```json
{
    "total_survivors": 2,
    "clean": 2,
    "percentage_clean": 100
}
```

* GET `/reports/infected` -> returns the total number of survivors (`total_survivors`), total of currently infected survivors (`infected_survivors`) and percentage of infected survivors (`percentage_infected`)
```json
{
    "total_survivors": 2,
    "infected_survivors": 0,
    "percentage_infected": 0
}
```

* GET `/reports/lost-point` -> returns the sum of all the lost inventories from infected survivors (data). and a success flag to determine if the request went well, as `0` can either mean there's no lost point or an error occurred.
```json
{
    "data": 0,
    "success": true
}
```

* GET `/reports/resources` -> returns the average of each inventory item that can be made available to each survivor rounded up to the nearest whole number.
```json
[
    {
        "item": "Medication",
        "balance": 697,
        "per_survivor": 348
    },
    {
        "item": "Water",
        "balance": 397,
        "per_survivor": 198
    },
    {
        "item": "Ammunition",
        "balance": 3982,
        "per_survivor": 1991
    },
    {
        "item": "Food",
        "balance": 600,
        "per_survivor": 300
    }
]
```

## Improvements

The trade endpoint is not idempotent, which means you can trigger a trade multiple times. To solve this, a failsafe/cooldown period might be deployed to make sure you don't execute the same trade with the same parameters to the same recipient within a period of time.

Transactions should be deployed to make the flow an "ACIDIC" flow. This will prevent inbalance in the inventory items.

Also, Message Queues and consumers can be introduced to make sure that the trades are cleaned up without the bottleneck of multiple inserts.