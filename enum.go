package main

import "fmt"

type ServerState int

const (
	UNDEPLOYED ServerState = iota
	IN_PROGRESSING
	DEPLOYED
)

func validate(state ServerState) error {
	if state == UNDEPLOYED || state == IN_PROGRESSING || state == DEPLOYED {
		return nil
	}
	return fmt.Errorf("%+v is not a valid state", state)
}

func main() {
	const NEW_STATE ServerState = 99
	const HACKING_STATE int = 2

	// IN_PROGRESSING works well, because it is in the enum
	if err := validate(IN_PROGRESSING); err != nil {
		fmt.Printf("%+v", err)
	}

	// NEW_STATE works well, because it is a var of type ServerState
	if err := validate(NEW_STATE); err != nil {
		fmt.Printf("%+v", err)
	}

	//  the following code wont compile, because it is not a var of type  ServerState
	//  this is due to the use of 'iota', compiler will ensure only specified enums are used
	if err := validate(HACKING_STATE); err != nil {
		fmt.Printf("%+v", err)
	}
}
