package nested

import (
	"JohnCosta27/go-bridge/test/test5/nested/morenested"
	"JohnCosta27/go-bridge/test/test5/nested/nested"
)

type Nested struct {
	DoubleNested string

	MoreNested     morenested.Nested
	MyDoubleNested nested.DoubleNested
}
