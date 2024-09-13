package nested

import "github.com/JohnCosta27/go-bridge/test/test7/nested/morenested"

type B struct {
	Hello string
	World struct {
		C morenested.D
	}
}
