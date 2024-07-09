package types

import "io"

type UploadAssetData struct {
	UserID       int64
	AssetName    string
	OriginalName string
	ContentType  string
	Body         io.Reader
}
