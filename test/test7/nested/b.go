package nested

import "johncosta.tech/go-bridge/test/test7/nested/morenested"

type B struct {
	World struct {
		C morenested.D
	}
}
