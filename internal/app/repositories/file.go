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
	writer *bufio.Writer
	reader *bufio.Reader
	file   *os.File
	logger *zap.Logger
	store  map[string][3]string
	cfg    configs.AppConfig
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
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, defaultRepoErrWrapper(err)
		}
		if record == nil {
			break
		}
		repo.store[record.ShortURL] = [3]string{record.OriginalURL, record.UserID, record.Deleted}
		repo.uuid, err = strconv.Atoi(record.UUID)
		if err != nil {
			return nil, fmt.Errorf("ошибка при получения uuid записи: %w", err)
		}
	}
	repo.uuid++
	return &repo, nil
}

func (r *FileRepository) readString() (*ShorterRecord, error) {
	// Читаем данные до символа переноса строки
	data, err := r.reader.ReadBytes('\n')
	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении строки %w", err)
	}
	// Преобразуем данные из JSON-представления в структуру
	record := ShorterRecord{}
	err = json.Unmarshal(data, &record)
	if err != nil {
		return nil, defaultRepoErrWrapper(err)
	}
	return &record, nil
}

func (r *FileRepository) writeString(record ShorterRecord) error {
	data, err := json.Marshal(&record)
	if err != nil {
		return fmt.Errorf("json marshal error: %w", err)
	}

	// Записываем событие в буфер
	if _, err := r.writer.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("ошибка при записи строки %w", err)
	}
	r.uuid++

	// Записываем буфер в файл
	if err := r.writer.Flush(); err != nil {
		return fmt.Errorf("ошибка записи буфера в файл : %w", err)
	}
	return nil
}

func (r *FileRepository) validateUniqueShortURL(ctx context.Context, shortURL string) error {
	select {
	case <-ctx.Done():
		return defaultRepoErrWrapper(ctx.Err())
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
		return "", defaultRepoErrWrapper(ctx.Err())
	default:
	}
	if err := r.validateUniqueShortURL(ctx, shortURL); err != nil {
		return "", defaultRepoErrWrapper(err)
	}
	for key, val := range r.store {
		if val[0] == fullURL {
			return key, ErrUniqueViolation
		}
	}
	r.store[shortURL] = [3]string{fullURL, userID, "f"}
	if err := r.writeString(ShorterRecord{strconv.Itoa(r.uuid), shortURL, fullURL, userID, "f"}); err != nil {
		delete(r.store, shortURL)
		return shortURL, defaultRepoErrWrapper(err)
	}
	return shortURL, nil
}

func (r *FileRepository) BatchCreate(ctx context.Context,
	data []models.BatchCreateData) ([]models.BatchCreateResponse, error) {
	select {
	case <-ctx.Done():
		return nil, defaultRepoErrWrapper(ctx.Err())
	default:
	}
	for _, val := range data {
		if err := r.validateUniqueShortURL(ctx, val.ShortURL); err != nil {
			return []models.BatchCreateResponse{
				models.BatchCreateResponse{ShortURL: val.ShortURL, UUID: val.UUID},
			}, defaultRepoErrWrapper(err)
		}
	}

	response := make([]models.BatchCreateResponse, 0, len(data))
	for _, val := range data {
		shortURL, err := url.JoinPath(r.cfg.BaseURL, val.ShortURL)
		if err != nil {
			r.logger.Error("ошибка при формировании url", zap.Error(err))
			return response, fmt.Errorf("ошибка при формировании url: %w", err)
		}
		_, err = r.Create(ctx, val.UserID, val.URL, val.ShortURL)
		if err != nil && !errors.Is(err, ErrUniqueViolation) {
			r.logger.Error("ошибка при записи url", zap.Error(err))
			return response, defaultRepoErrWrapper(err)
		}
		response = append(response, models.BatchCreateResponse{ShortURL: shortURL, UUID: val.UUID})
	}
	return response, nil
}

func (r *FileRepository) GetByShort(ctx context.Context, shortURL string) (string, error) {
	select {
	case <-ctx.Done():
		return "", defaultRepoErrWrapper(ctx.Err())
	default:
	}
	data, ok := r.store[shortURL]
	if ok {
		if data[2] == "t" {
			return "", ErrDeleted
		}
		return data[0], nil
	}
	return "", errWithVal(ErrURLNotFound, shortURL)
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
		return nil, defaultRepoErrWrapper(ctx.Err())
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
		return defaultRepoErrWrapper(ctx.Err())
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

func (r *FileRepository) Stats(ctx context.Context) (models.StatsData, error) {
	var data models.StatsData
	select {
	case <-ctx.Done():
		return data, defaultRepoErrWrapper(ctx.Err())
	default:
	}
	return data, nil
}
