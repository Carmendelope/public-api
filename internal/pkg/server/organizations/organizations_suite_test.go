/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package organizations

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"testing"
)

func TestOrganizationsPackage(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Organization package suite")
}
