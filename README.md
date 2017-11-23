# Resolve ASN to description

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
  
>go run main.go 
>EDGECAST - MCI Communications Services, Inc. d/b/a Verizon Business, US
 ```
