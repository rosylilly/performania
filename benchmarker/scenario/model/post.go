package model

import (
	"bytes"
	"io"
	"mime/multipart"
	"time"
)

type Post struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Body      string    `json:"body"`
	Photo     *Image    `json:"photo"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (p *Post) CreateRequestBody() (string, io.Reader, error) {
	buf := &bytes.Buffer{}
	w := multipart.NewWriter(buf)

	w.WriteField("body", p.Body)

	if p.Photo != nil {
		photo, err := w.CreateFormFile("photo", "photo.png")
		if err != nil {
			return w.FormDataContentType(), nil, err
		}
		if _, err := io.Copy(photo, p.Photo.Reader()); err != nil {
			return w.FormDataContentType(), nil, err
		}
	}

	if err := w.Close(); err != nil {
		return w.FormDataContentType(), nil, err
	}

	return w.FormDataContentType(), buf, nil
}
