package banner

import "fmt"

var Banner = `
  _________  ___  __  ____________ 
 / ___/ __ \/ _ |/ / / /_  __/ __ \
/ (_ / /_/ / __ / /_/ / / / / /_/ /
\___/\____/_/ |_\____/ /_/  \____/

		github.com/liancccc/goauto

`

func Print() {
	fmt.Print(Banner)
}
