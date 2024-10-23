package contentful

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// UploadService service
type UploadService service

type Upload struct {
	Sys Sys `json:"sys"`
}

// Uploads creates a new upload and returns a reference ID
func (service *UploadService) Uploads(ctx context.Context, spaceID string, file io.Reader) (*Upload, error) {
	path := fmt.Sprintf("/spaces/%s%s/uploads", spaceID, getEnvPath(service.c))
	method := http.MethodPost

	req, err := service.c.newRequest(ctx, method, path, nil, file, map[string]string{"Content-Type": "application/octet-stream"})
	if err != nil {
		return nil, err
	}
	var uploadResponse Upload
	if err := service.c.do(req, &uploadResponse); err != nil {
		return nil, err
	}
	return &uploadResponse, nil
}
