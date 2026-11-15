package favorites

import (
	"net/http"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/favorites"
)

// Negative: Проверяем запрет добавления собственного объявления в избранное
func (s *FavoritesSuite) TestAddOwnAdvertisement() {
	t := s.T()
	// Arrange – создаём объявление и используем токен его владельца.
	t.Log("Подготовка собственного объявления")
	created := s.createAdvertisement(s.ownerData.Token)

	// Act – добавляем своё объявление в избранное.
	t.Log("Добавляем своё объявление и ожидаем HTTP 403")
	result := favorites.AddExpectError(t, s.ownerData.Token, created.ID, http.StatusForbidden)

	// Assert – проверяем ошибку и избранное именно владельца объявления.
	require.Equal(t, "forbidden", result.Error, "неверный код ошибки")
	require.Equal(t, "Cannot add own advertisement to favorites", result.Message, "неверное сообщение об ошибке")
	require.NotContains(t, s.favoriteIDs(s.ownerData.Token), created.ID, "своё объявление появилось в избранном")
	t.Log("Добавление собственного объявления запрещено")
}
