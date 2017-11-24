[![Build Status](https://travis-ci.org/mehrdadrad/goasn.svg?branch=master)](https://travis-ci.org/mehrdadrad/goasn)
# Resolve ASN to description

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
  
      a := asn.Get(15133)
  
      println(a.Descr)
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
