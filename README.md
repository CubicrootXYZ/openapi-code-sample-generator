# OpenApi code sample generator

[![GitHub license](https://img.shields.io/github/license/CubicrootXYZ/openapi-code-sample-generator)](https://github.com/CubicrootXYZ/openapi-code-sample-generator/blob/main/LICENSE)
[![GitHub issues](https://img.shields.io/github/issues/CubicrootXYZ/openapi-code-sample-generator)](https://github.com/CubicrootXYZ/openapi-code-sample-generator/issues)
[![Actions Status](https://github.com/CubicrootXYZ/openapi-code-sample-generator/workflows/Main/badge.svg?branch=main)](https://github.com/CubicrootXYZ/openapi-code-sample-generator/workflows/actions)

Generates sample code for different languages and inserts them in you OpenApi specification file as `x-codeSamples`.

Tools like [redoc](https://github.com/Redocly/redoc) will provide those examples in the documentation.

## 📋 Supported Languages

* curl

### Supported Content-Types

* application/xml
* text/xml
* application/json
* text/json
* application/x-www-form-urlencoded
* multipart/form-data

### Supported authentication

* basic
* apikey
* cookie
* header
* openidconnect
* oauth2

## 🔍 Example

```
/pet/findByTags:
    get:
      ...
      x-codeSamples:
      - label: curl
        lang: curl
        source: 'curl "http://petstore.swagger.io/v2/pet/findByTags?tags%5B%5D=example-string"
          -H "Authorization: Bearer ${TOKEN}" -d "" -X GET'
```

## 📥 Installation

```
$ go install github.com/CubicrootXYZ/openapi-code-sample-generator@latest
```

or chose a version from the Github releases:

```
$ go install github.com/CubicrootXYZ/openapi-code-sample-generator@v0.0.1
```

## 📚 Usage

Depending on your environment the module is now available with
```
openapi-code-sample-generator
```

or you need to search for the binary, those are common places: 

```
/home/$user/go/bin/openapi-code-sample-generator
$GOPATH/bin
```

To add code samples: 

```
$ openapi-code-sample-generator generate --input-file example.yaml --output-file out.yaml
```

To convert OpenAPI 2 specifications to OpenAPI 3: 

```
$ openapi-code-sample-generator convert --file example.yaml --output-file out.yaml
```

For getting help:

```
$ openapi-code-sample-generator --help
$ openapi-code-sample-generator generate --help
$ openapi-code-sample-generator convert --help
```

## 🛠 Troubleshooting

### Can not install the module

Try it without the `@latest` version tag - older go versions might not be able to work with it.

Update your go version to go 1.16.

### Can not execute the module

`go install` usually builds a binary in `$GOPATH/bin` if you can not find it there you might just search the whole file system for it:

```
$ cd / && find * | grep openapi-code-sample-generator
```

Will show you all directories and files with the modules name. 

You then need to execute this binary, e.g.:

```
$ /home/michael/go/bin/openapi-code-sample-generator generate --help
```

Or add it to your `PATH` variable.

## 👨‍👨‍👧 Contributing

I really enjoy any help and contribution to this project. Feel free to work on open issues or new features. 

### Issues, ideas and more

Please submit your issues or specific feature requests as "Issues". Be as precise as possible.

### Contributing code

Fork this repository and add your changes. Open a pull request to merge them in the master branch of this repository.


## ❤️ Attribution

Great thanks to the maintainers of the used libraries and all contributors.

* [anyxml](https://github.com/clbanning/anyxml)
* [kin-openapi](https://github.com/getkin/kin-openapi)
* [ghodss/yaml](https://github.com/ghodss/yaml)
* [gorilla/schema](https://github.com/gorilla/schema)
* [cobra](https://github.com/spf13/cobra)