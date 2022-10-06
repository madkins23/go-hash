package crc16

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/sigurn/crc16"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/madkins23/go-hash/pkg/hash"
)

const fileLoremIpsum = "../../../testdata/lorem-ipsum.txt"

// testStandard is the standard '123456789' string.
// This appears to be a generally known 'check' value
// as discussed in https://pidlaboratory.com/9-rozne-algorytmy-crc/
// and displayed in https://crccalc.com/.
var testStandard = "123456789"

// testLoremIpsum is a larger sample
var testLoremIpsum = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam sodales aliquam finibus.
In bibendum purus id mattis rhoncus. Mauris nec libero laoreet, feugiat nisl at, ullamcorper justo.
Maecenas laoreet eleifend ultricies. Sed pulvinar convallis velit sit amet auctor.
Orci varius natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus.
Sed sed tortor eget est auctor tincidunt sit amet a nunc. Suspendisse laoreet turpis eget elit iaculis, congue iaculis enim egestas.`

// testCase associates a CRC16 method with an expected result.
type testCase struct {
	params     crc16.Params
	standard   uint16
	loremIpsum uint16
	fileResult uint16
}

// testCases configures the test parameters and results to be run.
// Results are taken from https://crccalc.com/ for the standard 'check' value
// as well as the value of testLoremIpsum and the supplied text file.
var testCases = []testCase{
	{crc16.CRC16_ARC, 0xBB3D, 0x6EFB, 0x7B09},
	{crc16.CRC16_AUG_CCITT, 0xE5CC, 0xAC8D, 0x1029},
	{crc16.CRC16_BUYPASS, 0xFEE8, 0x973E, 0x0116},
	{crc16.CRC16_CCITT_FALSE, 0x29B1, 0xEF31, 0xEDB4},
	{crc16.CRC16_CDMA2000, 0x4C06, 0x0304, 0x1580},
	{crc16.CRC16_DDS_110, 0x9ECF, 0x5B1A, 0xC64B},
	{crc16.CRC16_DECT_R, 0x007E, 0x1E3E, 0x87BF},
	{crc16.CRC16_DECT_X, 0x007F, 0x1E3F, 0x87BE},
	{crc16.CRC16_DNP, 0xEA82, 0xB5B9, 0x32AE},
	{crc16.CRC16_EN_13757, 0xC2B7, 0xA4E6, 0x9E86},
	{crc16.CRC16_GENIBUS, 0xD64E, 0x10CE, 0x124B},
	{crc16.CRC16_MAXIM, 0x44C2, 0x9104, 0x84F6},
	{crc16.CRC16_MCRF4XX, 0x6F91, 0x3596, 0x6A5C},
	{crc16.CRC16_RIELLO, 0x63D0, 0x9FA9, 0xF7AB},
	{crc16.CRC16_T10_DIF, 0xD0DB, 0x6E28, 0xDE75},
	{crc16.CRC16_TELEDISK, 0x0FB3, 0xCED9, 0xBB8D},
	{crc16.CRC16_TMS37157, 0x26B1, 0xA900, 0x97C7},
	{crc16.CRC16_USB, 0xB4C8, 0x2141, 0xA270},
	{crc16.CRC16_CRC_A, 0xBF05, 0x95D2, 0x00D1},
	{crc16.CRC16_KERMIT, 0x2189, 0x067D, 0x45AE},
	{crc16.CRC16_MODBUS, 0x4B37, 0xDEBE, 0x5D8F},
	{crc16.CRC16_X_25, 0x906E, 0xCA69, 0x95A3},
	{crc16.CRC16_XMODEM, 0x31C3, 0x38FD, 0xA240},
}

func testCaseLoop(t *testing.T, testData string,
	testWriter func(t *testing.T, h16 hash.Hash16, testData string, name string),
	testChecker func(t *testing.T, result uint16, tc testCase)) {
	var tc testCase
	for _, tc = range testCases {
		h16 := New(tc.params)
		require.NotNil(t, h16, "%s make hash", tc.params.Name)
		testWriter(t, h16, testData, tc.params.Name)
		crc := h16.Sum16()
		testChecker(t, crc, tc)
	}
}

func testWriterData(t *testing.T, h16 hash.Hash16, testData string, name string) {
	num, err := h16.Write([]byte(testData))
	assert.NoError(t, err, "%s hash write", name)
	assert.Equal(t, len(testData), num, "%s wrong data size from hash write", name)
}

func TestCheckerStandard(t *testing.T) {
	testCaseLoop(t, testStandard, testWriterData,
		func(t *testing.T, result uint16, tc testCase) {
			assert.Equal(t, tc.standard, result, "%s wrong CRC result", tc.params.Name)
		})
}

func TestCheckerLoremIpsum(t *testing.T) {
	testCaseLoop(t, testLoremIpsum, testWriterData,
		func(t *testing.T, result uint16, tc testCase) {
			assert.Equal(t, tc.loremIpsum, result, "%s wrong CRC result", tc.params.Name)
		})
}

func TestCheckerFileResult(t *testing.T) {
	testCaseLoop(t, fileLoremIpsum,
		func(t *testing.T, h16 hash.Hash16, fileName string, name string) {
			file, err := os.Open(fileName)
			require.NoError(t, err, "%s open file %s", name, fileName)
			_, err = io.Copy(h16, file)
			assert.NoError(t, err, "%s copy from file %s", name, fileName)
			_ = file.Close()
		}, func(t *testing.T, result uint16, tc testCase) {
			assert.Equal(t, tc.fileResult, result, "%s wrong CRC result", tc.params.Name)
		})
}

func ExampleHash16() {
	if file, err := os.Open(fileLoremIpsum); err == nil {
		h16 := New(crc16.CRC16_KERMIT)
		if _, err := io.Copy(h16, file); err == nil {
			fmt.Printf("%03d %s\n", h16.Sum16(), fileLoremIpsum)
		}
		_ = file.Close()
	} else {
		fmt.Printf("Unable to open file: %s\n", err.Error())
	}
	// Output:
	// 17838 ../../../testdata/lorem-ipsum.txt
}
