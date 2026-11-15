package favorites

import (
	"net/http"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/managers/favorites/models"
)

// Negative: Проверяем все ручки избранного без действующего токена
func (s *FavoritesSuite) TestFavoritesWithoutAuthorization() {
	authorizationCases := []struct {
		name    string
		token   string
		message string
	}{
		{
			name:    "Без токена",
			message: "Authorization header required",
		},
		{
			name:    "С неверным токеном",
			token:   "incorrect_token",
			message: "Invalid or expired token",
		},
	}
	operations := []struct {
		name        string
		method      string
		inFavorites bool
	}{
		{
			name:   "Получение списка",
			method: http.MethodGet,
		},
		{
			name:   "Добавление объявления",
			method: http.MethodPost,
		},
		{
			name:        "Удаление объявления",
			method:      http.MethodDelete,
			inFavorites: true,
		},
	}

	for _, operation := range operations {
		for _, tc := range authorizationCases {
			s.Run(operation.name+"/"+tc.name, func() {
				t := s.T()
				// Arrange – готовим объявление только для добавления и удаления.
				t.Log("Подготовка запроса без действующего токена")
				var adID string
				if operation.method != http.MethodGet {
					created := s.createAdvertisement(s.loginData.Token)
					adID = created.ID
					if operation.inFavorites {
						favorites.Add(t, s.loginData.Token, adID)
						require.Contains(t, s.favoriteIDs(s.loginData.Token), adID, "не удалось подготовить запись избранного")
					}
				}

				// Act – меняем только токен запроса.
				t.Log("Отправляем запрос и ожидаем HTTP 401")
				var result models.ErrorResponse
				switch operation.method {
				case http.MethodGet:
					result = favorites.GetExpectError(t, tc.token, http.StatusUnauthorized)
				case http.MethodPost:
					result = favorites.AddExpectError(t, tc.token, adID, http.StatusUnauthorized)
				case http.MethodDelete:
					result = favorites.RemoveExpectError(t, tc.token, adID, http.StatusUnauthorized)
				default:
					t.Fatalf("неизвестный HTTP-метод: %s", operation.method)
				}

				// Assert – проверяем ошибку и сохранение исходного состояния.
				require.Equal(t, "unauthorized", result.Error, "неверный код ошибки")
				require.Equal(t, tc.message, result.Message, "неверное сообщение об ошибке")
				if operation.method != http.MethodGet {
					ids := s.favoriteIDs(s.loginData.Token)
					if operation.inFavorites {
						require.Contains(t, ids, adID, "запись исчезла после запроса без авторизации")
					} else {
						require.NotContains(t, ids, adID, "объявление добавлено без авторизации")
					}
				}
				t.Log("Запрос отклонён с ожидаемой ошибкой")
			})
		}
	}
}
