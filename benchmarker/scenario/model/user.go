package model

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/textproto"
	"net/url"
	"strings"
	"time"
)

type User struct {
	ID          string    `json:"id"`
	Login       string    `json:"login"`
	DisplayName string    `json:"display_name"`
	Icon        *Image    `json:"icon"`
	Cover       *Image    `json:"cover"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (u *User) SignupRequestBody() (string, io.Reader, error) {
	buf := &bytes.Buffer{}
	w := multipart.NewWriter(buf)

	w.WriteField("login", u.Login)
	w.WriteField("display_name", u.DisplayName)

	ih := make(textproto.MIMEHeader)
	ih.Set("Content-Disposition", `form-data; name="icon"; filename="icon.png"`)
	ih.Set("Content-Type", "image/png")
	icon, err := w.CreatePart(ih)
	if err != nil {
		return w.FormDataContentType(), nil, err
	}
	if _, err := io.Copy(icon, u.Icon.Reader()); err != nil {
		return w.FormDataContentType(), nil, err
	}

	ch := make(textproto.MIMEHeader)
	ch.Set("Content-Disposition", `form-data; name="cover"; filename="cover.png"`)
	ch.Set("Content-Type", "image/png")
	cover, err := w.CreatePart(ch)
	if err != nil {
		return w.FormDataContentType(), nil, err
	}
	if _, err := io.Copy(cover, u.Cover.Reader()); err != nil {
		return w.FormDataContentType(), nil, err
	}

	if err := w.Close(); err != nil {
		return w.FormDataContentType(), nil, err
	}

	return w.FormDataContentType(), buf, err
}

func (u *User) LoginRequestBody() (string, io.Reader, error) {
	values := url.Values{}
	values.Add("login", u.Login)

	return "application/x-www-form-urlencoded", strings.NewReader(values.Encode()), nil
}
