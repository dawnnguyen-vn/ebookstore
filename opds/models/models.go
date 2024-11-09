package models

import (
	"os"

	"google.golang.org/api/drive/v3"
)

var (
	DRIVE_FOLDER_TYPE   = "application/vnd.google-apps.folder"
	NAVIGATION_TYPE     = "application/atom+xml;profile=opds-catalog;kind=acquisition"
	BASE_PATH           = os.Getenv("API_GATEWAY_STAGE") + "/opds/catalogs"
	DRIVE_DOWNLOAD_LINK = "https://drive.google.com/uc?export=download&id="
)

type Feed struct {
	Name     string
	Link     string
	MimeType string
	Type     string
}

func NewFeed(driveId string, driveName string, driveMimeType string) *Feed {
	MimeType := driveMimeType
	Type := "acquisition"
	Link := DRIVE_DOWNLOAD_LINK + driveId

	if MimeType == DRIVE_FOLDER_TYPE {
		MimeType = NAVIGATION_TYPE
		Type = "navigation"
		Link = BASE_PATH + "/" + driveId
	}

	return &Feed{
		Name:     driveName,
		Link:     Link,
		MimeType: MimeType,
		Type:     Type,
	}
}

func ConvertFilesToFeeds(files []*drive.File) []*Feed {
	feeds := []*Feed{}
	for _, file := range files {
		feed := NewFeed(file.Id, file.Name, file.MimeType)
		feeds = append(feeds, feed)
	}
	return feeds
}
