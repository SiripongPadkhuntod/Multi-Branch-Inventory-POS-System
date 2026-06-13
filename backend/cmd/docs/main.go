package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	spec := openAPISpec()
	raw, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		panic(err)
	}
	raw = append(raw, '\n')

	docsDir := filepath.Join("docs")
	mustWrite(filepath.Join(docsDir, "swagger.json"), raw)
	// JSON is valid YAML 1.2, so this keeps the YAML artifact dependency-free.
	mustWrite(filepath.Join(docsDir, "swagger.yaml"), raw)
	mustWrite(filepath.Join(docsDir, "openapi.go"), []byte(fmt.Sprintf(`package docs

const OpenAPISpec = %q
`, string(raw))))
}

func mustWrite(path string, data []byte) {
	if err := os.WriteFile(path, data, 0o644); err != nil {
		panic(err)
	}
}

func openAPISpec() map[string]any {
	return map[string]any{
		"openapi": "3.0.3",
		"info": map[string]any{
			"title":       "Multi-Branch Inventory & POS API",
			"version":     "1.0.0",
			"description": "Authentication, product, inventory, sales, branch, employee, dashboard, and template-style hexagonal endpoints.",
		},
		"servers": []map[string]string{
			{"url": "http://localhost:8080", "description": "Local Docker / Nginx"},
			{"url": "http://localhost:3000", "description": "Local frontend proxy"},
		},
		"tags": []map[string]string{
			{"name": "Auth"},
			{"name": "Products"},
			{"name": "Categories"},
			{"name": "Inventory"},
			{"name": "Sales"},
			{"name": "Users"},
			{"name": "Branches"},
			{"name": "Dashboard"},
			{"name": "Audit"},
			{"name": "POS Service"},
		},
		"components": components(),
		"paths":      paths(),
	}
}

func components() map[string]any {
	return map[string]any{
		"securitySchemes": map[string]any{
			"bearerAuth": map[string]string{"type": "http", "scheme": "bearer", "bearerFormat": "JWT"},
		},
		"schemas": map[string]any{
			"ApiResponse": schemaObject(map[string]any{
				"success": map[string]string{"type": "boolean"},
				"message": map[string]string{"type": "string"},
				"data":    map[string]any{},
			}),
			"TemplateSuccessResponse": schemaObject(map[string]any{
				"code":        map[string]string{"type": "string", "example": "SUCCESS"},
				"description": map[string]string{"type": "string", "example": "success"},
			}),
			"ErrorResponse": schemaObject(map[string]any{
				"success": map[string]any{"type": "boolean", "example": false},
				"message": map[string]string{"type": "string"},
			}),
			"TemplateErrorResponse": schemaObject(map[string]any{
				"code":        map[string]any{"type": "string", "example": "ERROR"},
				"description": map[string]string{"type": "string"},
			}),
			"LoginRequest": schemaObject(map[string]any{
				"email":    map[string]any{"type": "string", "format": "email", "example": "owner@example.com"},
				"password": map[string]any{"type": "string", "example": "password123"},
			}, "email", "password"),
			"LoginResponse": schemaObject(map[string]any{
				"access_token":  map[string]string{"type": "string"},
				"refresh_token": map[string]string{"type": "string"},
				"user":          ref("User"),
			}),
			"User": schemaObject(map[string]any{
				"id":         uuidProp(),
				"branch_id":  nullableUUIDProp(),
				"branch_ids": map[string]any{"type": "array", "items": uuidProp()},
				"role":       enumProp("OWNER", "MANAGER", "EMPLOYEE"),
				"name":       map[string]string{"type": "string"},
				"email":      map[string]any{"type": "string", "format": "email"},
				"status":     map[string]string{"type": "string"},
				"created_at": dateTimeProp(),
			}),
			"Branch": schemaObject(map[string]any{
				"id":         uuidProp(),
				"code":       map[string]any{"type": "string", "example": "BR01"},
				"name":       map[string]any{"type": "string", "example": "Siam Branch"},
				"address":    map[string]string{"type": "string"},
				"phone":      map[string]string{"type": "string"},
				"status":     map[string]any{"type": "string", "example": "ACTIVE"},
				"created_at": dateTimeProp(),
			}, "code", "name"),
			"Category": schemaObject(map[string]any{
				"id":          uuidProp(),
				"name":        map[string]any{"type": "string", "example": "Beverages"},
				"description": map[string]string{"type": "string"},
			}, "name"),
			"Product": schemaObject(map[string]any{
				"id":          uuidProp(),
				"sku":         map[string]any{"type": "string", "example": "SKU-001"},
				"barcode":     map[string]any{"type": "string", "example": "885000000001"},
				"qr_code":     map[string]string{"type": "string"},
				"name":        map[string]any{"type": "string", "example": "Arabica Coffee 250g"},
				"description": map[string]string{"type": "string"},
				"category_id": uuidProp(),
				"image_url":   map[string]any{"type": "string", "format": "uri"},
				"cost_price":  int64Prop(),
				"sell_price":  int64Prop(),
				"status":      map[string]any{"type": "string", "example": "ACTIVE"},
			}, "sku", "barcode", "name", "category_id"),
			"Inventory": schemaObject(map[string]any{
				"id":                uuidProp(),
				"branch_id":         uuidProp(),
				"product_id":        uuidProp(),
				"quantity":          int64Prop(),
				"reserved_quantity": int64Prop(),
				"reorder_threshold": int64Prop(),
				"updated_at":        dateTimeProp(),
			}),
			"InventoryMovement": schemaObject(map[string]any{
				"id":              uuidProp(),
				"branch_id":       uuidProp(),
				"branch_code":     map[string]string{"type": "string"},
				"branch_name":     map[string]string{"type": "string"},
				"product_id":      uuidProp(),
				"sku":             map[string]string{"type": "string"},
				"barcode":         map[string]string{"type": "string"},
				"product_name":    map[string]string{"type": "string"},
				"movement_type":   enumProp("RECEIVE", "SALE", "RETURN", "ADJUSTMENT", "TRANSFER_IN", "TRANSFER_OUT"),
				"quantity":        int64Prop(),
				"reference_id":    map[string]string{"type": "string"},
				"created_by":      uuidProp(),
				"created_by_name": map[string]string{"type": "string"},
				"created_at":      dateTimeProp(),
			}),
			"InventoryAdjustRequest": schemaObject(map[string]any{
				"branch_id":      uuidProp(),
				"product_id":     uuidProp(),
				"quantity_delta": int64Prop(),
				"reason":         map[string]string{"type": "string"},
			}, "branch_id", "product_id", "quantity_delta"),
			"TransferRequest": schemaObject(map[string]any{
				"from_branch_id": uuidProp(),
				"to_branch_id":   uuidProp(),
				"product_id":     uuidProp(),
				"quantity":       int64Prop(),
			}, "from_branch_id", "to_branch_id", "product_id", "quantity"),
			"BranchStockDetail": schemaObject(map[string]any{
				"branch_id":         uuidProp(),
				"branch_code":       map[string]string{"type": "string"},
				"branch_name":       map[string]string{"type": "string"},
				"quantity":          int64Prop(),
				"reserved_quantity": int64Prop(),
				"reorder_threshold": int64Prop(),
				"updated_at":        dateTimeProp(),
			}),
			"ProductStockSummary": schemaObject(map[string]any{
				"product_id":     uuidProp(),
				"sku":            map[string]string{"type": "string"},
				"barcode":        map[string]string{"type": "string"},
				"image_url":      map[string]any{"type": "string", "format": "uri"},
				"product_name":   map[string]string{"type": "string"},
				"total_quantity": int64Prop(),
				"total_reserved": int64Prop(),
				"branches":       arrayRef("BranchStockDetail"),
			}),
			"ProductListResponse": schemaObject(map[string]any{
				"code":        map[string]any{"type": "string", "example": "SUCCESS"},
				"description": map[string]any{"type": "string", "example": "success"},
				"data":        arrayRef("Product"),
			}),
			"ProductResponse": schemaObject(map[string]any{
				"code":        map[string]any{"type": "string", "example": "SUCCESS"},
				"description": map[string]any{"type": "string", "example": "success"},
				"data":        ref("Product"),
			}),
			"InventoryListResponse": schemaObject(map[string]any{
				"code":        map[string]any{"type": "string", "example": "SUCCESS"},
				"description": map[string]any{"type": "string", "example": "success"},
				"data":        arrayRef("Inventory"),
			}),
			"InventoryMovementsResponse": schemaObject(map[string]any{
				"code":        map[string]any{"type": "string", "example": "SUCCESS"},
				"description": map[string]any{"type": "string", "example": "success"},
				"data":        arrayRef("InventoryMovement"),
			}),
			"InventoryAllStockResponse": schemaObject(map[string]any{
				"code":        map[string]any{"type": "string", "example": "SUCCESS"},
				"description": map[string]any{"type": "string", "example": "success"},
				"data":        arrayRef("ProductStockSummary"),
			}),
			"CartItemInput": schemaObject(map[string]any{
				"product_id":      uuidProp(),
				"quantity":        int64Prop(),
				"final_price":     int64Prop(),
				"discount_amount": int64Prop(),
				"discount_reason": map[string]string{"type": "string"},
			}, "product_id", "quantity", "final_price"),
			"PaymentInput": schemaObject(map[string]any{
				"payment_method": enumProp("CASH", "PROMPTPAY", "BANK_TRANSFER", "CREDIT_CARD"),
				"amount":         int64Prop(),
			}, "payment_method", "amount"),
			"CreateSaleRequest": schemaObject(map[string]any{
				"branch_id": uuidProp(),
				"items":     arrayRef("CartItemInput"),
				"payments":  arrayRef("PaymentInput"),
				"discount":  int64Prop(),
				"tax":       int64Prop(),
			}, "branch_id", "items", "payments"),
			"Sale": schemaObject(map[string]any{
				"id":             uuidProp(),
				"receipt_number": map[string]string{"type": "string"},
				"branch_id":      uuidProp(),
				"employee_id":    uuidProp(),
				"subtotal":       int64Prop(),
				"discount":       int64Prop(),
				"tax":            int64Prop(),
				"total":          int64Prop(),
				"payment_status": map[string]string{"type": "string"},
				"refund_status":  map[string]string{"type": "string"},
				"created_at":     dateTimeProp(),
			}),
			"DashboardSummary": schemaObject(map[string]any{
				"daily_sales":       int64Prop(),
				"monthly_sales":     int64Prop(),
				"revenue":           int64Prop(),
				"profit":            int64Prop(),
				"low_stock":         map[string]any{"type": "array", "items": map[string]string{"type": "object"}},
				"top_products":      map[string]any{"type": "array", "items": map[string]string{"type": "object"}},
				"branch_comparison": map[string]any{"type": "array", "items": map[string]string{"type": "object"}},
			}),
			"AuditLog": schemaObject(map[string]any{
				"id":          uuidProp(),
				"user_id":     uuidProp(),
				"user_name":   map[string]string{"type": "string"},
				"user_email":  map[string]string{"type": "string"},
				"action":      map[string]string{"type": "string"},
				"entity_type": map[string]string{"type": "string"},
				"entity_id":   uuidProp(),
				"old_data":    map[string]string{"type": "string"},
				"new_data":    map[string]string{"type": "string"},
				"ip_address":  map[string]string{"type": "string"},
				"created_at":  dateTimeProp(),
			}),
		},
	}
}

func paths() map[string]any {
	return map[string]any{
		"/healthz":                           get("Dashboard", "Legacy health check", nil, "ApiResponse", false),
		"/swagger/doc.json":                  get("Dashboard", "OpenAPI JSON document", nil, "", false),
		"/api/v1/auth/login":                 post("Auth", "Login", "LoginRequest", "ApiResponse", false),
		"/api/v1/auth/refresh":               post("Auth", "Refresh token", "", "ApiResponse", false),
		"/api/v1/auth/logout":                post("Auth", "Logout", "", "ApiResponse", true),
		"/api/v1/auth/me":                    get("Auth", "Current user", nil, "ApiResponse", true),
		"/api/v1/products":                   pathItem(get("Products", "List products", productSearchParams(), "ApiResponse", true), post("Products", "Create product", "Product", "ApiResponse", true)),
		"/api/v1/products/{id}":              pathItem(get("Products", "Get product", idParam(), "ApiResponse", true), put("Products", "Update product", "Product", "ApiResponse", true), deleteOp("Products", "Delete product", true)),
		"/api/v1/products/{id}/image":        postMultipart("Products", "Upload product image", true),
		"/api/v1/products/barcode/{barcode}": get("Products", "Find product by barcode", []map[string]any{pathParam("barcode", "string")}, "ApiResponse", true),
		"/api/v1/categories":                 pathItem(get("Categories", "List categories", nil, "ApiResponse", true), post("Categories", "Create category", "Category", "ApiResponse", true)),
		"/api/v1/categories/{id}":            pathItem(put("Categories", "Update category", "Category", "ApiResponse", true), deleteOp("Categories", "Delete category", true)),
		"/api/v1/inventories":                get("Inventory", "List inventories", inventoryQueryParams(), "ApiResponse", true),
		"/api/v1/inventories/all-stock":      get("Inventory", "All stock across branches", queryParamOnly(), "ApiResponse", true),
		"/api/v1/inventories/movements":      get("Inventory", "Inventory movements", movementQueryParams(), "ApiResponse", true),
		"/api/v1/inventories/adjust":         post("Inventory", "Adjust stock", "InventoryAdjustRequest", "ApiResponse", true),
		"/api/v1/inventories/receive":        post("Inventory", "Receive stock", "InventoryAdjustRequest", "ApiResponse", true),
		"/api/v1/inventories/transfer":       post("Inventory", "Create transfer", "TransferRequest", "ApiResponse", true),
		"/api/v1/inventories/transfers":      get("Inventory", "List transfers", []map[string]any{queryParam("status", "string"), queryParam("limit", "integer")}, "ApiResponse", true),
		"/api/v1/sales":                      pathItem(get("Sales", "List sales", salesQueryParams(), "ApiResponse", true), post("Sales", "Create sale", "CreateSaleRequest", "ApiResponse", true)),
		"/api/v1/sales/{id}":                 get("Sales", "Sale detail", idParam(), "ApiResponse", true),
		"/api/v1/sales/refund":               post("Sales", "Refund / return sale", "", "ApiResponse", true),
		"/api/v1/users":                      pathItem(get("Users", "List users", nil, "ApiResponse", true), post("Users", "Create user", "", "ApiResponse", true)),
		"/api/v1/users/{id}":                 put("Users", "Update user", "", "ApiResponse", true),
		"/api/v1/users/sales-summary":        get("Users", "Employee sales summary", nil, "ApiResponse", true),
		"/api/v1/branches":                   pathItem(get("Branches", "List branches", nil, "ApiResponse", true), post("Branches", "Create branch", "Branch", "ApiResponse", true)),
		"/api/v1/branches/{id}":              put("Branches", "Update branch", "Branch", "ApiResponse", true),
		"/api/v1/dashboard":                  get("Dashboard", "Dashboard summary", []map[string]any{queryParam("branch_id", "string")}, "ApiResponse", true),
		"/api/v1/dashboard/branches":         get("Dashboard", "Accessible branches", nil, "ApiResponse", true),
		"/api/v1/audit-logs":                 get("Audit", "Audit logs", []map[string]any{queryParam("action", "string"), queryParam("entity_type", "string"), queryParam("q", "string"), queryParam("limit", "integer")}, "ApiResponse", true),
		"/pos-service/v1/health":             get("POS Service", "Template service health check", nil, "TemplateSuccessResponse", false),
		"/pos-service/v1/user/login":         post("POS Service", "Template login", "LoginRequest", "LoginResponse", false),
		"/pos-service/v1/user/refresh":       post("POS Service", "Template refresh", "", "LoginResponse", false),
		"/pos-service/v1/user/logout":        post("POS Service", "Template logout", "", "TemplateSuccessResponse", false),
		"/pos-service/v1/products":           get("POS Service", "Template list products", productSearchParams(), "ProductListResponse", true),
		"/pos-service/v1/products/barcode/{barcode}": get(
			"POS Service", "Template find product by barcode", []map[string]any{pathParam("barcode", "string")}, "ProductResponse", true,
		),
		"/pos-service/v1/inventories":           get("POS Service", "Template list inventories", inventoryQueryParams(), "InventoryListResponse", true),
		"/pos-service/v1/inventories/all-stock": get("POS Service", "Template all stock", queryParamOnly(), "InventoryAllStockResponse", true),
		"/pos-service/v1/inventories/movements": get("POS Service", "Template inventory movements", movementQueryParams(), "InventoryMovementsResponse", true),
		"/pos-service/v1/inventories/adjust":    post("POS Service", "Template adjust stock", "InventoryAdjustRequest", "TemplateSuccessResponse", true),
		"/pos-service/v1/inventories/receive":   post("POS Service", "Template receive stock", "InventoryAdjustRequest", "TemplateSuccessResponse", true),
	}
}

func schemaObject(properties map[string]any, required ...string) map[string]any {
	out := map[string]any{"type": "object", "properties": properties}
	if len(required) > 0 {
		out["required"] = required
	}
	return out
}

func ref(name string) map[string]string {
	return map[string]string{"$ref": "#/components/schemas/" + name}
}

func arrayRef(name string) map[string]any {
	return map[string]any{"type": "array", "items": ref(name)}
}

func uuidProp() map[string]string {
	return map[string]string{"type": "string", "format": "uuid"}
}

func nullableUUIDProp() map[string]any {
	return map[string]any{"type": "string", "format": "uuid", "nullable": true}
}

func dateTimeProp() map[string]string {
	return map[string]string{"type": "string", "format": "date-time"}
}

func int64Prop() map[string]string {
	return map[string]string{"type": "integer", "format": "int64"}
}

func enumProp(values ...string) map[string]any {
	return map[string]any{"type": "string", "enum": values}
}

func pathItem(ops ...map[string]any) map[string]any {
	item := map[string]any{}
	for _, op := range ops {
		for method, value := range op {
			item[method] = value
		}
	}
	return item
}

func get(tag, summary string, params []map[string]any, responseSchema string, secure bool) map[string]any {
	return map[string]any{"get": operation(tag, summary, "", responseSchema, params, secure)}
}

func post(tag, summary, requestSchema, responseSchema string, secure bool) map[string]any {
	return map[string]any{"post": operation(tag, summary, requestSchema, responseSchema, nil, secure)}
}

func put(tag, summary, requestSchema, responseSchema string, secure bool) map[string]any {
	return map[string]any{"put": operation(tag, summary, requestSchema, responseSchema, nil, secure)}
}

func deleteOp(tag, summary string, secure bool) map[string]any {
	return map[string]any{"delete": operation(tag, summary, "", "ApiResponse", idParam(), secure)}
}

func postMultipart(tag, summary string, secure bool) map[string]any {
	op := operation(tag, summary, "", "ApiResponse", idParam(), secure)
	op["requestBody"] = map[string]any{
		"required": true,
		"content": map[string]any{
			"multipart/form-data": map[string]any{
				"schema": schemaObject(map[string]any{"image": map[string]any{"type": "string", "format": "binary"}}, "image"),
			},
		},
	}
	return map[string]any{"post": op}
}

func operation(tag, summary, requestSchema, responseSchema string, params []map[string]any, secure bool) map[string]any {
	op := map[string]any{
		"tags":      []string{tag},
		"summary":   summary,
		"responses": responses(responseSchema),
	}
	if len(params) > 0 {
		op["parameters"] = params
	}
	if requestSchema != "" {
		op["requestBody"] = map[string]any{
			"required": true,
			"content":  map[string]any{"application/json": map[string]any{"schema": ref(requestSchema)}},
		}
	}
	if secure {
		op["security"] = []map[string][]string{{"bearerAuth": []string{}}}
	}
	return op
}

func responses(schema string) map[string]any {
	ok := map[string]any{"description": "Success"}
	if schema != "" {
		ok["content"] = map[string]any{"application/json": map[string]any{"schema": ref(schema)}}
	}
	return map[string]any{
		"200": ok,
		"400": map[string]any{"description": "Bad request"},
		"401": map[string]any{"description": "Unauthorized"},
		"500": map[string]any{"description": "Internal server error"},
	}
}

func idParam() []map[string]any {
	return []map[string]any{pathParam("id", "string")}
}

func pathParam(name, typ string) map[string]any {
	return map[string]any{"name": name, "in": "path", "required": true, "schema": map[string]string{"type": typ}}
}

func queryParam(name, typ string) map[string]any {
	return map[string]any{"name": name, "in": "query", "required": false, "schema": map[string]string{"type": typ}}
}

func queryParamOnly() []map[string]any {
	return []map[string]any{queryParam("q", "string")}
}

func productSearchParams() []map[string]any {
	return []map[string]any{queryParam("q", "string"), queryParam("limit", "integer"), queryParam("offset", "integer")}
}

func inventoryQueryParams() []map[string]any {
	return []map[string]any{queryParam("branch_id", "string"), queryParam("category_id", "string"), queryParam("q", "string")}
}

func movementQueryParams() []map[string]any {
	return []map[string]any{queryParam("branch_id", "string"), queryParam("q", "string"), queryParam("limit", "integer")}
}

func salesQueryParams() []map[string]any {
	return []map[string]any{queryParam("branch_id", "string"), queryParam("date_from", "string"), queryParam("date_to", "string")}
}
