# Bestiary-Crud

This is a simple CRUD app

This application is intended to store information about beasts in DnD 5e.

## App

### Endpoints

Operations may be performed using the following defined HTTP endpoints:

#### GET /

Health check endpoint.

Request:

```http
GET http://localhost:8080/
```

Response:

```json
{"status": "ok"}
```

#### GET /beasts

This endpoint will return items in the psql db.

Request:

```http
GET http://localhost:8080/beasts
```

Response:

```json
[
    {
        "BeastName": "Mind Flayer",
        "Type": "Aberration",
        "CR": "7",
        "Data": {
            "Attributes": {
                "CHA": "17 (+3)",
                "CON": "12 (+1)",
                "DEX": "12 (+1)",
                "INT": "19 (+4)",
                "STR": "11 (+0)",
                "WIS": "17 (+3)"
            },
            "Description": "Innate Spellcasting (Psionics). The mind flayer's innate spellcasting ability is Intelligence (spell save DC 15). It can innately cast the following spells, requiring no components:\nAt will: detect thoughts, levitate\n1/day each: dominate monster, plane shift (self only)\nMagic Resistance. The mind flayer has advantage on saving throws against spells and other magical effects."
        }
    },
    {
        "BeastName": "Mimic",
        "Type": "Monstrosity(Shapechanger)",
        "CR": "2",
        "Data": {
            "Attributes": {
                "CHA": "8 (-1)",
                "CON": "15 (+2)",
                "DEX": "12 (+1)",
                "INT": "5 (-3)",
                "STR": "17 (+3)",
                "WIS": "13 (+1)"
            },
            "Description": "Shapechanger. The mimic can use its action to polymorph into an object or back into its true, amorphous form. Its statistics are the same in each form. Any equipment it is wearing or carrying isn't transformed. It reverts to its true form if it dies.\nGrappler. The mimic has advantage on attack rolls against any creature grappled by it.\nAdhesive (Object Form Only). The mimic adheres to anything that touches it. A Huge or smaller creature adhered to the mimic is also grappled by it (escape DC 13). Ability checks made to escape this grapple have disadvantage.\nFalse Appearance (Object Form Only). While the mimic remains motionless, it is indistinguishable from an ordinary object."
        }
    }
]
```

#### POST /beasts

This endpoint will store given beast data in psql and return the object id.

Request

```http
POST http://localhost:8080/beasts

{
    "BeastName": "Mimic",
    "Type": "Monstrosity(Shapechanger)",
    "CR": "2",
    "Attributes": {
        "CHA": "8 (-1)",
        "CON": "15 (+2)",
        "DEX": "12 (+1)",
        "INT": "5 (-3)",
        "STR": "17 (+3)",
        "WIS": "13 (+1)"
    },
    "Description": "Shapechanger. The mimic can use its action to polymorph into an object or back into its true, amorphous form. Its statistics are the same in each form. Any equipment it is wearing or carrying isn't transformed. It reverts to its true form if it dies.\nGrappler. The mimic has advantage on attack rolls against any creature grappled by it.\nAdhesive (Object Form Only). The mimic adheres to anything that touches it. A Huge or smaller creature adhered to the mimic is also grappled by it (escape DC 13). Ability checks made to escape this grapple have disadvantage.\nFalse Appearance (Object Form Only). While the mimic remains motionless, it is indistinguishable from an ordinary object."
}
```

Response:

```json
{
    "BeastName": "Mimic"
}
```

#### GET /beasts/{key}

This endpoint will return the object of the given key in JSON format.

Request:

```http
GET /beasts/get/Mimic
```

Response:

```json
{
    "BeastName": "Mimic",
    "Type": "Monstrosity(Shapechanger)",
    "CR": "2",
    "Attributes": {
        "CHA": "8 (-1)",
        "CON": "15 (+2)",
        "DEX": "12 (+1)",
        "INT": "5 (-3)",
        "STR": "17 (+3)",
        "WIS": "13 (+1)"
    },
    "Description": "Shapechanger. The mimic can use its action to polymorph into an object or back into its true, amorphous form. Its statistics are the same in each form. Any equipment it is wearing or carrying isn't transformed. It reverts to its true form if it dies.\nGrappler. The mimic has advantage on attack rolls against any creature grappled by it.\nAdhesive (Object Form Only). The mimic adheres to anything that touches it. A Huge or smaller creature adhered to the mimic is also grappled by it (escape DC 13). Ability checks made to escape this grapple have disadvantage.\nFalse Appearance (Object Form Only). While the mimic remains motionless, it is indistinguishable from an ordinary object."
}
```

#### DELETE /beasts/{key}

This endpoint will delete the object of the given key

Request:

```http
DELETE /beasts/get/Mimic
```

Response:

```json
{
    "message": "Beast 'Mimic' deleted successfully."
}
```

#### PUT /beasts/{key}

This endpoint will update the object of the given key.

Request:

```http
PUT http://localhost:8080/beasts/Mimic

{
    "BeastName": "Mimic",
    "Type": "Monstrosity(Shapechanger)",
    "CR": "3",
    "Attributes": {
        "CHA": "8 (-1)",
        "CON": "16 (+3)",
        "DEX": "12 (+1)",
        "INT": "5 (-3)",
        "STR": "17 (+3)",
        "WIS": "13 (+1)"
    },
    "Description": "Shapechanger. The mimic can use its action to polymorph into an object or back into its true, amorphous form. Its statistics are the same in each form. Any equipment it is wearing or carrying isn't transformed. It reverts to its true form if it dies.\nGrappler. The mimic has advantage on attack rolls against any creature grappled by it.\nAdhesive (Object Form Only). The mimic adheres to anything that touches it. A Huge or smaller creature adhered to the mimic is also grappled by it (escape DC 13). Ability checks made to escape this grapple have disadvantage.\nFalse Appearance (Object Form Only). While the mimic remains motionless, it is indistinguishable from an ordinary object."
}
```

Response:

```json
{
    "message": "Beast 'Mimic' updated successfully."
}
```

### Code structure

- /api: Contains the db client and handlers for each endpoint.
- /config: Contains the config reading functions and the config files themselves in yaml format.
- /tests: Contains the unit tests for the testing stage. Has its own config.

### Running for development

```Shell
go run ./main.go
```

## Deployment

The application is fully dockerized and utilizes multi-stage docker builds for minimal release image footprint.
The final image runs on alpine 3.19. A minimal `docker-compose.yml` is also provided for local development.

### GitHub Actions

#### Build and test

Upon any merge action to the main branch, two pipelines trigger:

Build-and-test:

- Builds an image and tags it
- Runs tests via the test target of the multi stage build
- Pushes the container to AWS ECR

Deploy:

- Builds, tags and pushes the image
- Updates the task definition
- Deploys the task definition
