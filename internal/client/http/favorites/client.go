package favorites

import (
	"fmt"
	"net/http"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/favorites/models"
)

const favoritesPageSize = "50"

// HttpPostFavorite добавляет объявление в избранное пользователя.
func HttpPostFavorite(t *testing.T, token string, data models.AddRequest) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().Post(path.FavoritesPath).
		ContentType("application/json").JSON(data).Expect(t).End().Response
}

// HttpGetFavorites получает список с учётом лимита в 50 избранных объявлений.
func HttpGetFavorites(t *testing.T, token string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().Get(path.FavoritesPath).
		Query("page_size", favoritesPageSize).Expect(t).End().Response
}

// HttpDeleteFavorite удаляет объявление из избранного по его id.
func HttpDeleteFavorite(t *testing.T, token, adID string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().Delete(fmt.Sprintf(path.FavoritePath, adID)).
		Expect(t).End().Response
}
