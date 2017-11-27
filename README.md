[![Build Status](https://travis-ci.org/mehrdadrad/goasn.svg?branch=master)](https://travis-ci.org/mehrdadrad/goasn)
[![Go Report Card](https://goreportcard.com/badge/github.com/mehrdadrad/goasn)](https://goreportcard.com/report/github.com/mehrdadrad/goasn)
# Resolve ASN to description
It works based on the [bgp.potaroo.net](http://bgp.potaroo.net/) and resolves the AS number to description. It creates a database at your local host and doesn't touch the website for the future resolve's requests. it's fast and safe for goroutine.

# Installation

     go get github.com/mehrdadrad/goasn

# Usage
```
  package main

  import (
      "github.com/mehrdadrad/goasn"
  )

  func main() {
      asn := goasn.NewASN()
      asn.Init()

      if err, a := asn.Get(15133); err != nil {
        println(a.Descr)
      }
  }
```
```
>go run main.go
>EDGECAST - MCI Communications Services, Inc. d/b/a Verizon Business, US
 ```
## License
This project is licensed under MIT license. Please read the LICENSE file.


## Contribute
Welcomes any kind of contribution, please follow the next steps:

- Fork the project on github.com.
- Create a new branch.
- Commit changes to the new branch.
- Send a pull request.
