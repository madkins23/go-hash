package crc8

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/sigurn/crc8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go-hash/pkg/hash"
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

// testCase associates a CRC8 method with an expected result.
type testCase struct {
	params     crc8.Params
	standard   uint8
	loremIpsum uint8
	fileResult uint8
}

// testCases configures the test parameters and results to be run.
// Results are taken from https://crccalc.com/ for the standard 'check' value
// as well as the value of testLoremIpsum and the supplied text file.
var testCases = []testCase{
	{crc8.CRC8, 0xF4, 0x3C, 0xC0},
	{crc8.CRC8_CDMA2000, 0xDA, 0xEB, 0xF5},
	{crc8.CRC8_DARC, 0x15, 0xC7, 0x5B},
	{crc8.CRC8_DVB_S2, 0xBC, 0x11, 0x95},
	{crc8.CRC8_EBU, 0x97, 0x4E, 0x0C},
	{crc8.CRC8_I_CODE, 0x7E, 0x3D, 0x82},
	{crc8.CRC8_ITU, 0xA1, 0x69, 0x95},
	{crc8.CRC8_MAXIM, 0xA1, 0xA5, 0x14},
	{crc8.CRC8_ROHC, 0xD0, 0xD2, 0x11},
	{crc8.CRC8_WCDMA, 0x25, 0x3A, 0x36},
}

func testCaseLoop(t *testing.T, testData string,
	testWriter func(t *testing.T, h8 hash.Hash8, testData string, name string),
	testChecker func(t *testing.T, result uint8, tc testCase)) {
	var tc testCase
	for _, tc = range testCases {
		h8 := New(tc.params)
		require.NotNil(t, h8, "%s make hash", tc.params.Name)
		testWriter(t, h8, testData, tc.params.Name)
		crc := h8.Sum8()
		testChecker(t, crc, tc)
	}
}

func testWriterData(t *testing.T, h8 hash.Hash8, testData string, name string) {
	num, err := h8.Write([]byte(testData))
	assert.NoError(t, err, "%s hash write", name)
	assert.Equal(t, len(testData), num, "%s wrong data size from hash write", name)
}

func TestCheckerStandard(t *testing.T) {
	testCaseLoop(t, testStandard, testWriterData,
		func(t *testing.T, result uint8, tc testCase) {
			assert.Equal(t, tc.standard, result, "%s wrong CRC result", tc.params.Name)
		})
}

func TestCheckerLoremIpsum(t *testing.T) {
	testCaseLoop(t, testLoremIpsum, testWriterData,
		func(t *testing.T, result uint8, tc testCase) {
			assert.Equal(t, tc.loremIpsum, result, "%s wrong CRC result", tc.params.Name)
		})
}

func TestCheckerFileResult(t *testing.T) {
	testCaseLoop(t, fileLoremIpsum,
		func(t *testing.T, h8 hash.Hash8, fileName string, name string) {
			file, err := os.Open(fileName)
			require.NoError(t, err, "%s open file %s", name, fileName)
			_, err = io.Copy(h8, file)
			assert.NoError(t, err, "%s copy from file %s", name, fileName)
			_ = file.Close()
		}, func(t *testing.T, result uint8, tc testCase) {
			assert.Equal(t, tc.fileResult, result, "%s wrong CRC result", tc.params.Name)
		})
}

func ExampleHash8() {
	if file, err := os.Open(fileLoremIpsum); err == nil {
		h8 := New(crc8.CRC8)
		if _, err := io.Copy(h8, file); err == nil {
			fmt.Printf("%03d %s\n", h8.Sum8(), fileLoremIpsum)
		}
		_ = file.Close()
	} else {
		fmt.Printf("Unable to open file: %s\n", err.Error())
	}
	// Output:
	// 192 ../../../testdata/lorem-ipsum.txt
}
