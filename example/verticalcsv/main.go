package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/ibbbpbbbp/gotinycsv"
)

func main() {
	CSV := `header,A,B,C,D,E,F,G,H,I,J,K,L,M,N,O,P,Q,R
No,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18
Name,Alex,Bert,Carl,Dave,Elly,Fran,Gwen,Hank,Ivan,Jake,Kate,Luke,Myra,Neil,Omar,Page,Quin,Ruth
Sex,M,M,M,M,F,F,F,M,M,M,F,M,F,M,M,F,M,F
Age,41,42,32,39,30,33,26,30,53,32,47,34,23,36,38,31,29,28
Height(in),74,68,70,72,66,66,64,71,72,69,69,72,62,75,70,67,71,65
Weight(lbs),170,166,155,167,124,115,121,158,175,143,139,163,98,160,145,135,176,131
Birth,1999.01.01,2001.02.02,2002.05.05,1999.06.06,2003.03.03,1999.03.03,1998.04.04,2004.07.07,1999.11.11,1997.12.12,2002.10.10,1999.06.06,1995.09.09,1999.08.08,2001.03.03,1998.02.02,1997.04.04,1994.05.05`
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
	const leftmergin = 1
	const maxcols = 100
	if err := gotinycsv.LoadVertically(strings.NewReader(CSV), topmergin, leftmergin, maxcols, &perssonal); err != nil {
		fmt.Printf("%#v\n", err)
	}
	for _, v := range perssonal {
		fmt.Printf("Name:%s Age:%d Height:%f Weight:%f Birth:%s\n", v.Name, v.Age, v.Height, v.Weight, v.BirthDate.String())
	}
}
