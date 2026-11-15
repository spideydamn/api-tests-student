package favorites

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/favorites"
	"api-tests-template/internal/managers/favorites/models"
)

// Add добавляет объявление в избранное и возвращает созданную запись.
func Add(t *testing.T, token, adID string) models.Favorite {
	t.Helper()
	request := models.AddRequest{AdvertisementID: adID}
	var result models.Favorite
	readJSON(t, favorites.HttpPostFavorite(t, token, request), http.StatusCreated, &result)
	return result
}

// Get получает избранное авторизованного пользователя.
func Get(t *testing.T, token string) models.FavoriteList {
	t.Helper()
	var result models.FavoriteList
	readJSON(t, favorites.HttpGetFavorites(t, token), http.StatusOK, &result)
	require.NotNil(t, result.Items, "в ответе избранного отсутствует массив items")
	return result
}

// Remove удаляет запись и проверяет указанные статусы; при очистке допустим также 404.
func Remove(t *testing.T, token, adID string, expectedStatuses ...int) {
	t.Helper()
	response := favorites.HttpDeleteFavorite(t, token, adID)
	defer response.Body.Close()
	require.Contains(t, expectedStatuses, response.StatusCode, "неожиданный HTTP-статус удаления из избранного")
	t.Logf("результат удаления из избранного: HTTP %d", response.StatusCode)
}

// AddExpectError добавляет объявление с ошибочными условиями и разбирает ответ.
func AddExpectError(t *testing.T, token, adID string, expectedStatus int) models.ErrorResponse {
	t.Helper()
	request := models.AddRequest{AdvertisementID: adID}
	var result models.ErrorResponse
	readJSON(t, favorites.HttpPostFavorite(t, token, request), expectedStatus, &result)
	return result
}

// GetExpectError запрашивает избранное с ошибочными условиями и разбирает ответ.
func GetExpectError(t *testing.T, token string, expectedStatus int) models.ErrorResponse {
	t.Helper()
	var result models.ErrorResponse
	readJSON(t, favorites.HttpGetFavorites(t, token), expectedStatus, &result)
	return result
}

// RemoveExpectError удаляет запись с ошибочными условиями и разбирает ответ.
func RemoveExpectError(t *testing.T, token, adID string, expectedStatus int) models.ErrorResponse {
	t.Helper()
	var result models.ErrorResponse
	readJSON(t, favorites.HttpDeleteFavorite(t, token, adID), expectedStatus, &result)
	return result
}

// readJSON проверяет статус, читает и закрывает ответ, преобразует JSON в модель.
func readJSON(t *testing.T, response *http.Response, expectedStatus int, result any) {
	t.Helper()
	defer response.Body.Close()
	require.Equal(t, expectedStatus, response.StatusCode, "неожиданный HTTP-статус избранного")
	require.Contains(t, response.Header.Get("Content-Type"), "application/json", "ожидался ответ в формате JSON")
	body, err := io.ReadAll(response.Body)
	require.NoError(t, err, "не удалось прочитать ответ избранного")
	require.NoError(t, json.Unmarshal(body, result), "не удалось разобрать ответ избранного")
	t.Logf("результат запроса избранного: HTTP %d", response.StatusCode)
}
