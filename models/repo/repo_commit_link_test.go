// Copyright 2025 The Forgejo Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package repo

import (
	"testing"

	"forgejo.org/modules/setting"

	"github.com/stretchr/testify/assert"
)

func TestCommitLink(t *testing.T) {
	// Save original setting value
	originalForceFileOnly := setting.Git.ForceFileOnlyCommitDiffs
	defer func() {
		setting.Git.ForceFileOnlyCommitDiffs = originalForceFileOnly
	}()

	repo := &Repository{
		OwnerName: "testowner",
		Name:      "testrepo",
	}

	commitID := "abc123def456"

	t.Run("ForceFileOnlyCommitDiffs disabled", func(t *testing.T) {
		setting.Git.ForceFileOnlyCommitDiffs = false
		
		result := repo.CommitLink(commitID)
		expected := "/testowner/testrepo/commit/abc123def456"
		assert.Equal(t, expected, result)
	})

	t.Run("ForceFileOnlyCommitDiffs enabled", func(t *testing.T) {
		setting.Git.ForceFileOnlyCommitDiffs = true
		
		result := repo.CommitLink(commitID)
		expected := "/testowner/testrepo/commit/abc123def456?file-only=true"
		assert.Equal(t, expected, result)
	})

	t.Run("Empty commit ID", func(t *testing.T) {
		setting.Git.ForceFileOnlyCommitDiffs = true
		
		result := repo.CommitLink("")
		assert.Equal(t, "", result)
	})
}