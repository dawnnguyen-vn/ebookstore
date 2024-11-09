package models

import (
	"fmt"

	"google.golang.org/api/drive/v3"
)

const DRIVE_FOLDER_TYPE = "application/vnd.google-apps.folder"
const NAVIGATION_TYPE = "application/atom+xml;profile=opds-catalog;kind=acquisition"

type Feed struct {
	Name     string
	Link     string
	MimeType string
	Type     string
}

func NewFeed(driveId string, driveName string, driveMimeType string) Feed {
	MimeType := driveMimeType
	Type := "acquisition"
	Link := fmt.Sprintf("https://drive.google.com/uc?export=download&id=%s", driveId)

	if MimeType == DRIVE_FOLDER_TYPE {
		MimeType = NAVIGATION_TYPE
		Type = "navigation"
		Link = fmt.Sprintf("/opds/catalogs/%s", driveId)
	}

	return Feed{
		Name:     driveName,
		Link:     Link,
		MimeType: MimeType,
		Type:     Type,
	}
}

func ConvertFilesToFeeds(files []*drive.File) []Feed {
	feeds := []Feed{}
	for _, file := range files {
		feed := NewFeed(file.Id, file.Name, file.MimeType)
		feeds = append(feeds, feed)
	}
	return feeds
}
