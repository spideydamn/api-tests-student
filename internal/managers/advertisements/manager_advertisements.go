package advertisements

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
)

const (
	searchTimeout  = 15 * time.Second
	searchInterval = 500 * time.Millisecond
)

// CreateWithCleanup создаёт объявление и регистрирует его удаление после теста.
func CreateWithCleanup(t *testing.T, token string, data models.CreateRequest) models.Advertisement {
	t.Helper()
	response := advertisements.HttpPostAdvertisement(t, token, data)
	body := readBody(t, response)
	registerCleanup(t, token, body)
	checkStatus(t, response, http.StatusCreated)
	require.Contains(t, response.Header.Get("Content-Type"), "application/json", "ожидался ответ в формате JSON")
	var created models.Advertisement
	require.NoError(t, json.Unmarshal(body, &created), "не удалось разобрать созданное объявление")
	require.NotEmpty(t, created.ID, "сервер не вернул id созданного объявления")
	return created
}

// CreateExpectError отправляет запрос и разбирает ошибку, очищая неожиданно созданное объявление.
func CreateExpectError(t *testing.T, requestToken, cleanupToken string, data models.CreateRequest, expectedStatus int, omittedFields ...string) models.ErrorResponse {
	t.Helper()
	response := advertisements.HttpPostAdvertisement(t, requestToken, data, omittedFields...)
	body := readBody(t, response)
	id := registerCleanup(t, cleanupToken, body)
	checkStatus(t, response, expectedStatus)
	require.Empty(t, id, "при отклонённом запросе сервер вернул id объявления")
	require.Contains(t, response.Header.Get("Content-Type"), "application/json", "ожидался ответ в формате JSON")
	var result models.ErrorResponse
	require.NoError(t, json.Unmarshal(body, &result), "не удалось разобрать ошибку создания объявления")
	return result
}

// registerCleanup регистрирует удаление объявления до проверки остальных полей ответа.
func registerCleanup(t *testing.T, token string, body []byte) string {
	t.Helper()
	id := gjson.GetBytes(body, "id").String()
	if id != "" {
		t.Cleanup(func() {
			t.Logf("очистка: удаляем объявление %s", id)
			cleanupAdvertisement(t, token, id)
		})
	}
	return id
}

// cleanupAdvertisement проверяет наличие объявления, чтобы не удалять его повторно.
func cleanupAdvertisement(t *testing.T, token, id string) {
	t.Helper()
	response := advertisements.HttpGetAdvertisement(t, id)
	readBody(t, response)
	if response.StatusCode == http.StatusNotFound {
		t.Log("очистка: объявление уже удалено")
		return
	}
	checkStatus(t, response, http.StatusOK)
	Delete(t, token, id)
}

// Get получает и разбирает объявление по id.
func Get(t *testing.T, id string) models.Advertisement {
	t.Helper()
	var result models.Advertisement
	readJSON(t, advertisements.HttpGetAdvertisement(t, id), &result)
	return result
}

// WaitInSearch ждёт появления конкретного id, повторяя только запрос поиска.
func WaitInSearch(t *testing.T, title, id string) models.Advertisement {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), searchTimeout)
	defer cancel()
	for {
		var result models.AdvertisementList
		readJSON(t, advertisements.HttpGetAdvertisements(t, ctx, title), &result)
		for _, item := range result.Items {
			if item.ID == id {
				return item
			}
		}
		select {
		case <-ctx.Done():
			t.Fatalf("объявление %s не появилось в поиске по названию %q за %s", id, title, searchTimeout)
		case <-time.After(searchInterval):
		}
	}
}

// GetPhotos получает фотографии конкретного объявления.
func GetPhotos(t *testing.T, token, id string) []models.Photo {
	t.Helper()
	var photos []models.Photo
	readJSON(t, advertisements.HttpGetPhotos(t, token, id), &photos)
	return photos
}

// DownloadPhoto получает содержимое загруженного файла.
func DownloadPhoto(t *testing.T, photoURL string) []byte {
	t.Helper()
	response := advertisements.HttpGetPhoto(t, photoURL)
	body := readBody(t, response)
	checkStatus(t, response, http.StatusOK)
	require.NotEmpty(t, body, "сервер вернул пустую фотографию")
	return body
}

// Delete удаляет созданное тестом объявление и проверяет результат очистки.
func Delete(t *testing.T, token, id string) {
	t.Helper()
	response := advertisements.HttpDeleteAdvertisement(t, token, id)
	readBody(t, response)
	checkStatus(t, response, http.StatusNoContent)
}

// readJSON проверяет успешный ответ и преобразует JSON в модель.
func readJSON(t *testing.T, response *http.Response, result any) {
	t.Helper()
	body := readBody(t, response)
	checkStatus(t, response, http.StatusOK)
	require.Contains(t, response.Header.Get("Content-Type"), "application/json", "ожидался ответ в формате JSON")
	require.NoError(t, json.Unmarshal(body, result), "не удалось разобрать JSON-ответ")
}

// readBody читает и закрывает тело ответа.
func readBody(t *testing.T, response *http.Response) []byte {
	t.Helper()
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	require.NoError(t, err, "не удалось прочитать тело ответа")
	return body
}

// checkStatus проверяет и логирует статус без пароля и токена.
func checkStatus(t *testing.T, response *http.Response, expected int) {
	t.Helper()
	require.Equal(t, expected, response.StatusCode, "неожиданный HTTP-статус")
	t.Logf("результат запроса: HTTP %d", response.StatusCode)
}
