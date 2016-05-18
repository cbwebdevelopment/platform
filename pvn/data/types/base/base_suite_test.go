package base_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestBase(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "data/types/base")
}
