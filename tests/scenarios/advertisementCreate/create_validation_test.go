package advertisementCreate

import (
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
)

const (
	maxTitleLength             = 50
	maxDescriptionLength       = 500
	maxPrice                   = 1_000_000
	maxQuantity                = 100
	maxPhotos                  = 3
	titleRequiredMessage       = "title is required"
	descriptionRequiredMessage = "description is required"
	priceRangeMessage          = "price must be between 0 and 1000000 inclusive"
	quantityRangeMessage       = "quantity must be between 1 and 100 inclusive"
)

// Negative: Проверяем отказ в создании объявления с некорректными параметрами
func (s *CreateSuite) TestCreateWithInvalidFields() {
	cases := []struct {
		name         string
		change       func(*models.CreateRequest)
		omittedField string
		message      string
	}{
		{
			name:         "Отсутствует название",
			omittedField: "title",
			message:      titleRequiredMessage,
		},
		{
			name:         "Отсутствует описание",
			omittedField: "description",
			message:      descriptionRequiredMessage,
		},
		{
			name:    "Отсутствуют фотографии",
			change:  func(data *models.CreateRequest) { data.Photos = nil },
			message: "at least 1 photo is required",
		},
		{
			name:    "Пустое название",
			change:  func(data *models.CreateRequest) { data.Title = "" },
			message: titleRequiredMessage,
		},
		{
			name:    "Пустое описание",
			change:  func(data *models.CreateRequest) { data.Description = "" },
			message: descriptionRequiredMessage,
		},
		{
			name:    "Название длиннее допустимого на один символ",
			change:  func(data *models.CreateRequest) { data.Title = strings.Repeat("a", maxTitleLength+1) },
			message: fmt.Sprintf("title must not exceed %d characters provided %d", maxTitleLength, maxTitleLength+1),
		},
		{
			name:    "Описание длиннее допустимого на один символ",
			change:  func(data *models.CreateRequest) { data.Description = strings.Repeat("a", maxDescriptionLength+1) },
			message: fmt.Sprintf("description must be between 1 and %d characters", maxDescriptionLength),
		},
		{
			name:    "Отрицательная цена",
			change:  func(data *models.CreateRequest) { data.Price = -1 },
			message: priceRangeMessage,
		},
		{
			name:    "Цена выше максимальной",
			change:  func(data *models.CreateRequest) { data.Price = maxPrice + 1 },
			message: priceRangeMessage,
		},
		{
			name:    "Нулевое количество",
			change:  func(data *models.CreateRequest) { data.Quantity = 0 },
			message: quantityRangeMessage,
		},
		{
			name:    "Количество выше максимального",
			change:  func(data *models.CreateRequest) { data.Quantity = maxQuantity + 1 },
			message: quantityRangeMessage,
		},
		{
			name:    "Фотографий больше допустимого",
			change:  func(data *models.CreateRequest) { data.Photos = slices.Repeat([]string{photoPath}, maxPhotos+1) },
			message: fmt.Sprintf("maximum %d photos allowed", maxPhotos),
		},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			t := s.T()
			// Arrange – меняем одно условие в корректном запросе.
			t.Log("Подготовка данных для отдельного случая")
			data := newAdvertisementData()
			if tc.change != nil {
				tc.change(&data)
			}

			// Act – отправляем запрос с некорректными данными.
			t.Log("Отправляем запрос и ожидаем HTTP 400")
			result := advertisements.CreateExpectError(t, s.loginData.Token, s.loginData.Token, data, http.StatusBadRequest, tc.omittedField)

			// Assert – проверяем код ошибки и точное сообщение.
			require.Equal(t, "invalid_request", result.Error, "неверный код ошибки")
			require.Equal(t, tc.message, result.Message, "неверное сообщение об ошибке")
			t.Log("Запрос отклонён, код ошибки и сообщение совпадают с ожидаемыми")
		})
	}
}
