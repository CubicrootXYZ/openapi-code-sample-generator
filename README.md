# OpenApi code sample generator

Generates sample code for different languages and inserts them in you OpenApi specification file.

## Supported Languages

* curl

## Install

1. `git clone https://github.com/CubicrootXYZ/openapi-code-sample-generator.git`
2. `cd openapi-code-sample-generator`
3. `go install openapi-code-sample-generator`

Now you have installed the package. If it is automaticall added to your path you can use it with:
```
openapi-code-sample-generator
```

Otherwise you might find the binary in:
```
/home/$user/go/bin/openapi-code-sample-generator
```

## Use

To add samples: 
```
openapi-code-sample-generator generate --input-file example.yaml --output-file index.html
```