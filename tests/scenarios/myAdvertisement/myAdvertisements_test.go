package myAdvertisement

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"

	"api-tests-template/internal/managers/advertisements"
	advertisementModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/managers/myAdvertisements"
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

type TestSuite struct {
	suite.Suite
	loginData models.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Авторизация пользователя с кредами из переменных окружения и получение его параметров")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

// Positive: Проверяем получение собственных объявлений в профиле
func (s *TestSuite) TestGetMyAdvertisementsPositive() {
	t := s.T()
	// Arrange – создаём объявление только для этой проверки.
	t.Log("Подготовка собственного объявления")
	data := advertisementModels.CreateRequest{
		Title:       utils.RandomString(titleLength),
		Description: utils.RandomString(descriptionLength),
		Price:       advertisementPrice,
		Quantity:    advertisementQuantity,
		Photos:      []string{photoPath},
	}
	advertisements.CreateWithCleanup(t, s.loginData.Token, data)

	var advertisementsBody string
	s.Run("Получаем список собственных объявлений", func() {
		advertisementsBody = myAdvertisements.GetMyAdvertisements(s.T(), s.loginData.Token, http.StatusOK, 50)
	})

	var items gjson.Result
	s.Run("Проверяем, что у нас есть несколько добавленных ранее объявлений", func() {
		items = gjson.Get(advertisementsBody, "items")
		require.NotEmpty(s.T(), items.Array())
	})

	s.Run("Проверяем, что наши объявления имеют принадлежность к нашему пользователю и объявления имеют НЕ пустые параметры", func() {
		for _, item := range items.Array() {
			require.Equal(s.T(), s.loginData.User.Id, gjson.Get(item.String(), "user_id").String(), "user_id не совпадает")

			require.NotEmpty(s.T(), gjson.Get(item.String(), "id").String(), "id пустое")
			require.NotEmpty(s.T(), gjson.Get(item.String(), "created_at").String(), "created_at пустое")
			require.NotEmpty(s.T(), gjson.Get(item.String(), "updated_at").String(), "updated_at пустое")
			require.NotEmpty(s.T(), gjson.Get(item.String(), "title").String(), "Название пустое")
			require.NotEmpty(s.T(), gjson.Get(item.String(), "description").String(), "Описание пустое")
			require.True(s.T(), len(gjson.Get(item.String(), "photos").Array()) > 0, "Присутствуют пустые фотографии")
		}
	})
}

// Negative: Проверяем 401 ошибку при попытке получить собственные объявления без токена
func (s *TestSuite) TestGetMyAdvertisementsIncorrectToken() {
	var advertisementsBody string
	s.Run("Получаем список собственных объявлений с неправильным авторизационным токеном", func() {
		advertisementsBody = myAdvertisements.GetMyAdvertisements(s.T(), "incorrect_token", http.StatusUnauthorized, 50)
	})

	s.Run("Проверяем, что показывается правильная ошибка", func() {
		require.Equal(s.T(), "unauthorized", gjson.Get(advertisementsBody, "error").String(),
			"Некорректное значение error параметра")
		require.Equal(s.T(), "Invalid or expired token", gjson.Get(advertisementsBody, "message").String(),
			"Некорректное значение message параметра")
	})
}
