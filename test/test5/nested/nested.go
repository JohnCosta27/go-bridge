package nested

import (
	"johncosta.tech/go-bridge/test/test5/nested/morenested"
	"johncosta.tech/go-bridge/test/test5/nested/nested"
)

type Nested struct {
	DoubleNested string

	MoreNested     morenested.Nested
	MyDoubleNested nested.DoubleNested
}
