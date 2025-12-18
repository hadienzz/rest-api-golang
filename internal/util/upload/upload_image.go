package upload

import (
	"bytes"
	"context"
	"fmt"
	"go-fiber-api/internal/config"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/google/uuid"
)

type SupabaseUploadResult struct {
	PublicURL  string
	ObjectPath string
}

func UploadToSupabaseStorage(
	ctx context.Context,
	fileHeader *multipart.FileHeader,
	prefix string,
) (*SupabaseUploadResult, error) {

	configuration := config.Get()

	uploadedFile, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("open uploaded file failed: %w", err)
	}
	defer uploadedFile.Close()

	fileBytes, err := io.ReadAll(uploadedFile)
	if err != nil {
		return nil, fmt.Errorf("read uploaded file failed: %w", err)
	}

	fileExtension := "bin"
	if ext := path.Ext(fileHeader.Filename); ext != "" {
		fileExtension = ext[1:]
	}

	objectPath := fmt.Sprintf(
		"%s/%s.%s",
		prefix,
		uuid.NewString(),
		fileExtension,
	)

	// Build upload endpoint
	uploadURL := fmt.Sprintf(
		"%s/storage/v1/object/%s/%s",
		configuration.SupabaseURL,
		url.PathEscape("merchant"),
		url.PathEscape(objectPath),
	)

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		uploadURL,
		bytes.NewReader(fileBytes),
	)
	if err != nil {
		return nil, err
	}

	request.Header.Set(
		"Authorization",
		"Bearer "+configuration.SupabaseServiceKey,
	)
	request.Header.Set("Content-Type", fileHeader.Header.Get("Content-Type"))
	request.Header.Set("x-upsert", "true")

	httpClient := &http.Client{Timeout: 30 * time.Second}
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		responseBody, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf(
			"supabase storage upload failed [%d]: %s",
			response.StatusCode,
			string(responseBody),
		)
	}

	publicURL := fmt.Sprintf(
		"%s/storage/v1/object/public/%s/%s",
		configuration.SupabaseURL,
		"merchant",
		objectPath,
	)
	fmt.Println(objectPath)
	return &SupabaseUploadResult{
		PublicURL:  publicURL,
		ObjectPath: objectPath,
	}, nil
}
