# gotinycsv

## WHAT
An extremely small `Go` library that efficiently deserializes CSV format data.  
Support for special CSV with vertically aligned fields.

## Example
```go
package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/ibbbpbbbp/gotinycsv"
)

func main() {
	CSV := `No,Name,Sex,Age,Height(in),Weight(lbs),Birth
1,Alex,M,41,74,170,1999.01.01
2,Bert,M,42,68,166,2001.02.02
3,Carl,M,32,70,155,2002.05.05
4,Dave,M,39,72,167,1999.06.06
5,Elly,F,30,66,124,2003.03.03
`
	// slices do not need to be pre-allocated, gotinycsv will ensure optimal size.
	perssonal := []*struct {
		_         int64 // No (igonore)
		Name      string
		_         string // Sex (ignore)
		Age       int64
		Height    float64
		Weight    float64
		BirthDate time.Time
	}{}

	const topmergin = 1
	const maxrows = 100
	if err := gotinycsv.Load(strings.NewReader(CSV), topmergin, maxrows, &perssonal); err != nil {
		fmt.Printf("%#v\n", err)
	}
	for _, v := range perssonal {
		fmt.Printf("Name:%s Age:%d Height:%f Weight:%f Birth:%s\n", v.Name, v.Age, v.Height, v.Weight, v.BirthDate.String())
	}
}
```
> output
```shell
Name:Alex Age:41 Height:74.000000 Weight:170.000000 Birth:1999-01-01 00:00:00 +0000 UTC
Name:Bert Age:42 Height:68.000000 Weight:166.000000 Birth:2001-02-02 00:00:00 +0000 UTC
Name:Carl Age:32 Height:70.000000 Weight:155.000000 Birth:2002-05-05 00:00:00 +0000 UTC
Name:Dave Age:39 Height:72.000000 Weight:167.000000 Birth:1999-06-06 00:00:00 +0000 UTC
Name:Elly Age:30 Height:66.000000 Weight:124.000000 Birth:2003-03-03 00:00:00 +0000 UTC
```
Support for special CSVs with vertically aligned fields.  
I don't know if there is a demand for it :-)
```go
package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/ibbbpbbbp/gotinycsv"
)

func main() {
	CSV := `header,A,B,C,D,E
No,1,2,3,4,5
Name,Alex,Bert,Carl,Dave,Elly
Sex,M,M,M,M,F
Age,41,42,32,39,30
Height(in),74,68,70,72,66
Weight(lbs),170,166,155,167,124
Birth,1999.01.01,2001.02.02,2002.05.05,1999.06.06,2003.03.03
`
	// slices do not need to be pre-allocated, gotinycsv will ensure optimal size.
	perssonal := []struct {
		_         int64 // No (igonore)
		Name      string
		_         string // Sex (ignore)
		Age       int64
		Height    float64
		Weight    float64
		BirthDate time.Time
	}{}

	const topmergin = 1
	const leftmergin = 1
	const maxrows = 100
	if err := gotinycsv.LoadVertically(strings.NewReader(CSV), topmergin, leftmergin, maxrows, &perssonal); err != nil {
		fmt.Printf("%#v\n", err)
	}
	for _, v := range perssonal {
		fmt.Printf("Name:%s Age:%d Height:%f Weight:%f Birth:%s\n", v.Name, v.Age, v.Height, v.Weight, v.BirthDate.String())
	}
}
```
> output
```shell
Name:Alex Age:41 Height:74.000000 Weight:170.000000 Birth:1999-01-01 00:00:00 +0000 UTC
Name:Bert Age:42 Height:68.000000 Weight:166.000000 Birth:2001-02-02 00:00:00 +0000 UTC
Name:Carl Age:32 Height:70.000000 Weight:155.000000 Birth:2002-05-05 00:00:00 +0000 UTC
Name:Dave Age:39 Height:72.000000 Weight:167.000000 Birth:1999-06-06 00:00:00 +0000 UTC
Name:Elly Age:30 Height:66.000000 Weight:124.000000 Birth:2003-03-03 00:00:00 +0000 UTC
```
## Support Type
The types supported by `out interface{}`, the argument of `Load() or LoadVertically()`, are follows.   

```go
out = []struct{T} | []*struct{T}
T = string | int | int8 | int16 | int32 | int64 | float32 | float64 | time.Time
```
