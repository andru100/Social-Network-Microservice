package test

import (
	"fmt"
	//"time"
)


type inter1 interface {
	method1 () int 
	method2 () string
}

type testtype struct {
	answer int
}

type testtype2 struct {
	answer string
}

func (t *testtype) method1() int {
	return t.answer
}

func (t *testtype2) method2() string {
	return t.answer
}

func testinterface(i inter1) {
	fmt.Println(i.method1())
	fmt.Println(i.method2())
}

func (t *testtype) chainmeth() *testtype {
	t.answer = t.answer + 1
	return t
}

func (t *testtype) chainmeth2() *testtype {
	t.answer = t.answer + 1
	return t
}

func Test(i int) int {
	var t testtype
	// var t2 testtype2
	t.answer = i
	//t.answer = t.answer + 1
	//t.answer = t.answer + 1
	// t2.answer = "Hello World"
	// fmt.Println(t.method1())
	// fmt.Println(t2.method2())
	// testinterface(&t)
	// testinterface(&t2)
	return t.chainmeth().chainmeth2().answer
	//return t.answer

}

// interface.. types have to have a version of each method in the interface to satisfy the interface