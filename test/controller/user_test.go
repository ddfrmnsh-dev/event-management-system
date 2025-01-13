package controller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"event-management-system/controller"
	"event-management-system/models"
	modelUtil "event-management-system/utils/model_util"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserUseCase untuk menggantikan UserUseCase asli
type MockUserUseCase struct {
	mock.Mock
}

func (m *MockUserUseCase) FindAllUser() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

// Mock untuk FindUserById
func (m *MockUserUseCase) FindUserById(id int) (models.User, error) {
	args := m.Called(id)
	return args.Get(0).(models.User), args.Error(1)
}

// Mock untuk FindUserByUsername
func (m *MockUserUseCase) FindUserByUsername(username string) (models.User, error) {
	args := m.Called(username)
	return args.Get(0).(models.User), args.Error(1)
}

// Mock untuk FindUserByEmail
func (m *MockUserUseCase) FindUserByEmail(email string) (models.User, error) {
	args := m.Called(email)
	return args.Get(0).(models.User), args.Error(1)
}

// Mock untuk CreateUser
func (m *MockUserUseCase) CreateUser(input models.User) (models.User, error) {
	args := m.Called(input)
	return args.Get(0).(models.User), args.Error(1)
}

// Mock untuk UpdateUser
func (m *MockUserUseCase) UpdateUser(inputId models.GetCustomerDetailInput, user models.User) (models.User, error) {
	args := m.Called(inputId, user)
	return args.Get(0).(models.User), args.Error(1)
}

// Mock untuk DeleteUserById
func (m *MockUserUseCase) DeleteUserById(id int) (models.User, error) {
	args := m.Called(id)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserUseCase) FinByParams(params string, value bool) ([]models.User, error) {
	args := m.Called(params, value)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserUseCase) AddRoleToUser(userId int, roleId []int) (models.User, error) {
	args := m.Called(userId, roleId)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserUseCase) FindAllEventUser() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

// MockAuthMiddleware untuk menggantikan AuthMiddleware asli
type MockAuthMiddleware struct {
	mock.Mock
}

// Metode RequireToken yang menggunakan variadic parameter
func (m *MockAuthMiddleware) RequireToken(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Mocked handler yang mensimulasikan otorisasi berhasil
		ctx.Next()
	}
}

func TestCreateUser_Success(t *testing.T) {
	// Inisialisasi Gin router untuk unit test
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	rg := router.Group("/")

	// Mock dependencies
	mockUserUseCase := new(MockUserUseCase)
	mockAuthMiddleware := new(MockAuthMiddleware)
	userController := controller.NewUserController(mockUserUseCase, rg, mockAuthMiddleware)

	// Payload yang dikirimkan oleh client, biasanya dalam request body
	payload := models.User{
		Username: "John Doe",
		Email:    "johndoe@example.com",
		Password: "12345678", // password dalam teks biasa
		// Role:     "Admin",
	}

	// MockUser yang dikembalikan oleh mock userUseCase
	mockUser := models.User{
		Id:       1, // ID ditambahkan oleh sistem
		Username: "John Doe",
		Email:    "johndoe@example.com",
		Password: "hashedpassword123", // password sudah di-hash
		// Role:      "Admin",
		IsActive:  nil,        // misalnya di-set sesuai default dalam logika backend
		CreatedAt: time.Now(), // waktu saat user dibuat
		UpdatedAt: time.Now(), // waktu saat user dibuat pertama kali
	}

	// Setup ekspektasi mock
	mockUserUseCase.On("CreateUser", payload).Return(mockUser, nil)

	// Membuat request dan recorder
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Setup routing dan jalankan request
	userController.Route()
	router.ServeHTTP(w, req)

	// Verifikasi respons
	assert.Equal(t, http.StatusOK, w.Code)
	mockUserUseCase.AssertExpectations(t)

	// Parse response JSON untuk verifikasi data
	var response modelUtil.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err) // Pastikan tidak ada error saat parsing JSON

	// Verifikasi bahwa Data bisa di-unmarshal langsung menjadi models.User
	var actualUser models.User
	if response.Data != nil {
		actualUserBytes, err := json.Marshal(response.Data) // Convert data to JSON
		assert.NoError(t, err)
		err = json.Unmarshal(actualUserBytes, &actualUser) // Convert JSON to models.User
		assert.NoError(t, err)
	}

	// Verifikasi data tanpa membandingkan waktu secara langsung
	expectedUser := mockUser // mockUser yang diharapkan

	// Mengabaikan perbandingan waktu dengan mengatur CreatedAt dan UpdatedAt menjadi zero value
	expectedUser.CreatedAt = time.Time{}
	expectedUser.UpdatedAt = time.Time{}
	actualUser.CreatedAt = time.Time{}
	actualUser.UpdatedAt = time.Time{}

	// Verifikasi data
	assert.True(t, response.Status)
	assert.Equal(t, "Succes create user", response.Message)
	assert.Equal(t, expectedUser, actualUser)

}

func TestCreateUser_BadRequest(t *testing.T) {
	// Inisialisasi router Gin dan dependencies mock
	router := gin.Default()
	rg := router.Group("/")

	mockUserUseCase := new(MockUserUseCase)
	mockAuthMiddleware := new(MockAuthMiddleware)
	userController := controller.NewUserController(mockUserUseCase, rg, mockAuthMiddleware)

	// Buat request yang tidak valid
	invalidPayload := `{ "username": "","email":"", password:"", role:""}`
	req, _ := http.NewRequest("POST", "/users", bytes.NewBufferString(invalidPayload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Jalankan routing dan request
	userController.Route()
	router.ServeHTTP(w, req)

	// Verifikasi respons
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateUser_InternalServerError(t *testing.T) {
	router := gin.Default()
	rg := router.Group("/")

	// Mocking dependencies
	mockUserUseCase := new(MockUserUseCase)
	mockAuthMiddleware := new(MockAuthMiddleware)
	userController := controller.NewUserController(mockUserUseCase, rg, mockAuthMiddleware)

	// Payload yang valid
	payload := models.User{
		Username: "John Doe",
		Email:    "johndoe@example.com",
		Password: "12345678", // password dalam teks biasa
		// Role:     "Admin",
	}

	// Mocking kesalahan saat penyimpanan
	mockUserUseCase.On("CreateUser", payload).Return(models.User{}, errors.New("database error"))

	// Membuat request dan recorder
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Jalankan routing dan request
	userController.Route()
	router.ServeHTTP(w, req)

	// Verifikasi respons
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUserUseCase.AssertExpectations(t)
}
func TestFindUserById_Success(t *testing.T) {
	// Inisialisasi Gin router untuk unit test
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	rg := router.Group("/")

	// Mock dependencies
	mockUserUseCase := new(MockUserUseCase)
	mockAuthMiddleware := new(MockAuthMiddleware)
	userController := controller.NewUserController(mockUserUseCase, rg, mockAuthMiddleware)

	// Data user mock
	userID := 1
	mockUser := models.User{
		Id:       userID,
		Username: "John Doe",
		Email:    "johndoe@example.com",
		// Role:     "Admin",
		IsActive: nil,
	}

	// Setup ekspektasi mock
	mockUserUseCase.On("FindUserById", userID).Return(mockUser, nil)

	// Membuat request dan recorder
	req, _ := http.NewRequest("GET", fmt.Sprintf("/users/%d", userID), nil)
	w := httptest.NewRecorder()

	// Setup routing dan jalankan request
	userController.Route() // Pastikan route sudah di-setup dengan benar
	router.ServeHTTP(w, req)

	// Verifikasi respons
	assert.Equal(t, http.StatusOK, w.Code)
	mockUserUseCase.AssertExpectations(t)

	// Parse response JSON untuk verifikasi data
	var response modelUtil.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verifikasi bahwa data yang diterima sesuai dengan yang diharapkan
	var actualUser models.User
	if response.Data != nil {
		actualUserBytes, err := json.Marshal(response.Data)
		assert.NoError(t, err)
		err = json.Unmarshal(actualUserBytes, &actualUser)
		assert.NoError(t, err)
	}

	// Verifikasi data
	assert.True(t, response.Status)
	assert.Equal(t, "Succes Get User", response.Message)
	assert.Equal(t, mockUser, actualUser)
}

func TestFindUserById_NotFound(t *testing.T) {
	// Inisialisasi Gin router untuk unit test
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	rg := router.Group("/")

	// Mock dependencies
	mockUserUseCase := new(MockUserUseCase)
	mockAuthMiddleware := new(MockAuthMiddleware)
	userController := controller.NewUserController(mockUserUseCase, rg, mockAuthMiddleware)

	// Setup ekspektasi mock untuk kasus ID tidak ditemukan
	userID := 99 // ID yang tidak ada
	mockUserUseCase.On("FindUserById", userID).Return(models.User{}, fmt.Errorf("User not found"))

	// Membuat request dan recorder
	req, _ := http.NewRequest("GET", fmt.Sprintf("/users/%d", userID), nil)
	w := httptest.NewRecorder()

	// Setup routing dan jalankan request
	userController.Route()
	router.ServeHTTP(w, req)

	// Verifikasi respons
	assert.Equal(t, http.StatusNotFound, w.Code)
	mockUserUseCase.AssertExpectations(t)

	// Parse response JSON untuk verifikasi error message
	var response modelUtil.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verifikasi bahwa error message sesuai
	assert.False(t, response.Status)
	assert.Equal(t, "User not found", response.Message)
}
