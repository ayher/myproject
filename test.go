package main

import "fmt"

func main() {
	var a int = 1
	var b *int = &a
	var c **int = &b
	var x int = *b
	fmt.Println("a = 1",a) //1
	fmt.Println("&a = 0x1",&a) // 0x???
	fmt.Println("*&a = 1",*&a)	//1
	fmt.Println("b = 0x1",b)	//1
	fmt.Println("&b = 0x2",&b) // 0x??
	fmt.Println("*&b = 0x1",*&b) //1
	fmt.Println("*b = ",*b) // 1-??
	fmt.Println("c = ",c)// &b 0x??
	fmt.Println("*c = ",*c)// &c 1
	fmt.Println("&c = ",&c)// &c 0x???
	fmt.Println("*&c = ",*&c)// &b 0x???
	fmt.Println("**c = ",**c)// &a 0x??
	fmt.Println("***&*&*&*&c = ",***&*&*&*&*&c)//
	fmt.Println("x = ",x) // &a
}