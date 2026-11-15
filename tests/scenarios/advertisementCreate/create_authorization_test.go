package advertisementCreate

import (
	"net/http"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements"
)

// Negative: Проверяем отказ в создании объявления без действующего токена
func (s *CreateSuite) TestCreateWithoutAuthorization() {
	cases := []struct {
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

	for _, tc := range cases {
		s.Run(tc.name, func() {
			t := s.T()
			// Arrange – все поля объявления корректны.
			t.Log("Подготовка корректных данных объявления")
			data := newAdvertisementData()

			// Act – меняем только авторизацию запроса.
			t.Log("Отправляем запрос и ожидаем HTTP 401")
			result := advertisements.CreateExpectError(t, tc.token, s.loginData.Token, data, http.StatusUnauthorized)

			// Assert – проверяем код ошибки и сообщение.
			require.Equal(t, "unauthorized", result.Error, "неверный код ошибки")
			require.Equal(t, tc.message, result.Message, "неверное сообщение об ошибке")
			t.Log("Создание без авторизации отклонено с ожидаемой ошибкой")
		})
	}
}
