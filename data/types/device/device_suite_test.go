package device_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDeviceevent(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "data/types/device")
}
