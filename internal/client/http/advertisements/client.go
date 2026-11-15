package advertisements

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/advertisements/models"
)

// HttpPostAdvertisement отправляет форму и позволяет пропустить указанные поля.
func HttpPostAdvertisement(t *testing.T, token string, data models.CreateRequest, omittedFields ...string) *http.Response {
	t.Helper()
	for _, photo := range data.Photos {
		require.FileExists(t, photo, "не найден файл фотографии: %s", photo)
	}
	request := apiRunner.GetRunner().Auth(token).Create().Post(path.AdvertisementPath)
	fields := map[string]string{
		"title":       data.Title,
		"description": data.Description,
		"price":       strconv.FormatFloat(data.Price, 'f', -1, 64),
		"quantity":    strconv.Itoa(data.Quantity),
	}
	for name, value := range fields {
		if !slices.Contains(omittedFields, name) {
			request.MultipartFormData(name, value)
		}
	}
	return request.MultipartFile("photos", data.Photos...).Expect(t).End().Response
}

// HttpGetAdvertisement получает объявление по id.
func HttpGetAdvertisement(t *testing.T, id string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Create().Get(path.AdvertisementPath).
		Query("id", id).Expect(t).End().Response
}

// HttpGetAdvertisements ищет объявления по названию с ограничением времени.
func HttpGetAdvertisements(t *testing.T, ctx context.Context, title string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Create().Get(path.AdvertisementsPath).
		WithContext(ctx).Query("search", title).Expect(t).End().Response
}

// HttpGetPhotos получает список фотографий объявления.
func HttpGetPhotos(t *testing.T, token, id string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().
		Get(fmt.Sprintf(path.AdvertisementPhotosPath, id)).Expect(t).End().Response
}

// HttpGetPhoto скачивает файл из хранилища без передачи токена TestBoard.
func HttpGetPhoto(t *testing.T, photoURL string) *http.Response {
	t.Helper()
	client := &http.Client{Timeout: apiRunner.RequestTimeout}
	response, err := client.Get(photoURL)
	require.NoError(t, err, "не удалось скачать фотографию")
	return response
}

// HttpDeleteAdvertisement удаляет объявление его владельца.
func HttpDeleteAdvertisement(t *testing.T, token, id string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().Delete(path.AdvertisementPath).
		Query("id", id).Expect(t).End().Response
}
