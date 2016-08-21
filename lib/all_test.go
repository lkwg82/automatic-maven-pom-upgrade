package lib

var temporaryDirectoryForTests = &TemporaryDirectoryForTests{}

func init() {
	temporaryDirectoryForTests.Init()
}

func setup() {
	temporaryDirectoryForTests.Setup()
}

func cleanup() {
	temporaryDirectoryForTests.Cleanup()
}
