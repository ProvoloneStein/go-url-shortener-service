package handlersgrpc

import (
	"context"

	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/models"
	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/services"
	"github.com/ProvoloneStein/go-url-shortener-service/pkg/api/shorten"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	shorten.UnimplementedShortenServer
	logger *zap.Logger
	impl   *services.Service
}

func (h *Handler) CreateShortURL(ctx context.Context,
	req *shorten.CreateShortURLRequest) (*shorten.CreateShortURLResponse, error) {
	url, err := h.impl.CreateShortURL(ctx, req.GetUserId(), req.GetUrl())
	if err != nil {
		h.logger.Error("CreateShortURL service err", zap.Error(err))
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &shorten.CreateShortURLResponse{Result: url}, nil
}

func (h *Handler) BatchCreateShortURL(ctx context.Context,
	req *shorten.BatchCreateShortURLRequest) (*shorten.BatchCreateShortURLResponse, error) {
	items := []models.BatchCreateRequest{}
	for _, item := range req.GetItems() {
		items = append(items, models.BatchCreateRequest{URL: item.OriginalUrl, UUID: item.CorrelationId})
	}
	res, err := h.impl.BatchCreate(ctx, req.GetUserId(), items)
	if err != nil {
		h.logger.Error("BatchCreate service err", zap.Error(err))
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result := []*shorten.BatchCreateShortURLResponseData{}
	for _, item := range res {
		result = append(result, &shorten.BatchCreateShortURLResponseData{ShortUrl: item.ShortURL, CorrelationId: item.UUID})
	}

	return &shorten.BatchCreateShortURLResponse{Items: result}, nil
}

func (h *Handler) GetByShort(ctx context.Context,
	req *shorten.GetByShortRequest) (*shorten.GetByShortResponse, error) {
	full, err := h.impl.GetFullByID(ctx, req.GetUrl())
	if err != nil {
		h.logger.Error("GetFullByID service err", zap.Error(err))
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &shorten.GetByShortResponse{FullUrl: full}, nil
}

func (h *Handler) GetUserURLs(ctx context.Context,
	req *shorten.GetUserURLsRequest) (*shorten.GetUserURLsResponse, error) {
	urls, err := h.impl.GetListByUser(ctx, req.GetUserId())
	if err != nil {
		h.logger.Error("GetListByUser service err", zap.Error(err))
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	result := []*shorten.ShortenData{}
	for _, item := range urls {
		result = append(result, &shorten.ShortenData{ShortUrl: item.ShortURL, OriginalUrl: item.URL})
	}
	return &shorten.GetUserURLsResponse{Items: result}, nil
}

func (h *Handler) DeleteUserURLsBatch(ctx context.Context,
	req *shorten.DeleteUserURLsBatchRequest) (*shorten.DeleteUserURLsBatchResponse, error) {
	err := h.impl.DeleteUserURLsBatch(ctx, req.GetUserId(), req.GetUrls())
	if err != nil {
		h.logger.Error("DeleteUserURLsBatch service err", zap.Error(err))
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &shorten.DeleteUserURLsBatchResponse{}, nil
}

func NewServer(logger *zap.Logger, impl *services.Service) *Handler {
	return &Handler{logger: logger, impl: impl}
}
