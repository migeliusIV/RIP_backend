package repository

import (
	"context"
	"fmt"
	"front_start/internal/app/ds"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *Repository) GetGates() ([]ds.Gate, error) {
	var gates []ds.Gate

	err := r.db.Find(&gates).Error
	if err != nil {
		return nil, err
	}

	if len(gates) == 0 {
		return nil, fmt.Errorf("gates not found")
	}
	return gates, nil
}

func (r *Repository) GetGatesByName(title string) ([]ds.Gate, error) {
	var gates []ds.Gate
	err := r.db.Where("title ILIKE ?", "%"+title+"%").Find(&gates).Error
	if err != nil {
		return nil, err
	}
	return gates, nil
}

func (r *Repository) GetGateByID(id int) (*ds.Gate, error) {
	var gate ds.Gate
	err := r.db.First(&gate, id).Error
	if err != nil {
		return nil, err
	}
	return &gate, nil
}

// ---- JSON API helpers for services ----

func (r *Repository) ListGates(title string) ([]ds.Gate, error) {
	var gates []ds.Gate
	q := r.db
	if title != "" {
		q = q.Where("title ILIKE ?", "%"+title+"%")
	}
	if err := q.Find(&gates).Error; err != nil {
		return nil, err
	}
	return gates, nil
}

func (r *Repository) AddGate(g *ds.Gate) error {
	return r.db.Create(g).Error
}

func (r *Repository) UpdateGate(id uint, title, description, fullInfo, theAxis string, status *bool) (*ds.Gate, error) {
	var gate ds.Gate
	if err := r.db.First(&gate, id).Error; err != nil {
		return nil, err
	}
	if title != "" {
		gate.Title = title
	}
	if description != "" {
		gate.Description = description
	}
	if fullInfo != "" {
		gate.FullInfo = fullInfo
	}
	if theAxis != "" {
		gate.TheAxis = theAxis
	}
	if status != nil {
		gate.Status = *status
	}
	if err := r.db.Save(&gate).Error; err != nil {
		return nil, err
	}
	return &gate, nil
}

func (r *Repository) DeleteGate(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// First, get the gate to check if it has an image
		var gate ds.Gate
		if err := tx.First(&gate, id).Error; err != nil {
			return fmt.Errorf("gate with id %d not found: %w", id, err)
		}

		// Remove image from MinIO if exists
		if gate.Image != nil && *gate.Image != "" {
			oldImageURL, err := url.Parse(*gate.Image)
			if err == nil {
				oldObjectName := strings.TrimPrefix(oldImageURL.Path, fmt.Sprintf("/%s/", r.minio.bucketName))
				r.minio.client.RemoveObject(context.Background(), r.minio.bucketName, oldObjectName, minio.RemoveObjectOptions{})
			}
		}

		// Delete the gate record from database
		return tx.Delete(&ds.Gate{}, id).Error
	})
}

func (r *Repository) SetServiceImage(id uint, image string) error {
	return r.db.Model(&ds.Gate{}).Where("id_gate = ?", id).Update("image", image).Error
}

// SaveServiceImage uploads file to MinIO and updates gate image path with transaction and old image cleanup
func (r *Repository) SaveServiceImage(ctx context.Context, id uint, fileHeader *multipart.FileHeader) (string, error) {
	var finalImageURL string
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var gate ds.Gate
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&gate, id).Error; err != nil {
			return fmt.Errorf("gate with id %d not found: %w", id, err)
		}

		const imagePathPrefix = "img/"

		// Remove old image if exists
		if gate.Image != nil && *gate.Image != "" {
			oldImageURL, err := url.Parse(*gate.Image)
			if err == nil {
				oldObjectName := strings.TrimPrefix(oldImageURL.Path, fmt.Sprintf("/%s/", r.minio.bucketName))
				r.minio.client.RemoveObject(context.Background(), r.minio.bucketName, oldObjectName, minio.RemoveObjectOptions{})
			}
		}

		// Generate safe filename
		fileName := filepath.Base(fileHeader.Filename)
		ext := filepath.Ext(fileName)
		name := fileName[:len(fileName)-len(ext)]
		if name == "" {
			name = "image"
		}
		// Simple sanitization - replace non-alphanumeric with dash
		safeName := strings.ReplaceAll(name, " ", "-")
		fileName = safeName + ext
		objectName := imagePathPrefix + fileName

		// Upload to MinIO
		file, err := fileHeader.Open()
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = r.minio.client.PutObject(context.Background(), r.minio.bucketName, objectName, file, fileHeader.Size, minio.PutObjectOptions{
			ContentType: fileHeader.Header.Get("Content-Type"),
		})

		if err != nil {
			return fmt.Errorf("failed to upload to minio: %w", err)
		}

		// Generate full URL
		minioEndpoint := os.Getenv("MINIO_ENDPOINT")
		if minioEndpoint == "" {
			minioEndpoint = "127.0.0.1:9000"
		}
		imageURL := fmt.Sprintf("http://%s/%s/%s", minioEndpoint, r.minio.bucketName, objectName)

		// Update database
		if err := tx.Model(&gate).Update("image", imageURL).Error; err != nil {
			return fmt.Errorf("failed to update gate image url in db: %w", err)
		}

		finalImageURL = imageURL
		return nil
	})
	if err != nil {
		return "", err
	}
	return finalImageURL, nil
}
