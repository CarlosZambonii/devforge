package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/CarlosZambonii/devforge/internal/crypto"
	"github.com/CarlosZambonii/devforge/internal/repository"
)

type URLService struct {
	repo   *repository.URLRepository
	crypto *crypto.AESCrypto
}

func NewURLService(repo *repository.URLRepository, crypto *crypto.AESCrypto) *URLService {
	return &URLService{repo: repo, crypto: crypto}
}

func (s *URLService) Shorten(ctx context.Context, rawURL string) (string, error) {
	if rawURL == "" {
		return "", errors.New("url nao pode ser vazia")
	}

	code := uuid.New().String()[:8]

	encrypted, err := s.crypto.Encrypt(rawURL)
	if err != nil {
		return "", err
	}

	if err := s.repo.Save(ctx, code, encrypted); err != nil {
		return "", err
	}

	return code, nil
}

func (s *URLService) Resolve(ctx context.Context, code string) (string, error) {
	if code == "" {
		return "", errors.New("code nao pode ser vazio")
	}

	encrypted, err := s.repo.Get(ctx, code)
	if err != nil {
		return "", errors.New("code nao encontrado")
	}

	rawURL, err := s.crypto.Decrypt(encrypted)
	if err != nil {
		return "", err
	}

	return rawURL, nil
}

func (s *URLService) Delete(ctx context.Context, code string) error {
	if code == "" {
		return errors.New("code nao pode ser vazio")
	}
	return s.repo.Delete(ctx, code)
}
