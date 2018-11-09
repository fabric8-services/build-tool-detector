/*

Package types implements a simple way to
create a GoaBuildToolDetector struct with
the recognized build type.

The GoaBuildToolDetector struct will result
in the following JSON response:

	[[ maven ]]

		{
		"build-tool-type": "maven"
		}


	[[ unknown ]]

		{
		"build-tool-type": "unknown"
		}

*/
package types

import (
	"github.com/fabric8-services/build-tool-detector/app"
)

const (
	// Maven build type detected Maven.
	Maven  = "maven"
	pomXML = "pom.xml"

	// NodeJS build type detected node.
	NodeJS      = "nodejs"
	packageJSON = "package.json"

	// Unknown build type detected Unknown.
	Unknown = "unknown"
)

// BuildType TODO
type BuildType struct {
	BuildType string
	File      string
}

// NewMaven will create a detect
// struct with the BuildToolType set
// to maven.
func NewMaven() *app.GoaDetect {
	return &app.GoaDetect{
		BuildToolType: Maven,
	}
}

// NewNodeJS will create a detect
// struct with the BuildToolType set
// to NodeJS.
func NewNodeJS() *app.GoaDetect {
	return &app.GoaDetect{
		BuildToolType: NodeJS,
	}
}

// NewUnknown will create a detect
// struct with the BuildToolType set
// to unknown.
func NewUnknown() *app.GoaDetect {
	return &app.GoaDetect{
		BuildToolType: Unknown,
	}
}

// GetTypes returns the BuilType for all
// supported build tools.
func GetTypes() []BuildType {
	buildTypes := make([]BuildType, 2)

	buildTypes[0] = getTypeMaven()
	buildTypes[1] = getTypeNodeJS()

	return buildTypes
}

// getTypeMaven returns BuildType for maven.
func getTypeMaven() BuildType {
	return BuildType{Maven, pomXML}
}

// getTypeNodeJS returns BuildType for nodejs.
func getTypeNodeJS() BuildType {
	return BuildType{NodeJS, packageJSON}
}
