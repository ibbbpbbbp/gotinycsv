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
6,Fran,F,33,66,115,1999.03.03
7,Gwen,F,26,64,121,1998.04.04
8,Hank,M,30,71,158,2004.07.07
9,Ivan,M,53,72,175,1999.11.11
10,Jake,M,32,69,143,1997.12.12
11,Kate,F,47,69,139,2002.10.10
12,Luke,M,34,72,163,1999.06.06
13,Myra,F,23,62,98,1995.09.09
14,Neil,M,36,75,160,1999.08.08
15,Omar,M,38,70,145,2001.03.03
16,Page,F,31,67,135,1998.02.02
17,Quin,M,29,71,176,1997.04.04
18,Ruth,F,28,65,131,1994.05.05
`
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
