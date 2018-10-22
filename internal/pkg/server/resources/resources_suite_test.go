/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package resources

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"testing"
)

func TestResourcesPackage(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Resources package suite")
}
