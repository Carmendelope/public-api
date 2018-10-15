/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package users

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"testing"
)

func TestUsersPackage(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Users package suite")
}
