package nested

import (
	"johncosta.tech/struct-to-types/test/test5/nested/morenested"
)

type Nested struct {
	DoubleNested string

	MoreNested morenested.Nested
}
