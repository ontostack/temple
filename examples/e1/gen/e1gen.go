package main

import "io"

func makeLoop(writer io.Writer) {

	@
	package main
	import "fmt"

	func main() {
	@

	for _, s := range []string{"Alice", "Bob"} {
		x := struct{y string}{y: s}
			@
			fmt.Println(#(x.y))
			fmt.Println(#s)
			@
	}

	@}@

}
