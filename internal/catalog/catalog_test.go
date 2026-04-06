package catalog

import (
	"testing"

	"github.com/arduino/go-paths-helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLoad provides coverage for the `Load` function.
func TestLoad(t *testing.T) {
	workingDirectory, err := paths.Getwd()
	require.NoError(t, err)

	catalogPath := workingDirectory.Join("testdata", "TestLoad", "inoplatforms.tsv")
	catalogData, err := Load(catalogPath)
	require.NoError(t, err)

	assertion := Type{
		{
			"ADLINK SAMD Boards",
			"adlink",
			"samd",
			"https://github.com/ADLINK/OT2IT",
			"https://raw.githubusercontent.com/ADLINK/OT2IT/refs/heads/main/package_adlink_ot2it.json",
			"/",
			"master",
			"https://github.com/ADLINK/OT2IT",
			"/",
			"main",
			"",
			"Hard fork of https://github.com/arduino/ArduinoCore-samd.",
			"",
		},
		{
			"Tlera Corp STM32L0 Boards",
			"A21Y",
			"stm32l0",
			"https://github.com/OS-Q/A21Y",
			"",
			"/",
			"master",
			"",
			"",
			"",
			"",
			"Hard fork of https://github.com/GrumpyOldPizza/ArduinoCore-stm32l0.",
			"",
		},
	}

	assert.Equal(t, assertion, catalogData)
}

// TestTypeWrite provides coverage for the `(Type) Write` method.
func TestTypeWrite(t *testing.T) {
	data := Type{
		{
			"ADLINK SAMD Boards",
			"adlink",
			"samd",
			"https://github.com/ADLINK/OT2IT",
			"https://raw.githubusercontent.com/ADLINK/OT2IT/refs/heads/main/package_adlink_ot2it.json",
			"/",
			"master",
			"https://github.com/ADLINK/OT2IT",
			"/",
			"main",
			"",
			"Hard fork of https://github.com/arduino/ArduinoCore-samd.",
			"",
		},
		{
			"Tlera Corp STM32L0 Boards",
			"A21Y",
			"stm32l0",
			"https://github.com/OS-Q/A21Y",
			"",
			"/",
			"master",
			"",
			"",
			"",
			"",
			"Hard fork of https://github.com/GrumpyOldPizza/ArduinoCore-stm32l0.",
			"",
		},
	}

	writeFolder, err := paths.TempDir().MkTempDir("ino-platform-discovery-TestTypeWrite")
	require.NoError(t, err)

	writePath := writeFolder.Join("inoplatforms.csv")

	err = data.Write(writePath)
	require.NoError(t, err, "Function does not return an error.")

	exists, err := writePath.ExistCheck()
	require.NoError(t, err)
	require.True(t, exists, "The file was created.")

	actual, err := writePath.ReadFileAsLines()
	require.NoError(t, err)

	assertion := []string{
		"Name	Vendor	Architecture	Repository	Boards Manager URL	Repository Data Folder	Branch Name	Package Index Repository	Package Index Folder	Package Index Branch	Reference	Notes	Suppress",
		"ADLINK SAMD Boards	adlink	samd	https://github.com/ADLINK/OT2IT	https://raw.githubusercontent.com/ADLINK/OT2IT/refs/heads/main/package_adlink_ot2it.json	/	master	https://github.com/ADLINK/OT2IT	/	main		Hard fork of https://github.com/arduino/ArduinoCore-samd.	",
		"Tlera Corp STM32L0 Boards	A21Y	stm32l0	https://github.com/OS-Q/A21Y		/	master					Hard fork of https://github.com/GrumpyOldPizza/ArduinoCore-stm32l0.	",
		"",
	}

	assert.Equal(t, assertion, actual, "Written file has expected content.")
}
