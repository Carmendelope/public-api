/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package nodes

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"testing"
)

func TestNodesPackage(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Nodes package suite")
}
