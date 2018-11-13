/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package applications

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"testing"
)

func TestApplicationsPackage(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Applications package suite")
}
