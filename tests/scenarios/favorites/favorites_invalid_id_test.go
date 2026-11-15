package favorites

import (
	"fmt"
	"net/http"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/managers/favorites/models"
)

const (
	invalidAdvertisementID = "incorrect_id"
	notFoundCode           = "not_found"
	invalidRequestCode     = "invalid_request"
)

// Negative: Проверяем несуществующий и некорректный id объявления в избранном
func (s *FavoritesSuite) TestFavoritesWithInvalidID() {
	cases := []struct {
		name                    string
		method                  string
		id                      string
		useDeletedAdvertisement bool
		expectedStatus          int
		expectedCode            string
		expectedMessage         string
	}{
		{
			name:                    "Добавление несуществующего объявления",
			method:                  http.MethodPost,
			useDeletedAdvertisement: true,
			expectedStatus:          http.StatusNotFound,
			expectedCode:            notFoundCode,
			expectedMessage:         "Advertisement not found",
		},
		{
			name:                    "Удаление несуществующего объявления",
			method:                  http.MethodDelete,
			useDeletedAdvertisement: true,
			expectedStatus:          http.StatusNotFound,
			expectedCode:            notFoundCode,
			expectedMessage:         "Favorite entry not found",
		},
		{
			name:            "Добавление с некорректным id",
			method:          http.MethodPost,
			id:              invalidAdvertisementID,
			expectedStatus:  http.StatusBadRequest,
			expectedCode:    invalidRequestCode,
			expectedMessage: fmt.Sprintf("invalid UUID length: %d", len(invalidAdvertisementID)),
		},
		{
			name:            "Удаление с некорректным id",
			method:          http.MethodDelete,
			id:              invalidAdvertisementID,
			expectedStatus:  http.StatusBadRequest,
			expectedCode:    invalidRequestCode,
			expectedMessage: "Invalid advertisement ID",
		},
		{
			name:            "Добавление с пустым id",
			method:          http.MethodPost,
			expectedStatus:  http.StatusBadRequest,
			expectedCode:    invalidRequestCode,
			expectedMessage: "invalid UUID length: 0",
		},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			t := s.T()
			// Arrange – для отсутствующего id создаём и удаляем своё объявление.
			t.Log("Подготовка id для отдельного случая")
			id := tc.id
			if tc.useDeletedAdvertisement {
				created := s.createAdvertisement(s.loginData.Token)
				advertisements.Delete(t, s.ownerData.Token, created.ID)
				id = created.ID
			}

			// Act – отправляем запрос с несуществующим или некорректным id.
			t.Logf("Отправляем запрос и ожидаем HTTP %d", tc.expectedStatus)
			var result models.ErrorResponse
			switch tc.method {
			case http.MethodPost:
				result = favorites.AddExpectError(t, s.loginData.Token, id, tc.expectedStatus)
			case http.MethodDelete:
				result = favorites.RemoveExpectError(t, s.loginData.Token, id, tc.expectedStatus)
			default:
				t.Fatalf("неизвестный HTTP-метод: %s", tc.method)
			}

			// Assert – проверяем код ошибки и точное сообщение.
			require.Equal(t, tc.expectedCode, result.Error, "неверный код ошибки")
			require.Equal(t, tc.expectedMessage, result.Message, "неверное сообщение об ошибке")
			t.Log("Запрос отклонён с ожидаемой ошибкой")
		})
	}
}
