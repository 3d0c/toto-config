package models

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/3d0c/toto-config/pkg/config"
	"github.com/3d0c/toto-config/pkg/helpers"
	"github.com/3d0c/toto-config/pkg/log"
)

var (
	// Created samples list is used for benchmarks
	testSamples []*ConfigScheme
	// Sample which used as "expected" for FindBy()
	expectedSample *ConfigScheme
)

func TestFindBy(t *testing.T) {
	var (
		cm     *Config
		sample *ConfigScheme
		err    error
	)

	cm, err = NewConfigModel()
	assert.Nil(t, err)

	fmt.Printf("Looking for: %s, %s, %d\n", expectedSample.Package, expectedSample.CountryCode, expectedSample.PercentileMax)
	sample, err = cm.FindBy(expectedSample.Package, expectedSample.CountryCode, expectedSample.PercentileMax)
	assert.Nil(t, err)
	assert.Equal(t, expectedSample.MainSKU, sample.MainSKU)
}

func BenchmarkFindBy(b *testing.B) {
	var (
		cm  *Config
		n   = len(testSamples)
		j   = 0
		err error
	)

	for i := 0; i < b.N; i++ {
		if j >= n {
			j = 0
		}
		if _, err = cm.FindBy(testSamples[j].Package, testSamples[j].CountryCode, testSamples[j].PercentileMax); err != nil {
			panic(err)
		}
	}
}

// - generate random number of samples. max number is 10,000
// - insert them into toto_configuration table
// - return random sample to test FindBy method
func prepareDatabase() error {
	var (
		n         = helpers.RandomSeed(1000, 10000)
		rndSample = helpers.RandomSeed(1, n)
		cm        *Config
		err       error
	)

	if cm, err = NewConfigModel(); err != nil {
		return fmt.Errorf("error initializing Config model - %s", err)
	}

	testSamples = make([]*ConfigScheme, 0, n)

	for i := 0; i < n; i++ {
		sample := &ConfigScheme{
			Package:       helpers.RandomString(32),
			CountryCode:   helpers.RandomString(1),
			PercentileMin: helpers.RandomSeed(1, 30),
			PercentileMax: helpers.RandomSeed(31, 100),
			MainSKU:       helpers.RandomString(32),
		}

		if err = cm.Create(sample); err != nil {
			return fmt.Errorf("error creating sample - %s", err)
		}

		testSamples = append(testSamples, sample)

		if i == rndSample {
			expectedSample = sample
		}
	}

	return nil
}

func TestMain(m *testing.M) {
	const (
		// root directory of the project
		dbFileName = "/tmp/models_test.db"
	)
	var (
		err error
	)

	log.InitLogger(config.Logger{
		Level:     "debug",
		AddCaller: true,
	})

	config.TheConfig().Database.DSN = dbFileName
	config.TheConfig().Database.Dialect = "sqlite3"

	if err = prepareDatabase(); err != nil {
		fmt.Printf("Error preparing testing environment - %s\n", err)
		os.Exit(-1)
	}

	exitval := m.Run()

	// if err := os.Remove(dbFileName); err != nil {
	// 	fmt.Printf("ERROR: removing testing database, error - %s\n", err)
	// }

	os.Exit(exitval)
}
