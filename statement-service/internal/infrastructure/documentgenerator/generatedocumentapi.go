package documentgenerator

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"

	"github.com/spf13/viper"
)

type GenerateDocumentApiInterface interface {
	GenerateFromHtml(html string) (string, error)
}

type GenerateDocumentApi struct {
	httpClient http.Client
}

func NewGenerateDocumentApi(httpClient http.Client) GenerateDocumentApiInterface {
	return &GenerateDocumentApi{
		httpClient: httpClient,
	}
}

func (api *GenerateDocumentApi) GenerateFromHtml(html string) (string, error) {
	baseUrl := viper.GetString("documentGenerator.baseUrl")
	url := fmt.Sprintf("%v/forms/chromium/convert/html", baseUrl)

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	part, err := writer.CreateFormFile("files", "index.html")
	if err != nil {
		slog.Error("error creating form file", "err", err)
		return "", err
	}

	_, err = io.Copy(part, bytes.NewBufferString(html))
	if err != nil {
		slog.Error("error copying HTML content", "err", err)
		return "", err
	}

	err = writer.Close()
	if err != nil {
		slog.Error("error closing writer", "err", err)
		return "", err
	}

	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		slog.Error("error creating HTTP request", "err", err)
		return "", err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("error making HTTP request", "err", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("failed to convert HTML to PDF, status code", "StatusCode", resp.StatusCode)
		return "", fmt.Errorf("failed to convert HTML to PDF, status code: %d", resp.StatusCode)
	}

	pdfBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("error saving PDF", "err", err)
		return "", err
	}

	base64PDF := base64.StdEncoding.EncodeToString(pdfBytes)

	slog.Info("pdf generated from html")

	return base64PDF, nil
}
