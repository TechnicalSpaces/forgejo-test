// Copyright 2025 The Forgejo Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package integration

import (
	"net/http"
	"strings"
	"testing"

	"forgejo.org/modules/setting"
	"forgejo.org/tests"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

func TestForceFileOnlyCommitDiffs(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	// Save original setting value
	originalForceFileOnly := setting.Git.ForceFileOnlyCommitDiffs
	defer func() {
		setting.Git.ForceFileOnlyCommitDiffs = originalForceFileOnly
	}()

	session := loginUser(t, "user2")

	t.Run("ForceFileOnlyCommitDiffs disabled", func(t *testing.T) {
		setting.Git.ForceFileOnlyCommitDiffs = false
		
		req := NewRequest(t, "GET", "/user2/repo1/commits/branch/master")
		resp := session.MakeRequest(t, req, http.StatusOK)
		assert.Equal(t, http.StatusOK, resp.Code)

		doc := NewHTMLParser(t, resp.Body)
		commitLinks := doc.doc.Find("#commits-table tbody tr td.sha a[href*='/commit/']")
		assert.Positive(t, commitLinks.Length(), "Should have commit links")
		
		// Verify that commit links don't contain ?file-only=true
		commitLinks.Each(func(i int, s *goquery.Selection) {
			href, _ := s.Attr("href")
			assert.NotContains(t, href, "?file-only=true", "Commit link should not have file-only parameter when disabled")
		})
	})

	t.Run("ForceFileOnlyCommitDiffs enabled", func(t *testing.T) {
		setting.Git.ForceFileOnlyCommitDiffs = true
		
		req := NewRequest(t, "GET", "/user2/repo1/commits/branch/master")
		resp := session.MakeRequest(t, req, http.StatusOK)
		assert.Equal(t, http.StatusOK, resp.Code)

		doc := NewHTMLParser(t, resp.Body)
		commitLinks := doc.doc.Find("#commits-table tbody tr td.sha a[href*='/commit/']")
		assert.Positive(t, commitLinks.Length(), "Should have commit links")
		
		// Verify that commit links contain ?file-only=true
		foundFileOnlyLink := false
		commitLinks.Each(func(i int, s *goquery.Selection) {
			href, _ := s.Attr("href")
			if strings.Contains(href, "?file-only=true") {
				foundFileOnlyLink = true
			}
		})
		assert.True(t, foundFileOnlyLink, "At least one commit link should have file-only parameter when enabled")
	})
}