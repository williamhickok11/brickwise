package service

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"brickwise/backend/internal/errors"
	"brickwise/backend/internal/models"
	"brickwise/backend/internal/repository"
)

const (
	defaultPropertyImagesDir = "./data/property-images"
	maxPropertyImageBytes    = 8 << 20 // 8 MiB
)

var allowedImageTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/jpg":  ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
	"image/gif":  ".gif",
}

// PropertyImageService stores multiple photos per property on local disk.
type PropertyImageService struct {
	repo         repository.PropertyImageRepository
	propertyRepo repository.PropertyRepository
	dir          string
}

func NewPropertyImageService(repo repository.PropertyImageRepository, propertyRepo repository.PropertyRepository) *PropertyImageService {
	dir := os.Getenv("PROPERTY_IMAGES_DIR")
	if strings.TrimSpace(dir) == "" {
		dir = defaultPropertyImagesDir
	}
	return &PropertyImageService{
		repo:         repo,
		propertyRepo: propertyRepo,
		dir:          dir,
	}
}

func (s *PropertyImageService) List(propertyID int) ([]*models.PropertyImage, error) {
	if _, err := s.propertyRepo.GetByID(propertyID); err != nil {
		return nil, errors.ErrNotFound
	}
	images, err := s.repo.ListByPropertyID(propertyID)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrInternal)
	}
	if images == nil {
		images = []*models.PropertyImage{}
	}
	return images, nil
}

// Upload saves an image for the property. The first image becomes main automatically.
func (s *PropertyImageService) Upload(propertyID int, fileName string, contentType string, size int64, body io.Reader) (*models.PropertyImage, error) {
	if _, err := s.propertyRepo.GetByID(propertyID); err != nil {
		return nil, errors.ErrNotFound
	}

	if size > maxPropertyImageBytes {
		return nil, errors.New("FILE_TOO_LARGE", "image must be 8MB or smaller", errors.ErrInvalidInput.Status)
	}

	ext, err := normalizeImageType(fileName, contentType)
	if err != nil {
		return nil, err
	}

	count, err := s.repo.CountByPropertyID(propertyID)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrInternal)
	}

	if err := os.MkdirAll(s.propertyDir(propertyID), 0o755); err != nil {
		return nil, errors.Wrap(err, errors.ErrInternal)
	}

	// Insert metadata first so we get a stable id for the file name.
	img := &models.PropertyImage{
		PropertyID:  propertyID,
		FileName:    safeFileName(fileName, ext),
		ContentType: canonicalContentType(ext),
		Size:        0,
		IsMain:      count == 0,
	}
	if err := s.repo.Create(img); err != nil {
		return nil, errors.Wrap(err, errors.ErrInternal)
	}

	storagePath := s.storagePath(propertyID, img.ID)
	tmpPath := storagePath + ".tmp"
	out, err := os.Create(tmpPath)
	if err != nil {
		_ = s.repo.Delete(img.ID)
		return nil, errors.Wrap(err, errors.ErrInternal)
	}

	written, copyErr := io.Copy(out, io.LimitReader(body, maxPropertyImageBytes+1))
	closeErr := out.Close()
	if copyErr != nil {
		os.Remove(tmpPath)
		_ = s.repo.Delete(img.ID)
		return nil, errors.Wrap(copyErr, errors.ErrInternal)
	}
	if closeErr != nil {
		os.Remove(tmpPath)
		_ = s.repo.Delete(img.ID)
		return nil, errors.Wrap(closeErr, errors.ErrInternal)
	}
	if written == 0 {
		os.Remove(tmpPath)
		_ = s.repo.Delete(img.ID)
		return nil, errors.New("INVALID_FILE", "file is empty", errors.ErrInvalidInput.Status)
	}
	if written > maxPropertyImageBytes {
		os.Remove(tmpPath)
		_ = s.repo.Delete(img.ID)
		return nil, errors.New("FILE_TOO_LARGE", "image must be 8MB or smaller", errors.ErrInvalidInput.Status)
	}

	if err := os.Rename(tmpPath, storagePath); err != nil {
		os.Remove(tmpPath)
		_ = s.repo.Delete(img.ID)
		return nil, errors.Wrap(err, errors.ErrInternal)
	}

	// Size is stored at insert time as 0; update via delete+recreate is heavy —
	// re-read is fine for API response; persist size with a lightweight update:
	img.Size = written
	if err := s.repo.UpdateSize(img.ID, written); err != nil {
		os.Remove(storagePath)
		_ = s.repo.Delete(img.ID)
		return nil, errors.Wrap(err, errors.ErrInternal)
	}

	return img, nil
}

func (s *PropertyImageService) SetMain(propertyID, imageID int) (*models.PropertyImage, error) {
	if _, err := s.propertyRepo.GetByID(propertyID); err != nil {
		return nil, errors.ErrNotFound
	}
	img, err := s.repo.GetByID(imageID)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	if img.PropertyID != propertyID {
		return nil, errors.ErrNotFound
	}
	if err := s.repo.SetMain(propertyID, imageID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(err, errors.ErrInternal)
	}
	img.IsMain = true
	return img, nil
}

func (s *PropertyImageService) Delete(propertyID, imageID int) error {
	if _, err := s.propertyRepo.GetByID(propertyID); err != nil {
		return errors.ErrNotFound
	}
	img, err := s.repo.GetByID(imageID)
	if err != nil {
		return errors.ErrNotFound
	}
	if img.PropertyID != propertyID {
		return errors.ErrNotFound
	}

	wasMain := img.IsMain
	if err := s.repo.Delete(imageID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return errors.ErrNotFound
		}
		return errors.Wrap(err, errors.ErrInternal)
	}
	os.Remove(s.storagePath(propertyID, imageID))

	if wasMain {
		remaining, listErr := s.repo.ListByPropertyID(propertyID)
		if listErr == nil && len(remaining) > 0 {
			_ = s.repo.SetMain(propertyID, remaining[0].ID)
		}
	}
	return nil
}

// CleanupForProperty removes on-disk files for a property (DB rows cascade separately).
func (s *PropertyImageService) CleanupForProperty(propertyID int) {
	os.RemoveAll(s.propertyDir(propertyID))
}

// OpenForDownload returns metadata and an open file handle. Caller must close the file.
func (s *PropertyImageService) OpenForDownload(propertyID, imageID int) (*models.PropertyImage, *os.File, error) {
	if _, err := s.propertyRepo.GetByID(propertyID); err != nil {
		return nil, nil, errors.ErrNotFound
	}
	img, err := s.repo.GetByID(imageID)
	if err != nil {
		return nil, nil, errors.ErrNotFound
	}
	if img.PropertyID != propertyID {
		return nil, nil, errors.ErrNotFound
	}
	f, err := os.Open(s.storagePath(propertyID, imageID))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, errors.ErrNotFound
		}
		return nil, nil, errors.Wrap(err, errors.ErrInternal)
	}
	return img, f, nil
}

func (s *PropertyImageService) propertyDir(propertyID int) string {
	return filepath.Join(s.dir, fmt.Sprintf("%d", propertyID))
}

func (s *PropertyImageService) storagePath(propertyID, imageID int) string {
	return filepath.Join(s.propertyDir(propertyID), fmt.Sprintf("%d", imageID))
}

func normalizeImageType(fileName, contentType string) (string, error) {
	ct := strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
	if ext, ok := allowedImageTypes[ct]; ok {
		return ext, nil
	}

	lowerName := strings.ToLower(fileName)
	switch {
	case strings.HasSuffix(lowerName, ".jpg"), strings.HasSuffix(lowerName, ".jpeg"):
		return ".jpg", nil
	case strings.HasSuffix(lowerName, ".png"):
		return ".png", nil
	case strings.HasSuffix(lowerName, ".webp"):
		return ".webp", nil
	case strings.HasSuffix(lowerName, ".gif"):
		return ".gif", nil
	default:
		return "", errors.New("INVALID_FILE_TYPE", "only JPEG, PNG, WebP, or GIF images are allowed", errors.ErrInvalidInput.Status)
	}
}

func canonicalContentType(ext string) string {
	switch ext {
	case ".png":
		return "image/png"
	case ".webp":
		return "image/webp"
	case ".gif":
		return "image/gif"
	default:
		return "image/jpeg"
	}
}

func safeFileName(fileName, ext string) string {
	base := filepath.Base(fileName)
	if base == "." || base == "/" || base == "" {
		return "image" + ext
	}
	if !strings.Contains(strings.ToLower(base), ".") {
		return base + ext
	}
	return base
}
