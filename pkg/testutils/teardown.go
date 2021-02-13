package testutils

type Teardown func()

// CombineInto returns one Teardown function, which applies tds in order
//goland:noinspection SpellCheckingInspection
func CombineInto(tds ...Teardown) Teardown {
	return func() {
		for _, teardown := range tds {
			teardown()
		}
	}
}

// CombineInto returns one Teardown function, which applies {first: Teardown} first
// and then {other []Teardown} in order
// remarks: note that {first: Teardown} does not get updated!
func (first Teardown) CombineInto(other ...Teardown) Teardown {
	tdFuncs := append([]Teardown{first}, other...)
	return CombineInto(tdFuncs...)
}
