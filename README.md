# Microservice for retrieving page with seo data and products

## Usage

### Prerequisites for running application

- Linux
- Network access (necessary for downloading dependencies)
- Ports 8080 available
- docker and docker-compose

### Running the application
To run application execute command:

```bash
make build-and-run
```

or if the docker image is already build:

```bash
make run
```

To load dummy data from resources/mongodb/ to database run:
```bash
make setup-dummy-db-data
```

### Using the API

#### */pages/{id}* endpoint
##### GET

Returns page with id given as path parameter

Sample request:
```bash
curl --request GET \
  --url http://localhost:8080/pages/1
```

Sample response:
```json
{
  "SEO": {
    "PageId": 1,
    "Title": "title1",
    "Description": "description1",
    "Robots": "robots1"
  },
  "Products": [
    {
      "Id": 2,
      "PageId": 1,
      "Name": "name2",
      "Description": "description2",
      "Price": 20.99
    },
    {
      "Id": 1,
      "PageId": 1,
      "Name": "name1",
      "Description": "description2",
      "Price": 2.5
    }
  ]
}
```
## Development

### Building project with tests

```bash
make test
```

### Building project without tests

```bash
make build
```

### Building docker image

Image will be built with tag: pages-ms:latest

```bash
make build-docker
```

## Next steps

- Add integration tests
- Add logging library
- Add liveness/readiness checks
- Add OpenAPI
- Improve error messages
- Improve context handling
