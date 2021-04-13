package dropbox

import (
	"bytes"
	"context"
	"io"
	"os"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/sharing"
)

type DropboxSdkFiles interface {
	Upload(arg *files.CommitInfo, content io.Reader) (res *files.FileMetadata, err error)
	DeleteV2(arg *files.DeleteArg) (res *files.DeleteResult, err error)
}

type DropboxSdkSharing interface {
	CreateSharedLinkWithSettings(arg *sharing.CreateSharedLinkWithSettingsArg) (res sharing.IsSharedLinkMetadata, err error)
}

type Dropbox interface {
	Uploader(filename string, buffer []byte) (path string, err error)
	CreateSharedLink(path string) (string, error)
	Delete(path string) error
}

type dbx struct {
	ctx     context.Context
	token   string
	files   DropboxSdkFiles
	sharing DropboxSdkSharing
}

func NewDropboxRepository(ctx context.Context) Dropbox {
	token := os.Getenv("DROPBOX_ACCESS_TOKEN")
	dbxConfig := dropbox.Config{
		Token: token,
	}
	dbxFiles := files.New(dbxConfig)
	dbxSharing := sharing.New(dbxConfig)

	return &dbx{ctx, token, dbxFiles, dbxSharing}
}

func (repo dbx) Uploader(filename string, buffer []byte) (path string, err error) {
	payload := new(bytes.Buffer)
	payload.Write(buffer)

	args := &files.CommitInfo{
		Path: "/getprint/" + filename,
		Mode: &files.WriteMode{
			Tagged: dropbox.Tagged{Tag: "add"},
		},
		Autorename:     true,
		Mute:           false,
		StrictConflict: false,
	}
	res, err := repo.files.Upload(args, payload)
	if err != nil {
		return "", err
	}

	return res.PathLower, nil
}

func (repo dbx) CreateSharedLink(path string) (string, error) {
	args := &sharing.CreateSharedLinkWithSettingsArg{
		Path:     path,
		Settings: nil,
	}

	res, err := repo.sharing.CreateSharedLinkWithSettings(args)
	if err != nil {
		return "", err
	}
	linkRes := res.(*sharing.FileLinkMetadata)

	return linkRes.Url, nil
}

func (repo dbx) Delete(path string) error {
	args := &files.DeleteArg{
		Path: path,
	}
	_, err := repo.files.DeleteV2(args)
	if err != nil {
		return err
	}

	return nil
}
