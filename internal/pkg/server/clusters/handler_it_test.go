/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package clusters

import (
	"github.com/nalej/public-api/internal/pkg/utils"
	"github.com/onsi/ginkgo"
	"github.com/rs/zerolog/log"
	"os"
)

var _ = ginkgo.Describe("Clusters", func() {

	if ! utils.RunIntegrationTests() {
		log.Warn().Msg("Integration tests are skipped")
		return
	}

	var (
		inframgrModelAddress= os.Getenv("IT_INFRAMGR_ADDRESS")
	)

	if inframgrModelAddress == "" {
		ginkgo.Fail("missing environment variables")
	}

	

})