package repositories

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ProvoloneStein/go-url-shortener-service/configs"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/models"
	"go.uber.org/zap"
	"io"
	"net/url"
	"os"
	"strconv"
)

type ShorterRecord struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	UserID      string `json:"user_id"`
}

type FileRepository struct {
	cfg    configs.AppConfig
	logger *zap.Logger
	file   *os.File
	writer *bufio.Writer
	reader *bufio.Reader
	store  map[string][]string
	uuid   int
}

func NewFileRepository(cfg configs.AppConfig, logger *zap.Logger, file *os.File) (*FileRepository, error) {

	repo := FileRepository{
		cfg:    cfg,
		logger: logger,
		file:   file,
		store:  make(map[string][]string),
		writer: bufio.NewWriter(file),
		reader: bufio.NewReader(file),
		uuid:   0,
	}

	for {
		record, err := repo.ReadString()
		if err != nil {
			return nil, err
		}
		if record == nil {
			break
		}
		repo.store[record.ShortURL] = []string{record.OriginalURL, record.UserID}
		repo.uuid, err = strconv.Atoi(record.UUID)
		if err != nil {
			return nil, fmt.Errorf("ошибка при получения uuid записи: %s", err)
		}
	}
	repo.uuid += 1
	return &repo, nil
}

func (r *FileRepository) ReadString() (*ShorterRecord, error) {
	// читаем данные до символа переноса строки
	data, err := r.reader.ReadBytes('\n')
	if err == io.EOF {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	// преобразуем данные из JSON-представления в структуру
	record := ShorterRecord{}
	err = json.Unmarshal(data, &record)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *FileRepository) WriteString(record ShorterRecord) error {
	data, err := json.Marshal(&record)
	if err != nil {
		return err
	}

	// записываем событие в буфер
	if _, err := r.writer.Write(data); err != nil {
		return err
	}

	// добавляем перенос строки
	if err := r.writer.WriteByte('\n'); err != nil {
		return err
	}

	r.uuid += 1

	// записываем буфер в файл
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
		return "", ctx.Err()
	default:
	}
	if err := r.validateUniqueShortURL(ctx, shortURL); err != nil {
		return "", err
	}
	for key, val := range r.store {
		if val[0] == fullURL {
			return key, ErrorUniqueViolation
		}
	}
	r.store[shortURL] = []string{fullURL, userID}
	if err := r.WriteString(ShorterRecord{strconv.Itoa(r.uuid), shortURL, fullURL, userID}); err != nil {
		delete(r.store, shortURL)
		return shortURL, err
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
			return []models.BatchCreateResponse{models.BatchCreateResponse{ShortURL: val.ShortURL, UUID: val.UUID}}, err
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
			return response, err
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
		return nil, ctx.Err()
	default:
	}
	for key, val := range r.store {
		if userID == val[1] {
			result = append(result, models.GetURLResponse{ShortURL: key, URL: val[0]})
		}
	}
	return result, nil
}

func (r *FileRepository) ValidateUniqueUser(ctx context.Context, userID string) error {

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	for _, val := range r.store {
		if val[1] == userID {
			return fmt.Errorf("%w: %s", ErrUserExists, userID)
		}
	}
	return nil
}
