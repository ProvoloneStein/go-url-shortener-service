package repositories

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"

	"go.uber.org/zap"

	"github.com/ProvoloneStein/go-url-shortener-service/configs"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/models"
)

type ShorterRecord struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	UserID      string `json:"user_id"`
	Deleted     string `json:"deleted"`
}

type FileRepository struct {
	cfg    configs.AppConfig
	logger *zap.Logger
	file   *os.File
	writer *bufio.Writer
	reader *bufio.Reader
	store  map[string][3]string
	uuid   int
}

func NewFileRepository(cfg configs.AppConfig, logger *zap.Logger, file *os.File) (*FileRepository, error) {
	repo := FileRepository{
		cfg:    cfg,
		logger: logger,
		file:   file,
		store:  make(map[string][3]string),
		writer: bufio.NewWriter(file),
		reader: bufio.NewReader(file),
		uuid:   0,
	}

	for {
		record, err := repo.readString()
		if err != nil {
			return nil, fmt.Errorf("repository: %w", err)
		}
		if record == nil {
			break
		}
		repo.store[record.ShortURL] = [3]string{record.OriginalURL, record.UserID, record.Deleted}
		repo.uuid, err = strconv.Atoi(record.UUID)
		if err != nil {
			return nil, fmt.Errorf("repository: ошибка при получения uuid записи: %s", err)
		}
	}
	repo.uuid += 1
	return &repo, nil
}

func (r *FileRepository) readString() (*ShorterRecord, error) {
	// Читаем данные до символа переноса строки
	data, err := r.reader.ReadBytes('\n')
	if err == io.EOF {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("ошибка при чтении строки %w", err)
	}
	// Преобразуем данные из JSON-представления в структуру
	record := ShorterRecord{}
	err = json.Unmarshal(data, &record)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *FileRepository) writeString(record ShorterRecord) error {
	data, err := json.Marshal(&record)
	if err != nil {
		return err
	}

	// Записываем событие в буфер
	if _, err := r.writer.Write(data); err != nil {
		return fmt.Errorf("ошибка при записи строки %w", err)
	}

	// Добавляем перенос строки
	if err := r.writer.WriteByte('\n'); err != nil {
		return fmt.Errorf("ошибка при записи строки %w", err)
	}

	r.uuid += 1

	// Записываем буфер в файл
	return r.writer.Flush()
}

func (r *FileRepository) validateUniqueShortURL(ctx context.Context, shortURL string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	if _, ok := r.store[shortURL]; !ok {
		return nil
	}
	return ErrShortURLExists
}

func (r *FileRepository) Create(ctx context.Context, userID, fullURL, shortURL string) (string, error) {
	select {
	case <-ctx.Done():
		return "", fmt.Errorf("repository: %w", ctx.Err())
	default:
	}
	if err := r.validateUniqueShortURL(ctx, shortURL); err != nil {
		return "", fmt.Errorf("repository: %w", err)
	}
	for key, val := range r.store {
		if val[0] == fullURL {
			return key, fmt.Errorf("repository: %w", ErrorUniqueViolation)
		}
	}
	r.store[shortURL] = [3]string{fullURL, userID, "f"}
	if err := r.writeString(ShorterRecord{strconv.Itoa(r.uuid), shortURL, fullURL, userID, "f"}); err != nil {
		delete(r.store, shortURL)
		return shortURL, fmt.Errorf("repository: %w", err)
	}
	return shortURL, nil
}

func (r *FileRepository) BatchCreate(ctx context.Context, data []models.BatchCreateData) ([]models.BatchCreateResponse, error) {
	var response []models.BatchCreateResponse

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	for _, val := range data {
		if err := r.validateUniqueShortURL(ctx, val.ShortURL); err != nil {
			return []models.BatchCreateResponse{models.BatchCreateResponse{ShortURL: val.ShortURL, UUID: val.UUID}}, fmt.Errorf("repository: %w", err)
		}
	}
	for _, val := range data {
		shortURL, err := url.JoinPath(r.cfg.BaseURL, val.ShortURL)
		if err != nil {
			r.logger.Error("ошибка при формировании url", zap.Error(err))
			return response, err
		}
		_, err = r.Create(ctx, val.UserID, val.URL, val.ShortURL)
		if err != nil && !errors.Is(err, ErrorUniqueViolation) {
			r.logger.Error("ошибка при записи url", zap.Error(err))
			return response, fmt.Errorf("repository: %w", err)
		}
		response = append(response, models.BatchCreateResponse{ShortURL: shortURL, UUID: val.UUID})
	}
	return response, nil
}

func (r *FileRepository) GetByShort(ctx context.Context, shortURL string) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}
	data, ok := r.store[shortURL]
	if ok {
		if data[2] == "t" {
			return "", fmt.Errorf("repository: %w", ErrDeleted)
		}
		return data[0], nil
	}
	return "", fmt.Errorf("%w: %s", ErrURLNotFound, shortURL)
}

func (r *FileRepository) Ping() error {
	return nil
}

func (r *FileRepository) Close() error {
	return nil
}

func (r *FileRepository) GetListByUser(ctx context.Context, userID string) ([]models.GetURLResponse, error) {
	var result []models.GetURLResponse
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("repository: %w", ctx.Err())
	default:
	}
	for key, val := range r.store {
		if userID == val[1] && val[2] == "f" {
			result = append(result, models.GetURLResponse{ShortURL: key, URL: val[0]})
		}
	}
	return result, nil
}

func (r *FileRepository) DeleteUserURLsBatch(ctx context.Context, userID string, data []string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	for _, short := range data {
		row, ok := r.store[short]
		if ok {
			if row[1] == userID {
				row[2] = "t"
			}
			// todo исправить в файле
		}
	}
	return nil
}

func (r *FileRepository) ValidateUniqueUser(ctx context.Context, userID string) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("repository: %w", ctx.Err())
	default:
	}
	for _, val := range r.store {
		if val[1] == userID {
			return fmt.Errorf("repository: %w: %s", ErrUserExists, userID)
		}
	}
	return nil
}
