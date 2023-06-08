package repositories

import (
	"bufio"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"io"
	"os"
	"strconv"
)

type ShorterRecord struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type FileRepository struct {
	logger *zap.Logger
	file   *os.File
	writer *bufio.Writer
	reader *bufio.Reader
	store  map[string]string
	uuid   int
}

func NewFileRepository(logger *zap.Logger, file *os.File) *FileRepository {

	repo := FileRepository{
		logger: logger,
		file:   file,
		store:  make(map[string]string),
		writer: bufio.NewWriter(file),
		reader: bufio.NewReader(file),
		uuid:   0,
	}

	for {
		record, err := repo.ReadString()
		if err != nil {
			logger.Error("Ошибка при попытке иницилизации репозитория", zap.Error(err))
		}
		if record == nil {
			break
		}
		repo.store[record.ShortURL] = record.OriginalURL
		repo.uuid, err = strconv.Atoi(record.UUID)
		if err != nil {
			logger.Fatal("Ошибка при получения uuid записи", zap.Error(err))
		}
	}
	repo.uuid += 1
	return &repo
}

func (r *FileRepository) ReadString() (*ShorterRecord, error) {
	// читаем данные до символа переноса строки
	data, err := r.reader.ReadBytes('\n')
	if err == io.EOF {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		// преобразуем данные из JSON-представления в структуру
		record := ShorterRecord{}
		err = json.Unmarshal(data, &record)
		if err != nil {
			return nil, err
		}
		return &record, nil
	}
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

func (r *FileRepository) Create(fullURL string) (string, error) {
	var shortURL string
	for {
		shortURL = randomString()
		if _, ok := r.store[shortURL]; !ok {
			r.store[shortURL] = fullURL
			if err := r.WriteString(ShorterRecord{strconv.Itoa(r.uuid), shortURL, fullURL}); err != nil {
				return "", err
			}
			return shortURL, nil
		}
	}
}

func (r *FileRepository) GetByShort(shortURL string) (string, error) {
	fullURL, ok := r.store[shortURL]
	if ok {
		return fullURL, nil
	}
	return "", errors.New("url not found")
}
