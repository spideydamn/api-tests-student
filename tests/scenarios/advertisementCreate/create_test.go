package advertisementCreate

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"
	base "api-tests-template/tests"
)

const (
	titleLength           = 20
	descriptionLength     = 40
	advertisementPrice    = 137.50
	advertisementQuantity = 2
	photoPath             = "testdata/photo.png"
)

type CreateSuite struct {
	suite.Suite
	loginData authModels.LoginOkResponse
}

func TestCreateAdvertisement(t *testing.T) {
	suite.Run(t, new(CreateSuite))
}

// SetupSuite загружает настройки и один раз получает токен для набора.
func (s *CreateSuite) SetupSuite() {
	base.SetupSuite()
	base.Precondition("Авторизация пользователя с данными из переменных окружения")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
	require.NotEmpty(s.T(), s.loginData.Token, "сервер не вернул токен")
	require.NotEmpty(s.T(), s.loginData.User.Id, "сервер не вернул id пользователя")
}

// TearDownSuite завершает набор тестов через общий метод шаблона.
func (s *CreateSuite) TearDownSuite() {
	base.TearDownSuite()
}

// Positive: Проверяем создание объявления со всеми параметрами
func (s *CreateSuite) TestCreateWithAllFields() {
	t := s.T()

	// Arrange – готовим свои данные для этого запуска.
	data := newAdvertisementData()
	t.Logf("подготовка: название %q, цена %.2f, количество %d, фотографий %d", data.Title, data.Price, data.Quantity, len(data.Photos))

	// Act – создаём объявление; менеджер сразу регистрирует очистку.
	t.Log("действие: создаём объявление со всеми параметрами")
	created := advertisements.CreateWithCleanup(t, s.loginData.Token, data)
	t.Logf("результат создания: id %s", created.ID)

	// Assert – проверяем результат через чтение, поиск и загрузку фотографии.
	var actual models.Advertisement
	if !s.Run("Получаем объявление по id и проверяем сохранённые поля", func() {
		actual = s.checkAdvertisementFields(data, created.ID)
	}) {
		return
	}

	if !s.Run("Проверяем наличие объявления в поиске по полному названию", func() {
		found := advertisements.WaitInSearch(s.T(), data.Title, created.ID)
		require.Equal(s.T(), data.Title, found.Title, "в поиске другое название")
		s.T().Logf("результат поиска: найден id %s", found.ID)
	}) {
		return
	}

	s.Run("Проверяем фотографии объявления и доступность файла", func() {
		s.checkAdvertisementPhotos(actual, len(data.Photos))
	})
}

// checkAdvertisementFields сравнивает сохранённые поля с отправленными данными.
func (s *CreateSuite) checkAdvertisementFields(data models.CreateRequest, id string) models.Advertisement {
	t := s.T()
	t.Helper()
	actual := advertisements.Get(t, id)
	require.Equal(t, id, actual.ID, "вернулось другое объявление")
	require.Equal(t, s.loginData.User.Id, actual.UserID, "неверный владелец объявления")
	require.Equal(t, data.Title, actual.Title, "название отличается от отправленного")
	require.Equal(t, data.Description, actual.Description, "описание отличается от отправленного")
	require.Equal(t, data.Price, actual.Price, "цена отличается от отправленной")
	require.Equal(t, data.Quantity, actual.Quantity, "количество отличается от отправленного")
	t.Log("проверка полей завершена")
	return actual
}

// checkAdvertisementPhotos проверяет фотографии объявления и доступность файла.
func (s *CreateSuite) checkAdvertisementPhotos(actual models.Advertisement, expectedCount int) {
	t := s.T()
	t.Helper()
	photos := advertisements.GetPhotos(t, s.loginData.Token, actual.ID)
	require.Len(t, photos, expectedCount, "число фотографий отличается от загруженного")
	require.Len(t, actual.Photos, expectedCount, "в карточке неверное число фотографий")
	photo := photos[0]
	require.NotEmpty(t, photo.ID, "у фотографии нет id")
	require.Equal(t, actual.Photos[0].ID, photo.ID, "фотография не совпадает с карточкой")
	require.NotEmpty(t, photo.URL, "у фотографии нет ссылки")
	advertisements.DownloadPhoto(t, photo.URL)
	t.Log("результат: фотография доступна для скачивания")
}

// newAdvertisementData готовит корректные данные отдельно для каждого случая.
func newAdvertisementData() models.CreateRequest {
	return models.CreateRequest{
		Title:       utils.RandomString(titleLength),
		Description: utils.RandomString(descriptionLength),
		Price:       advertisementPrice,
		Quantity:    advertisementQuantity,
		Photos:      []string{photoPath},
	}
}
