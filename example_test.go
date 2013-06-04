package gomruby_test

import (
	"."
	"fmt"
)

func Example() {
	mruby := gomruby.New()
	context := mruby.NewLoadContext("filter.rb")

	userCode := `
		message.include? "500"
	`

	code := fmt.Sprintf(`
		def filter
			message = ARGV[0]
			%s
		end`, userCode)
	_, err := context.Load(code)
	if err != nil {
		panic(err)
	}

	for _, message := range []string{
		`1.2.3.1 - - [04/Jun/2013:18:02:01 +0000] host "GET /foo HTTP/1.0" 200`,
		`1.2.3.2 - - [04/Jun/2013:18:02:02 +0000] host "GET /bar HTTP/1.0" 300`,
		`1.2.3.3 - - [04/Jun/2013:18:02:03 +0000] host "GET /baz HTTP/1.0" 400`,
		`1.2.3.4 - - [04/Jun/2013:18:02:04 +0000] host "GET /bzr HTTP/1.0" 500`,
	} {
		res, err := context.Load("filter", message)
		if err != nil {
			panic(err)
		}
		fmt.Println(res.(bool))
	}

	// Output:
	// false
	// false
	// false
	// true
}
