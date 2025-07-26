package handlers

import (
	"bytes"
	"emtest/api-service/db"
	"emtest/api-service/subscription"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type HandlersTestSuite struct {
	suite.Suite
	app      *fiber.App
	testDB   *gorm.DB
	pool     *dockertest.Pool
	resource *dockertest.Resource
}

type TestUpdateSubscription struct {
	ServiceName string
	Price       int
	UserId      uuid.UUID
	StartDate   string
}

func (suite *HandlersTestSuite) SetupSuite() {

	pool, err := dockertest.NewPool("")
	if err != nil {
		logrus.Fatalf("Couldn't connect to docker: %s", err)
	}
	logrus.Info("Pool initialized...")
	suite.pool = pool

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "15",
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=test",
			"POSTGRES_DB=test_db",
			"listen_address = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		logrus.Fatalf("Couldn't start resource: %s", err)
	}
	logrus.Info("Resource initialized...")
	suite.resource = resource

	var testDB *gorm.DB
	dns := fmt.Sprintf(
		"host=localhost user=test password=secret dbname=test_db port=%s sslmode=disable TimeZone=UTC",
		resource.GetPort("5432/tcp"),
	)

	if err := pool.Retry(func() error {
		logrus.Info("Attempt to connect to db...")
		var err error
		testDB, err = gorm.Open(
			postgres.Open(dns),
			&gorm.Config{},
		)

		if err != nil {
			logrus.Infof("Database connection failed: %s", err)
			return err
		}

		sqlDB, err := testDB.DB()
		if err != nil {
			logrus.Infof("Failed to get sql.DB: %s", err)
			return err
		}

		if err := sqlDB.Ping(); err != nil {
			logrus.Infof("Database ping failed: %s", err)
			return err
		}

		logrus.Info("DB connected")
		return nil
	}); err != nil {
		logrus.Fatalf("Failed to connecto to docker: %s", err)
	}
	suite.testDB = testDB
	db.DB = testDB
	logrus.Info("Test database initialized...")

	err = suite.testDB.AutoMigrate(&subscription.Subscription{})
	if err != nil {
		logrus.Fatalf("Could not migrate database: %s", err)
	}

	suite.app = fiber.New()
	suite.app.Post("/api/v1/subscriptions", CreateSubscription)
	suite.app.Get("/api/v1/subscriptions", GetSubscriptions)
	suite.app.Put("/api/v1/subscriptions", UpdateSubscription)
	suite.app.Delete("/api/v1/subscriptions", DeleteSubscription)
	suite.app.Get("/api/v1/subscriptions/calculate", CalculateTotalCost)

	go func() {
		suite.app.Listen(":8081")
	}()

	time.Sleep(100 * time.Millisecond)
}

func (suite *HandlersTestSuite) TearDownSuite() {
	sqlDB, _ := suite.testDB.DB()
	sqlDB.Close()

	if suite.pool != nil && suite.resource != nil {
		if err := suite.pool.Purge(suite.resource); err != nil {
			logrus.Fatalf("Couldn't purge resource: %s", err)
		}
	}
}

func (suite *HandlersTestSuite) SetupTest() {
	if suite.testDB != nil {
		suite.testDB.Where("1 = 1").Delete(&subscription.Subscription{})
	}
}

func (suite *HandlersTestSuite) makeRequest(method, endpoint string, body interface{}) (*http.Response, error) {

	var reqBody io.Reader

	if body != nil {
		jsonBody, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, "http://localhost:8081"+endpoint, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	return client.Do(req)
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestHandlersSuite(t *testing.T) {
	suite.Run(t, new(HandlersTestSuite))
}

func (suite *HandlersTestSuite) TestCreateSubscription_Success() {

	sub := subscription.Subscription{
		ServiceName: "Test Yandex",
		Price:       900,
		UserId:      uuid.New(),
		StartDate:   "01-2025",
	}

	resp, err := suite.makeRequest("POST", "/api/v1/subscriptions", sub)
	assert.NoError(suite.T(), err)

	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	assert.NotEmpty(suite.T(), resp.Body, "Response body should not be empty")

	body, err := io.ReadAll(resp.Body)
	assert.NoError(suite.T(), err)

	var createdSub subscription.Subscription
	err = json.Unmarshal(body, &createdSub)

	assert.NoError(suite.T(), err)
	assert.NotEqual(suite.T(), uuid.Nil, createdSub.ID)
	assert.Equal(suite.T(), "Test Yandex", createdSub.ServiceName)
	assert.Equal(suite.T(), 900, createdSub.Price)
	assert.Equal(suite.T(), sub.UserId, createdSub.UserId)
	assert.Equal(suite.T(), "01-2025", createdSub.StartDate)
}

func (suite *HandlersTestSuite) TestCreateSubscription_InvalidJSON() {

	sub := subscription.Subscription{
		Price: 900,
	}

	resp, err := suite.makeRequest("POST", "/api/v1/subscriptions", sub)
	assert.NoError(suite.T(), err)

	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
}

func (suite *HandlersTestSuite) TestCreateSubscription_InvalidDateFormat() {

	sub := subscription.Subscription{
		ServiceName: "Test Yandex",
		Price:       900,
		UserId:      uuid.New(),
		StartDate:   "15-01-2025",
	}

	resp, err := suite.makeRequest("POST", "/api/v1/subscriptions", sub)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
}

func (suite *HandlersTestSuite) TestGetSubscriptions_All() {

	sub1 := subscription.Subscription{
		ServiceName: "Test Yandex",
		Price:       900,
		UserId:      uuid.New(),
		StartDate:   "01-2025",
	}
	sub2 := subscription.Subscription{
		ServiceName: "Test Yandex 2",
		Price:       600,
		UserId:      uuid.New(),
		StartDate:   "02-2025",
	}

	suite.testDB.Create(&sub1)
	suite.testDB.Create(&sub2)

	resp, err := suite.makeRequest("GET", "/api/v1/subscriptions", nil)
	assert.NoError(suite.T(), err)

	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	assert.NotEmpty(suite.T(), resp.Body, "Response body should not be empty")

	body, err := io.ReadAll(resp.Body)
	assert.NoError(suite.T(), err)

	var subs []subscription.Subscription
	err = json.Unmarshal(body, &subs)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), subs, 2)
}

func (suite *HandlersTestSuite) TestGetSubscriptions_All_Zero() {

	resp, err := suite.makeRequest("GET", "/api/v1/subscriptions", nil)
	assert.NoError(suite.T(), err)

	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	assert.NotEmpty(suite.T(), resp.Body, "Response body should not be empty")

	body, err := io.ReadAll(resp.Body)
	assert.NoError(suite.T(), err)

	var subs []subscription.Subscription
	err = json.Unmarshal(body, &subs)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), subs, 0)
}

func (suite *HandlersTestSuite) TestGetSubscriptions_ByID_Found() {

	sub := subscription.Subscription{
		ServiceName: "Test Yandex",
		Price:       900,
		UserId:      uuid.New(),
		StartDate:   "01-2025",
	}

	suite.testDB.Create(&sub)

	resp, err := suite.makeRequest("GET", fmt.Sprintf("/api/v1/subscriptions?id=%s", sub.ID.String()), nil)
	assert.NoError(suite.T(), err)

	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	assert.NotEmpty(suite.T(), resp.Body, "Response body should not be empty")

	body, err := io.ReadAll(resp.Body)
	assert.NoError(suite.T(), err)

	var returnedSub subscription.Subscription
	err = json.Unmarshal(body, &returnedSub)

	assert.NoError(suite.T(), err)
	assert.NotEqual(suite.T(), uuid.Nil, returnedSub.ID)
	assert.Equal(suite.T(), "Test Yandex", returnedSub.ServiceName)
	assert.Equal(suite.T(), 900, returnedSub.Price)
	assert.Equal(suite.T(), sub.UserId, returnedSub.UserId)
	assert.Equal(suite.T(), "01-2025", returnedSub.StartDate)
}

func (suite *HandlersTestSuite) TestGetSubscriptions_ByID_NotFound() {

	resp, err := suite.makeRequest("GET", fmt.Sprintf("/api/v1/subscriptions?id=%s", uuid.New().String()), nil)
	assert.NoError(suite.T(), err)

	defer resp.Body.Close()
	assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)
}

func (suite *HandlersTestSuite) TestUpdateSubscription_Success() {

	sub := subscription.Subscription{
		ServiceName: "Test Yandex",
		Price:       900,
		UserId:      uuid.New(),
		StartDate:   "01-2025",
	}
	result := suite.testDB.Create(&sub)
	assert.NoError(suite.T(), result.Error)
	assert.NotEqual(suite.T(), uuid.Nil, sub.ID)

	updateReq := map[string]interface{}{
		"service_name": "Test Yandex",
		"price":        100,
		"user_id":      sub.UserId,
		"start_date":   "02-2025",
	}

	resp, err := suite.makeRequest("PUT", fmt.Sprintf("/api/v1/subscriptions?id=%s", sub.ID.String()), updateReq)
	assert.NoError(suite.T(), err)

	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	assert.NotEmpty(suite.T(), resp.Body, "Response body should not be empty")

	body, err := io.ReadAll(resp.Body)
	assert.NoError(suite.T(), err)

	var updatedSub subscription.Subscription
	err = json.Unmarshal(body, &updatedSub)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 100, updatedSub.Price)
	assert.Equal(suite.T(), "02-2025", updatedSub.StartDate)

}

func (suite *HandlersTestSuite) TestUpdateSubscription_NotFound() {

	updatedData := subscription.Subscription{
		ServiceName: "Test Yandex",
		Price:       100,
		UserId:      uuid.New(),
		StartDate:   "02-2025",
	}

	resp, err := suite.makeRequest("PUT", fmt.Sprintf("/api/v1/subscriptions?id=%s", uuid.New().String()), updatedData)

	assert.NoError(suite.T(), err)

	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)
}

func (suite *HandlersTestSuite) TestDeleteSubscription_Success() {

	sub := subscription.Subscription{
		ServiceName: "Test Yandex",
		Price:       900,
		UserId:      uuid.New(),
		StartDate:   "01-2025",
	}
	suite.testDB.Create(&sub)

	resp, err := suite.makeRequest("DELETE", fmt.Sprintf("/api/v1/subscriptions?id=%s", sub.ID.String()), nil)
	assert.NoError(suite.T(), err)
	defer resp.Body.Close()

	var count int64
	suite.testDB.Model(&subscription.Subscription{}).Where("id = ?", sub.ID).Count(&count)
	assert.Equal(suite.T(), int64(0), count)
}

func (suite *HandlersTestSuite) TestDeleteSubscription_NotFound() {
	resp, err := suite.makeRequest("DELETE", fmt.Sprintf("/api/v1/subscriptions?id=%s", uuid.New().String()), nil)
	assert.NoError(suite.T(), err)

	defer resp.Body.Close()
	assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)
}

func (suite *HandlersTestSuite) TestCalcTotalCost_NoFilter() {
	subs := []subscription.Subscription{
		{ServiceName: "Test Yandex", Price: 100, UserId: uuid.New(), StartDate: "01-2024"},
		{ServiceName: "Test Google", Price: 200, UserId: uuid.New(), StartDate: "02-2024"},
		{ServiceName: "Test Yahoo", Price: 300, UserId: uuid.New(), StartDate: "03-2024"},
	}

	for _, sub := range subs {
		suite.testDB.Create(&sub)
	}

	resp, err := suite.makeRequest("GET", "/api/v1/subscriptions/calculate", nil)
	assert.NoError(suite.T(), err)

	defer resp.Body.Close()
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(suite.T(), err)

	var result map[string]int
	err = json.Unmarshal(body, &result)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 600, result["total"])

}

func (suite *HandlersTestSuite) TestCalcTotalCost_Filter_UserId() {
	userId := uuid.New()
	subs := []subscription.Subscription{
		{ServiceName: "Test Yandex", Price: 100, UserId: uuid.New(), StartDate: "01-2024"},
		{ServiceName: "Test Google", Price: 200, UserId: uuid.New(), StartDate: "02-2024"},
		{ServiceName: "Test Yahoo", Price: 300, UserId: uuid.New(), StartDate: "03-2024"},
		{ServiceName: "Test Yandex", Price: 101, UserId: userId, StartDate: "01-2024"},
		{ServiceName: "Test Google", Price: 202, UserId: userId, StartDate: "02-2024"},
		{ServiceName: "Test Yahoo", Price: 303, UserId: userId, StartDate: "03-2024"},
	}

	for _, sub := range subs {
		result := suite.testDB.Create(&sub)
		if result.Error != nil {
			suite.T().Fatalf("Failed to create subscription: %v", result.Error)
		}
	}

	resp, err := suite.makeRequest("GET", fmt.Sprintf("/api/v1/subscriptions/calculate?user_id=%s", userId.String()), nil)
	assert.NoError(suite.T(), err)

	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	assert.NotEmpty(suite.T(), resp.Body, "Response body should not be empty")

	body, err := io.ReadAll(resp.Body)
	assert.NoError(suite.T(), err)

	var result SuccessCostResponse

	err = json.Unmarshal(body, &result)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 606, result.Total)
}

func (suite *HandlersTestSuite) TestCalcTotalCost_Filter_UserIdStartDate() {
	userId := uuid.New()
	subs := []subscription.Subscription{
		{ServiceName: "Test Yandex", Price: 100, UserId: uuid.New(), StartDate: "01-2024"},
		{ServiceName: "Test Google", Price: 200, UserId: uuid.New(), StartDate: "02-2024"},
		{ServiceName: "Test Yahoo", Price: 300, UserId: uuid.New(), StartDate: "03-2024"},
		{ServiceName: "Test Yandex", Price: 101, UserId: userId, StartDate: "01-2024"},
		{ServiceName: "Test Google", Price: 202, UserId: userId, StartDate: "02-2024"},
		{ServiceName: "Test Yahoo", Price: 303, UserId: userId, StartDate: "03-2024"},
	}

	for _, sub := range subs {
		suite.testDB.Create(&sub)
	}

	resp, err := suite.makeRequest(
		"GET",
		fmt.Sprintf("/api/v1/subscriptions/calculate?user_id=%s&start_date=%s", userId.String(), "02-2024"),
		nil,
	)
	assert.NoError(suite.T(), err)

	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	assert.NotEmpty(suite.T(), resp.Body, "Response body should not be empty")

	body, err := io.ReadAll(resp.Body)
	assert.NoError(suite.T(), err)

	var result map[string]int
	err = json.Unmarshal(body, &result)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 505, result["total"])

}

func (suite *HandlersTestSuite) TestCalcTotalCost_Filter_UserIdStartEndDate() {
	userId := uuid.New()
	subs := []subscription.Subscription{
		{ServiceName: "Test Yandex", Price: 100, UserId: uuid.New(), StartDate: "01-2024"},
		{ServiceName: "Test Google", Price: 200, UserId: uuid.New(), StartDate: "02-2024"},
		{ServiceName: "Test Yahoo", Price: 300, UserId: uuid.New(), StartDate: "03-2024"},
		{ServiceName: "Test Yandex", Price: 101, UserId: userId, StartDate: "01-2024"},
		{ServiceName: "Test Google", Price: 202, UserId: userId, StartDate: "02-2024"},
		{ServiceName: "Test Yahoo", Price: 303, UserId: userId, StartDate: "03-2024"},
		{ServiceName: "Test Yahoo", Price: 404, UserId: userId, StartDate: "04-2024"},
		{ServiceName: "Test Yahoo", Price: 505, UserId: userId, StartDate: "05-2024"},
	}

	for _, sub := range subs {
		suite.testDB.Create(&sub)
	}

	resp, err := suite.makeRequest(
		"GET",
		fmt.Sprintf(
			"/api/v1/subscriptions/calculate?user_id=%s&start_date=%s&end_date=%s",
			userId.String(), "02-2024", "04-2025",
		),
		nil,
	)
	assert.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	assert.NotEmpty(suite.T(), resp.Body, "Response body should not be empty")

	body, err := io.ReadAll(resp.Body)
	assert.NoError(suite.T(), err)

	var result map[string]int
	err = json.Unmarshal(body, &result)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 909, result["total"])

}
