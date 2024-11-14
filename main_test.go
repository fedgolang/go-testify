package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Запрос сформирован корректно, сервис возвращает код ответа 200 и тело ответа не пустое
func TestMainHandlerHappyResult(t *testing.T) {
	expectedCode := 200                                                 // Ожидаемый код
	req := httptest.NewRequest("GET", "/cafe?count=4&city=moscow", nil) // Отправим запрос с параметрами count=4&city=moscow

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.Equal(t, responseRecorder.Code, expectedCode) // Сверим, что код ответа = expectedCode, если !ОК останавливаем тест
	assert.NotEmpty(t, responseRecorder.Body.String())    // Проверим, что тело ответа не пустое
}

// Город, который передаётся в параметре city, не поддерживается. Сервис возвращает код ответа 400 и ошибку wrong city value в теле ответа.
func TestMainHandlerWrongCity(t *testing.T) {
	expectedCode := 400
	expectedBody := "wrong city value"
	req := httptest.NewRequest("GET", "/cafe?count=4&city=moscoww", nil) // Отправим запрос с некорректным городом в параметре city

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.Equal(t, responseRecorder.Code, expectedCode)         // Сверим, что код ответа = expectedCode, если !ОК останавливаем тест
	assert.Equal(t, responseRecorder.Body.String(), expectedBody) // сверим, что в тело вернулось ожидаемое значение expectedBody
}

// Если в параметре count указано больше, чем есть всего, должны вернуться все доступные кафе
func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	expectedCode := 200
	expectedRes := "Мир кофе,Сладкоежка,Кофе и завтраки,Сытый студент"
	totalCount := 4                                                                                                // Кол-во кафешек в moscow
	req := httptest.NewRequest("GET", "/cafe?count="+fmt.Sprintf("%d", totalCount+rand.Int())+"&city=moscow", nil) // Выведем рандомное число больше известного кол-ва, так постараемся избежать парадокса пестицида

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.Equal(t, responseRecorder.Code, expectedCode) // Проверяем, что при большем кол-ве всё равно вернём ОК, если !ОК останавливаем тест

	body := responseRecorder.Body.String() // Запишем тело в строку
	list := strings.Split(body, ",")       // Создадим слайс значений со всеми кафешками

	assert.Equal(t, len(list), totalCount) // проверим, что кол-во вернувшихся кафешек равно ожидаемому значению
	assert.Equal(t, body, expectedRes)     // Для точности сверим, что все кафешки были выведены
}
