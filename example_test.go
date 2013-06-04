package gomruby_test

import (
	"."
	"fmt"
)

func Example() {
	// create new VM instance and load context
	mruby := gomruby.New()
	defer mruby.Delete()
	context := mruby.NewLoadContext("select.rb")
	defer context.Delete()

	// this is user-supplied code
	userCode := `
		message.include? "500"
	`

	// define method
	code := fmt.Sprintf(`
		def select
			message = ARGV[0]
			%s
		end`, userCode)
	_, err := context.Load(code)
	if err != nil {
		panic(err)
	}

	// iterate over messages and select interesting
	for _, message := range []string{
		`1.2.3.1 - - [04/Jun/2013:18:02:01 +0000] host "GET /foo HTTP/1.0" 200`,
		`1.2.3.2 - - [04/Jun/2013:18:02:02 +0000] host "GET /bar HTTP/1.0" 300`,
		`1.2.3.3 - - [04/Jun/2013:18:02:03 +0000] host "GET /baz HTTP/1.0" 400`,
		`1.2.3.4 - - [04/Jun/2013:18:02:04 +0000] host "GET /bzr HTTP/1.0" 500`,
	} {
		res, err := context.Load("select", message)
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
