package conf

import (
	"io/ioutil"
	"log"
	"os"
)

// Version is a string of the currently running application version
var Version string = getVersion()

// Get the application version from version.txt file or GIT_VERSION env variable if not present
func getVersion() string {
	data, err := ioutil.ReadFile("version.txt")
	if err != nil {
		log.Println("version.txt file not found, fallback to GIT_VERSION env variable")
		return os.Getenv("GIT_VERSION")
	}
	return string(data)
}
