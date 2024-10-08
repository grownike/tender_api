package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"avito_tenders/internal/db"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func setupTestApp() (*gin.Engine, *db.Database) {
	// Загрузим тестовые переменные окружения
	err := godotenv.Load("test.env")
	if err != nil {
		log.Fatalf("Error loading test.env file")
	}

	// Инициализируем базу данных
	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize test database: %v", err)
	}

	// Настраиваем маршруты
	router := SetupRoutes(database)

	return router, database
}

func teardownTestDB(database *db.Database) {
	tables := []string{
		"bid_versions",
		"bid",
		"tender_version",
		"tender",
		"organization_responsible",
		"organization",
		"employee",
	}
	for _, table := range tables {
		database.DB.Exec("TRUNCATE TABLE " + table + " CASCADE")
	}
}

func TestIntegrationCycle(t *testing.T) {
	router, database := setupTestApp()
	defer teardownTestDB(database)

	server := httptest.NewServer(router)
	defer server.Close()

	client := &http.Client{}

	//Создание первого пользователя
	username1 := "user1"
	createUser_ApiRequest(t, client, server.URL, username1, "User1FirstName", "User1LastName")

	//Создание первой компании
	organizationID1 := createOrganization_ApiRequest(t, client, server.URL, "Company A", "Description A", "LLC")

	//Назначение первого пользователя ответственным за первую компанию
	assignResponsible_ApiRequest(t, client, server.URL, organizationID1, username1)

	//Создание тендера от имени первой компании
	tenderID, status := createTender_ApiRequest(t, client, server.URL, organizationID1, username1, "Tender 1", "Description of Tender 1", "Service Type 1")

    if status != 200{
        t.Fatalf("Ожидался статус 200, получен %d", status)
    }
	//Создание второго пользователя
	username2 := "user2"
	createUser_ApiRequest(t, client, server.URL, username2, "User2FirstName", "User2LastName")

	//Создание второй компании
	organizationID2 := createOrganization_ApiRequest(t, client, server.URL, "Company B", "Description B", "JSC")

	//Назначение второго пользователя ответственным за вторую компанию
	assignResponsible_ApiRequest(t, client, server.URL, organizationID2, username2)

	// println(organizationID2)

	//Создание предложения на тендер от ответственного лица второй компании
	createBid_ApiRequest(t, client, server.URL, tenderID, organizationID2, "User", username2, "Bid 1", "Bid 1 Description")

	//Создание предложения на тендер от второй компании напрямую
	createBid_ApiRequest(t, client, server.URL, tenderID, organizationID2, "Organization", username2, "Bid 2", "Bid 2 Description")

}

func TestIntegrationErrorTender(t *testing.T) {
	router, database := setupTestApp()
	defer teardownTestDB(database)

	server := httptest.NewServer(router)
	defer server.Close()

	client := &http.Client{}

	//Создание первого пользователя
	username1 := "user1"
	createUser_ApiRequest(t, client, server.URL, username1, "User1FirstName", "User1LastName")

	//Создание первой компании
	organizationID1 := createOrganization_ApiRequest(t, client, server.URL, "Company A", "Description A", "LLC")

	

	//Создание тендера от имени первой компании, но user1 не может это сделать!
	_, status := createTender_ApiRequest(t, client, server.URL, organizationID1, username1, "Tender 1", "Description of Tender 1", "Service Type 1")

    if status != 403{
        t.Fatalf("Ожидался статус 403, получен %d", status)
    }

    //Назначение первого пользователя ответственным за первую компанию
	assignResponsible_ApiRequest(t, client, server.URL, organizationID1, username1)

    _, status = createTender_ApiRequest(t, client, server.URL, organizationID1, username1, "Tender 1", "Description of Tender 1", "Service Type 1")

    if status != 200{
        t.Fatalf("Ожидался статус 200, получен %d", status)
    }

    _, status = createTender_ApiRequest(t, client, server.URL, organizationID1, "blablabla", "Tender 1", "Description of Tender 1", "Service Type 1")

    if status != 401{
        t.Fatalf("Ожидался статус 401, получен %d", status)
    }
    
    fakeOrgId := "83bbe121-02e7-410e-933f-15d8544ae07b"
    _, status = createTender_ApiRequest(t, client, server.URL, fakeOrgId, username1, "Tender 1", "Description of Tender 1", "Service Type 1")

    if status != 401{
        t.Fatalf("Ожидался статус 401, получен %d", status)
    }
	
}



func createUser_ApiRequest(t *testing.T, client *http.Client, baseURL, username, firstName, lastName string) {
	user := map[string]string{
		"username":  username,
		"firstName": firstName,
		"lastName":  lastName,
	}
	body, _ := json.Marshal(user)
	resp, err := client.Post(baseURL+"/api/employees/new", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Не удалось создать пользователя: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Ожидался статус 200, получен %d", resp.StatusCode)
	}
}

func createOrganization_ApiRequest(t *testing.T, client *http.Client, baseURL, name, description, orgType string) string {
	org := map[string]string{
		"name":        name,
		"description": description,
		"type":        orgType, // Возможные значения: 'IE', 'LLC', 'JSC'
	}
	body, _ := json.Marshal(org)
	resp, err := client.Post(baseURL+"/api/organizations/new", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Не удалось создать компанию: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Ожидался статус 200, получен %d", resp.StatusCode)
	}

	var result struct {
		ID string `json:"id"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	return result.ID
}

func assignResponsible_ApiRequest(t *testing.T, client *http.Client, baseURL, organizationID, username string) {
	assignment := map[string]string{}
	body, _ := json.Marshal(assignment)

	resp, err := client.Post(baseURL+"/api/newAssign/"+organizationID+"/"+username, "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Не удалось назначить ответственного: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Ожидался статус 200, получен %d", resp.StatusCode)
	}
}

func createTender_ApiRequest(t *testing.T, client *http.Client, baseURL, organizationID, creatorUsername, name, description, serviceType string) (string, int) {
	tender := map[string]interface{}{
		"name":            name,
		"description":     description,
		"serviceType":     serviceType,
		"organizationId":  organizationID,
		"creatorUsername": creatorUsername,
	}
	body, _ := json.Marshal(tender)
	resp, err := client.Post(baseURL+"/api/tenders/new", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Не удалось создать тендер: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		ID string `json:"id"`
	}
	json.NewDecoder(resp.Body).Decode(&result)

	return result.ID, resp.StatusCode
}

func createBid_ApiRequest(t *testing.T, client *http.Client, baseURL, tenderID, OrgId, authorType, creatorUsername, name, description string) {
	bid := map[string]interface{}{
		"name":            name,
		"description":     description,
		"tenderId":        tenderID,
		"authorType":      authorType,
		"creatorUsername": creatorUsername,
		"AuthorId":        OrgId,
	}
	body, _ := json.Marshal(bid)
	resp, err := client.Post(baseURL+"/api/bids/new", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Не удалось создать предложение: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Ожидался статус 200, получен %d", resp.StatusCode)
	}
}