/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package roles

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"testing"
)

func TestRolesPackage(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Roles package suite")
}
