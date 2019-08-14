package record_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRecord(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Record Suite")
}
