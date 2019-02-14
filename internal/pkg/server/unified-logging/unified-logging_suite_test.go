/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package unified_logging

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"testing"
)

func TestUnifiedLoggingPackage(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Unified Logging package suite")
}
