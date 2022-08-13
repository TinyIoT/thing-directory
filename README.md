# TinyIoT Thing Directory
This is an implementation of the [W3C Web of Things (WoT) Thing Description Directory (TDD)](https://w3c.github.io/wot-discovery/), a registry of [Thing Descriptions](https://www.w3.org/TR/wot-thing-description/).

## Getting Started
Visit the following pages to get started:
* [Deployment](../../wiki/Deployment): How to deploy the software, as Docker container, Debian package, or platform-specific binary distributions
* [Configuration](../../wiki/Configuration): How to configure the server software with JSON files and environment variables
* [API Documentation][1]: How to interact with the networking APIs

**Further documentation are available in the [wiki](../../wiki)**.

## Features
* Service Discovery
  * [DNS-SD registration](../../wiki/Discovery-with-DNS-SD)
* RESTful API
  * [HTTP API][1]
    * Things API - TD creation, read, update (put/patch), deletion, and listing (pagination) 
    * Search API - [JSONPath query language](../../wiki/Query-Language)
    * Events API
    * TD validation with JSON Schema(s)
    * Request [authentication](https://github.com/linksmart/go-sec/wiki/Authentication) and [authorization](https://github.com/linksmart/go-sec/wiki/Authorization)
    * JSON-LD response format
* Persistent Storage
  * LevelDB
* CI/CD ([Github Actions](https://github.com/tinyiot/thing-directory/actions?query=workflow:CICD))
  * Automated testing
  * Automated builds and releases

## Development

Clone this repo:
```bash
git clone https://github.com/tinyiot/thing-directory.git
cd thing-directory
```

Compile from source:
```bash
go build
```
This will result in an executable named `thing-directory` (linux/macOS) or `thing-directory.exe` (windows).

Get the CLI arguments help (linux/macOS):
```bash
$ ./thing-directory -help
Usage of ./thing-directory:
  -conf string
        Configuration file path (default "conf/thing-directory.json")
  -version
        Print the API version
```

Run (linux/macOS):
```bash
$ ./thing-directory --conf=sample_conf/thing-directory.json
```

To build and run together:
```bash
go run . --conf=sample_conf/thing-directory.json
```

Test all packages (add `-v` flag for verbose results):
```bash
go test ./...
```


## Contributing
Contributions are welcome. 

Please fork, make your changes, and submit a pull request. For major changes, please open an issue first and discuss it with the other authors.


[1]: https://petstore.swagger.io?url=https://raw.githubusercontent.com/tinyiot/thing-directory/master/apidoc/openapi-spec.yml
