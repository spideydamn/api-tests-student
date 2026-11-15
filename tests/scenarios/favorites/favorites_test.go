package favorites

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/managers/advertisements"
	advertisementModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/utils"
	base "api-tests-template/tests"
)

const (
	titleLength           = 20
	descriptionLength     = 40
	advertisementPrice    = 100
	advertisementQuantity = 1
	photoPath             = "../advertisementCreate/testdata/photo.png"
)

type FavoritesSuite struct {
	suite.Suite
	ownerData authModels.LoginOkResponse
	loginData authModels.LoginOkResponse
}

func TestFavorites(t *testing.T) {
	suite.Run(t, new(FavoritesSuite))
}

// SetupSuite загружает настройки и авторизует владельца объявления и второго пользователя.
func (s *FavoritesSuite) SetupSuite() {
	base.SetupSuite()
	t := s.T()
	base.Precondition("Авторизация владельца объявления и пользователя избранного")
	s.ownerData = auth.Login(t, os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
	s.loginData = auth.Login(t, os.Getenv("FAVORITES_LOGIN"), os.Getenv("FAVORITES_PASSWORD"))
	require.NotEmpty(t, s.ownerData.Token, "сервер не вернул токен владельца")
	require.NotEmpty(t, s.loginData.Token, "сервер не вернул токен пользователя избранного")
	require.NotEmpty(t, s.ownerData.User.Id, "сервер не вернул id владельца")
	require.NotEmpty(t, s.loginData.User.Id, "сервер не вернул id пользователя избранного")
	require.NotEqual(t, s.ownerData.User.Id, s.loginData.User.Id, "для избранного нужны два разных пользователя")
}

// TearDownSuite завершает набор через общий метод шаблона.
func (s *FavoritesSuite) TearDownSuite() {
	base.TearDownSuite()
}

// Positive: Проверяем добавление, удаление из избранного и удаление самого объявления
func (s *FavoritesSuite) TestFavoritesLifecycle() {
	t := s.T()

	// Arrange – создаём объявление другого пользователя и регистрируем очистку до проверок.
	t.Log("Подготовка объявления другого пользователя")
	created := s.createAdvertisement(s.loginData.Token)

	// Act и Assert – выполняем цикл и проверяем состояние после каждого действия.
	if !s.Run("Добавляем объявление в избранное и проверяем список", func() {
		t := s.T()
		t.Log("Добавляем объявление в избранное")
		added := favorites.Add(t, s.loginData.Token, created.ID)
		require.NotEmpty(t, added.FavoriteID, "сервер не вернул id записи избранного")
		require.Equal(t, created.ID, added.Advertisement.ID, "в избранное добавлено другое объявление")
		require.Contains(t, s.favoriteIDs(s.loginData.Token), created.ID, "добавленное объявление отсутствует в избранном")
		t.Log("Объявление найдено в избранном")
	}) {
		return
	}

	if !s.Run("Удаляем объявление из избранного и проверяем список", func() {
		t := s.T()
		t.Log("Удаляем объявление из избранного")
		favorites.Remove(t, s.loginData.Token, created.ID, http.StatusNoContent)
		require.NotContains(t, s.favoriteIDs(s.loginData.Token), created.ID, "объявление осталось в избранном после удаления")
		t.Log("Объявление отсутствует в избранном")
	}) {
		return
	}

	if !s.Run("Повторно добавляем объявление в избранное", func() {
		t := s.T()
		t.Log("Повторно добавляем объявление в избранное перед удалением самого объявления")
		favorites.Add(t, s.loginData.Token, created.ID)
		require.Contains(t, s.favoriteIDs(s.loginData.Token), created.ID, "объявление не появилось в избранном перед удалением")
		t.Log("Объявление снова найдено в избранном")
	}) {
		return
	}

	s.Run("Удаляем само объявление и проверяем избранное", func() {
		t := s.T()
		t.Log("Удаляем объявление от имени владельца")
		advertisements.Delete(t, s.ownerData.Token, created.ID)
		require.NotContains(t, s.favoriteIDs(s.loginData.Token), created.ID, "удалённое объявление осталось в избранном")
		t.Log("После удаления объявления запись исчезла из избранного")
	})
}

// createAdvertisement создаёт объявление владельца и регистрирует очистку для указанного пользователя избранного.
func (s *FavoritesSuite) createAdvertisement(favoritesToken string) advertisementModels.Advertisement {
	t := s.T()
	t.Helper()
	data := advertisementModels.CreateRequest{
		Title:       utils.RandomString(titleLength),
		Description: utils.RandomString(descriptionLength),
		Price:       advertisementPrice,
		Quantity:    advertisementQuantity,
		Photos:      []string{photoPath},
	}
	created := advertisements.CreateWithCleanup(t, s.ownerData.Token, data)
	t.Cleanup(func() {
		t.Logf("очистка: удаляем объявление %s из избранного", created.ID)
		favorites.Remove(t, favoritesToken, created.ID, http.StatusNoContent, http.StatusNotFound)
	})
	return created
}

// favoriteIDs получает id объявлений из избранного указанного пользователя.
func (s *FavoritesSuite) favoriteIDs(token string) []string {
	t := s.T()
	t.Helper()
	list := favorites.Get(t, token)
	ids := make([]string, 0, len(list.Items))
	for _, item := range list.Items {
		require.NotEmpty(t, item.Advertisement.ID, "у объявления в избранном отсутствует id")
		ids = append(ids, item.Advertisement.ID)
	}
	return ids
}
