# Web crawler

This is my _go_ at creating a web crawler using Golang.

## Configure
The following elements are accepted as environment variables in `./settings.yaml`
```yaml
baseURL: "https://google.com" // The site to be crawled
printerType: "json" // The desired format of the results ["raw","json"]
persist: true // If you wish for the results to be written to a file
httpTimeout: 10s // The time to wait for the HTTP Client before returning an error
```

## Build
This will lint the codebase as well create a binary.
```makefile
make
```

## Execute
After building the binary, start it by using the following
```
./crawler serve
```

## Design
This project was designed with
- [Twelve-Factor](https://12factor.net/) in mind
- A commonly adopted project [structure](https://github.com/golang-standards/project-layout)
- Golang's official code review [guide](https://github.com/golang/go/wiki/CodeReviewComments)

