# Price Watcher Rest Service

price watcher is an application that allows users to pass a url to an item on ebay/amazon and get notified when the price
changes.
The backend api is written in Go with the Gin(api) and Gorm(database orm) libraries

## Routes

The base url for all routes is {site}/api

POST /user/login
login a user with an email and password

GET /items
returns a basic json list of items
eg.

```json
{
  "items": [
    {
      "name": "the items name",
      "nickname": "user defined item nickname",
      "current_price": "last checked price of the item",
      "id": "items id",
      "last_checked": "12/12/12"
    },
    ...
  ]
}
```

GET /items/{id}
returns a single item by id

```json
{
  "name": "the items name",
  "nickname": "user defined item nickname",
  "current_price": "last checked price of the item",
  "id": "items id",
  "url": "item url",
  "last_checked": "12/12/12",
  "next_check": "13/12/12",
  "starting_price": 12, // the price of the item when the user first added it
  "price_difference": "50" // the price difference between the starting price and the last checked price
}
```

POST /items/{user session id}
add a new item to watch, pass in the user session id returned from login

```json
{
  "url": "product url"
}
```
