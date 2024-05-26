package nested

import (
	"github.com/JohnCosta27/go-bridge/test/test5/nested/morenested"
	"github.com/JohnCosta27/go-bridge/test/test5/nested/nested"
)

type Nested struct {
	DoubleNested string

	MoreNested     morenested.Nested
	MyDoubleNested nested.DoubleNested
}
