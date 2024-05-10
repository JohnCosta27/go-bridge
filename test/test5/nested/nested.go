package nested

import (
	"johncosta.tech/struct-to-types/test/test5/nested/morenested"
	"johncosta.tech/struct-to-types/test/test5/nested/nested"
)

type Nested struct {
	DoubleNested string

	MoreNested     morenested.Nested
	MyDoubleNested nested.DoubleNested
}
