
go get github.com/rs/xid

curl --location --request POST '<http://localhost:8080/recipes>'
\
--header 'Content-Type: application/json' \
--data-raw '{
"name": "Homemade Pizza",
"tags" : ["italian", "pizza", "dinner"],
"ingredients": [
"1 1/2 cups (355 ml) warm water (105°F-115°F)",
"1 package (2 1/4 teaspoons) of active dry yeast",
"3 3/4 cups (490 g) bread flour",
"feta cheese, firm mozzarella cheese, grated"
],
"instructions": [

    "Step 1.",
"Step 2.",
"Step 3."
]
}


> swagger generate spec -o ./swagger.json
> swagger serve ./swagger.json
> swagger serve -F swagger ./swagger.json


# Commands
- go get go.mongodb.org/mongo-driver/v2/mongo

- MONGO_URI="mongodb://admin:password@localhost:27017/test?authSource=admin" MONGO_DATABASE=demo go run main.go

- mongoimport --username admin --password password --authenticationDatabase admin --db demo --collection recipes --file recipes.json --jsonArray