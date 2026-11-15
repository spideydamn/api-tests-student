package favorites

import (
	"fmt"
	"net/http"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/favorites"
)

const maxFavorites = 50

// Negative: Проверяем отказ при превышении лимита избранного пользователя
func (s *FavoritesSuite) TestFavoritesLimit() {
	t := s.T()
	// Arrange – учитываем существующие записи и заполняем только свободные места.
	t.Log("Подготовка избранного с максимальным числом записей")
	initial := favorites.Get(t, s.loginData.Token)
	require.LessOrEqual(t, initial.Total, maxFavorites, "исходное избранное уже превышает лимит")
	require.Len(t, initial.Items, initial.Total, "получена не вся страница избранного")
	for i := initial.Total; i < maxFavorites; i++ {
		created := s.createAdvertisement(s.loginData.Token)
		favorites.Add(t, s.loginData.Token, created.ID)
	}
	filled := favorites.Get(t, s.loginData.Token)
	require.Equal(t, maxFavorites, filled.Total, "не удалось заполнить избранное до лимита")
	require.Len(t, filled.Items, maxFavorites, "в списке меньше записей, чем ожидалось")
	extraAdvertisement := s.createAdvertisement(s.loginData.Token)

	// Act – пытаемся добавить ещё одно объявление после достижения лимита.
	t.Log("Превышаем лимит избранного и ожидаем HTTP 400")
	result := favorites.AddExpectError(t, s.loginData.Token, extraAdvertisement.ID, http.StatusBadRequest)

	// Assert – проверяем ошибку, размер списка и отсутствие нового объявления.
	require.Equal(t, "limit_reached", result.Error, "неверный код ошибки")
	require.Equal(t, fmt.Sprintf("Favorites limit of %d reached", maxFavorites), result.Message, "неверное сообщение об ошибке")
	actual := favorites.Get(t, s.loginData.Token)
	require.Equal(t, maxFavorites, actual.Total, "число записей изменилось после отклонённого добавления")
	require.Len(t, actual.Items, maxFavorites, "размер списка изменился после отклонённого добавления")
	for _, item := range actual.Items {
		require.NotEqual(t, extraAdvertisement.ID, item.Advertisement.ID, "объявление добавлено сверх лимита")
	}
	t.Log("Лимит соблюдён, дополнительное объявление не добавлено")
}
