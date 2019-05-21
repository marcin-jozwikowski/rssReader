package cli

var verboseLevel = DefaultVerbose

const DefaultVerbose = levelNormal
const levelDebug = 3
const levelInfo = 2
const levelNormal = 1
const levelNone = 0

func IsVerbose() bool {
	return verboseLevel > levelNone
}

func IsVerboseInfo() bool {
	return verboseLevel >= levelInfo
}

func IsVerboseDebug() bool {
	return verboseLevel >= levelDebug
}

func SetVerbose(newVerboseLevel int) {
	verboseLevel = newVerboseLevel
}
