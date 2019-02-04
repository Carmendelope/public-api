/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package devices

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"testing"
)

func TestDevicesPackage(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Devices package suite")
}
