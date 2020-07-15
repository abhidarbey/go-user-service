package seed

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/abhidarbey/user-service/pkg/models"
	"github.com/abhidarbey/user-service/pkg/settings/db/mongodb"
	"github.com/abhidarbey/user-service/pkg/utils/resterr"
	"go.mongodb.org/mongo-driver/bson"
)

// GetOrCreate permission
func GetOrCreate(permission models.Permission) (*models.Permission, *resterr.RestErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userDB := mongodb.Client.Database("erp-user-service")

	permissionCollection := userDB.Collection("permission")

	// GetPermission
	filter := bson.M{"name": permission.Name}
	err := permissionCollection.FindOne(ctx, filter).Decode(&permission)
	if err == nil {
		return nil, nil
	}

	// Create permission
	_, err = permissionCollection.InsertOne(ctx, bson.M{
		"name": permission.Name,
	})
	if err != nil {
		return nil, resterr.NewInternalServerError(err.Error())
	}
	return &permission, nil
}

// AddPermissions to the db
func AddPermissions() {
	// SeedPermissions struct
	type SeedPermissions struct {
		Permissions []models.Permission
	}
	// Open our orgJSON
	orgJSON, err := os.Open("static/json/permissions/org_permissions.json")
	catalogueJSON, err := os.Open("static/json/permissions/catalogue_permissions.json")
	inventoryJSON, err := os.Open("static/json/permissions/inventory_permissions.json")
	orderJSON, err := os.Open("static/json/permissions/order_permissions.json")
	warehouseJSON, err := os.Open("static/json/permissions/warehouse_permissions.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	// defer the closing of our orgJSON so that we can parse it later on
	defer orgJSON.Close()
	defer catalogueJSON.Close()
	defer inventoryJSON.Close()
	defer orderJSON.Close()
	defer warehouseJSON.Close()

	orgPermissionsData, _ := ioutil.ReadAll(orgJSON)
	cataloguePermissionsData, _ := ioutil.ReadAll(catalogueJSON)
	inventoryPermissionsData, _ := ioutil.ReadAll(inventoryJSON)
	orderPermissionsData, _ := ioutil.ReadAll(orderJSON)
	warehousePermissionsData, _ := ioutil.ReadAll(warehouseJSON)

	var orgPermissions SeedPermissions
	var cataloguePermissions SeedPermissions
	var inventoryPermissions SeedPermissions
	var orderPermissions SeedPermissions
	var warehousePermissions SeedPermissions

	if err := json.Unmarshal(orgPermissionsData, &orgPermissions); err != nil {
		fmt.Println("Error while unmarshalling org permissions json")
	}
	if err := json.Unmarshal(cataloguePermissionsData, &cataloguePermissions); err != nil {
		fmt.Println("Error while unmarshalling catalogue permissions json")
	}
	if err := json.Unmarshal(inventoryPermissionsData, &inventoryPermissions); err != nil {
		fmt.Println("Error while unmarshalling inventory permissions json")
	}
	if err := json.Unmarshal(orderPermissionsData, &orderPermissions); err != nil {
		fmt.Println("Error while unmarshalling order permissions json")
	}
	if err := json.Unmarshal(warehousePermissionsData, &warehousePermissions); err != nil {
		fmt.Println("Error while unmarshalling warehouse permissions json")
	}

	var addCounter int = 0
	var existingCounter int = 0

	// Seed data to db
	for i := 0; i < len(orgPermissions.Permissions); i++ {
		permission := models.Permission{}
		permission.Name = orgPermissions.Permissions[i].Name
		_, err := GetOrCreate(permission)
		if err != nil {
			existingCounter++
		} else {
			addCounter++
		}
	}

	for i := 0; i < len(cataloguePermissions.Permissions); i++ {
		permission := models.Permission{}
		permission.Name = cataloguePermissions.Permissions[i].Name
		_, err := GetOrCreate(permission)
		if err != nil {
			existingCounter++
		} else {
			addCounter++
		}
	}

	for i := 0; i < len(inventoryPermissions.Permissions); i++ {
		permission := models.Permission{}
		permission.Name = inventoryPermissions.Permissions[i].Name
		_, err := GetOrCreate(permission)
		if err != nil {
			existingCounter++
		} else {
			addCounter++
		}
	}

	for i := 0; i < len(orderPermissions.Permissions); i++ {
		permission := models.Permission{}
		permission.Name = orderPermissions.Permissions[i].Name
		_, err := GetOrCreate(permission)
		if err != nil {
			existingCounter++
		} else {
			addCounter++
		}
	}

	for i := 0; i < len(warehousePermissions.Permissions); i++ {
		permission := models.Permission{}
		permission.Name = warehousePermissions.Permissions[i].Name
		_, err := GetOrCreate(permission)
		if err != nil {
			existingCounter++
		} else {
			addCounter++
		}
	}

	fmt.Println("\n//****************************************************//")
	fmt.Printf("%v permissions already exists and %v new permissions added.", existingCounter, addCounter)
	fmt.Println("\n//****************************************************//")
}
