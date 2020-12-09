// +build !js

package engine

import (
	"fmt"
	"log"
)

func prepareInterop() {

}

func getGlobalContext() *GlobalContext {
	return newGlobalContext()
}

func commencePanic(reason, message string) {
	log.Println("Beginning to panic")
	log.Println(fmt.Sprintf("Reason: %s", reason))
	log.Println("Message:")
	log.Fatal(message)
}

func frontEndCloseScene() {}
