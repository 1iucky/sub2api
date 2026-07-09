package dto

import (
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"
)

func TestAnnouncementDTOIncludesCategory(t *testing.T) {
	createdAt := time.Unix(1776790020, 0)
	ann := AnnouncementFromService(&service.Announcement{
		ID:         1,
		Title:      "Model update",
		Content:    "New models",
		Status:     service.AnnouncementStatusActive,
		NotifyMode: service.AnnouncementNotifyModeSilent,
		Category:   service.AnnouncementCategoryModelUpdate,
		CreatedAt:  createdAt,
		UpdatedAt:  createdAt,
	})

	require.NotNil(t, ann)
	require.Equal(t, service.AnnouncementCategoryModelUpdate, ann.Category)
}

func TestUserAnnouncementDTOIncludesCategory(t *testing.T) {
	createdAt := time.Unix(1776790020, 0)
	ann := UserAnnouncementFromService(&service.UserAnnouncement{
		Announcement: service.Announcement{
			ID:         1,
			Title:      "Release notes",
			Content:    "Changed",
			Status:     service.AnnouncementStatusActive,
			NotifyMode: service.AnnouncementNotifyModeSilent,
			Category:   service.AnnouncementCategoryChangelog,
			CreatedAt:  createdAt,
			UpdatedAt:  createdAt,
		},
	})

	require.NotNil(t, ann)
	require.Equal(t, service.AnnouncementCategoryChangelog, ann.Category)
}
