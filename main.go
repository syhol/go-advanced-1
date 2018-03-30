package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

func main() {
	usingCases()
	usingSelect()
	annonStruct()
	rot13DecoratedReader()
	doSomeShouting()
	mustDoSomeThings()
}

func usingCases() {
	var (
		zeros int
		even  int
		odd   int
		total int
	)

	numbers := []int{1, 3, 0, 5, 7, -2, 8, 9, 10, 10}

Loop:
	for _, i := range numbers {
		total++

		switch {
		case i < 0:
			continue Loop
		case i == 0:
			zeros++
		case i%2 == 0:
			even++
		default:
			odd++
		}
	}

	fmt.Printf("zeros: %d, even: %d, odd %d, total %d\n", zeros, even, odd, total)

}

func usingSelect() {
	c := make(chan string, 4) // String channel with a limit of 4

	var s string

Fill:
	for {
		select {
		case c <- "":
		default:
			break Fill
		}
	}

Infinite:
	for {
		select {
		case s = <-c:
			if s == "" {
				fmt.Println("found a string")
			} else {
				fmt.Println("END")
				break Infinite
			}
		default:
			fmt.Println("Channel empty, adding 2 blank and 1 end")
			c <- ""
			c <- ""
			c <- "END"
		}
	}
}

func annonStruct() {
	var people []struct {
		Name string
	}

	people = append(people, struct{ Name string }{Name: "Tim"})
	people = append(people, struct{ Name string }{Name: "Bob"})
	fmt.Printf("%+v\n", people)
}

type foo []string

func (f foo) floop() {
	return
}

func typeAssertions() {
	// var people []struct {
	// 	Name string
	// }

	myFoo := foo{}

	func(f interface{}) {
		switch f.(type) {
		case int:
			fmt.Println("I am int")
		case []string:
			fmt.Println("I am []string")
		// case foo:
		// 	fmt.Println("I am foo")
		case (interface {
			floop()
		}):
			fmt.Println("I am a flooper")
		}

		_, isFlooper := f.(interface {
			floop()
		})

		fmt.Println("Flooper status:", isFlooper)
	}(myFoo)
}

type rot13Reader struct {
	r io.Reader
}

func rot13(b byte) byte {
	var a, z byte
	switch {
	case 'a' <= b && b <= 'z':
		a, z = 'a', 'z'
	case 'A' <= b && b <= 'Z':
		a, z = 'A', 'Z'
	default:
		return b
	}
	return (b-a+13)%(z-a+1) + a
}

func (r *rot13Reader) Read(bs []byte) (int, error) {
	count, err := r.r.Read(bs)
	if err != nil {
		return count, err
	}
	for i := 0; i < count; i++ {
		bs[i] = rot13(bs[i])
	}
	return count, nil
}

func getCode() io.Reader {
	s := strings.NewReader("Lbh penpxrq gur pbqr!\n")
	return s
}

func decorateDecoder(s io.Reader) io.Reader {
	r := rot13Reader{s}
	return &r
}

func rot13DecoratedReader() {
	io.Copy(os.Stdout, getCode())
	io.Copy(os.Stdout, decorateDecoder(getCode()))

	r := decorateDecoder(getCode())
	// Custom 4 bytes at a time
	bs := make([]byte, 4, 4)
	var count int
	var err error
	for err != io.EOF {
		count, err = r.Read(bs)
		if len(bs) != count {
			bs = bs[:count]
		}
		os.Stdout.Write(bs)
	}
}

type myString struct {
	str string
}

type shouting struct {
	myString
}

func newMyString(s string) myString {
	return myString{str: s}
}

func newShouting(s string) shouting {
	shouting := shouting{}
	shouting.str = strings.ToUpper(s)
	return shouting
}

func (ms myString) Output() {
	fmt.Println(ms.str)
}

func (ms shouting) Output() {
	fmt.Println("Shouting:", ms.str)
}

func doSomeShouting() {
	newMyString("Yellow Word").Output()
	newShouting("Yellow Word").Output()
}

func mustDoSomeThings() {
	exp := regexp.MustCompile("[1-9]+")
	fmt.Println(exp.MatchString("Foo 1 Bar"))
	fmt.Println(exp.MatchString("Foo Bar"))
}
