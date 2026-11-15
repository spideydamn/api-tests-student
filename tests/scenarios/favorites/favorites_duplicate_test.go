package favorites

import (
	"net/http"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/favorites"
)

// Negative: Проверяем отказ при повторном добавлении объявления в избранное
func (s *FavoritesSuite) TestAddDuplicateFavorite() {
	t := s.T()
	// Arrange – создаём объявление и добавляем его в избранное.
	t.Log("Подготовка объявления, которое уже есть в избранном")
	created := s.createAdvertisement(s.loginData.Token)
	favorites.Add(t, s.loginData.Token, created.ID)
	require.Contains(t, s.favoriteIDs(s.loginData.Token), created.ID, "не удалось подготовить запись избранного")

	// Act – повторно добавляем то же объявление.
	t.Log("Повторяем добавление и ожидаем HTTP 409")
	result := favorites.AddExpectError(t, s.loginData.Token, created.ID, http.StatusConflict)

	// Assert – проверяем ошибку и отсутствие дубликата.
	require.Equal(t, "already_exists", result.Error, "неверный код ошибки")
	require.Equal(t, "Advertisement is already in favorites", result.Message, "неверное сообщение об ошибке")
	count := 0
	for _, id := range s.favoriteIDs(s.loginData.Token) {
		if id == created.ID {
			count++
		}
	}
	require.Equal(t, 1, count, "в избранном должна остаться ровно одна запись объявления")
	t.Log("Повторное добавление отклонено, дубликат не появился")
}
