/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package clusters

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"testing"
)

func TestClustersPackage(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Clusters package suite")
}
