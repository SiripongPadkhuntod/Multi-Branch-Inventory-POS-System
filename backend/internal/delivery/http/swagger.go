package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func registerSwagger(r *gin.Engine) {
	r.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/swagger/index.html")
	})
	r.GET("/swagger/index.html", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(swaggerHTML))
	})
	r.GET("/swagger/doc.json", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json; charset=utf-8", []byte(openAPISpec))
	})
}

const swaggerHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>Multi-Branch POS API Docs</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
  <style>
    body { margin: 0; background: #f8fafc; }
    .topbar { display: none; }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = function () {
      window.ui = SwaggerUIBundle({
        url: "/swagger/doc.json",
        dom_id: "#swagger-ui",
        deepLinking: true,
        persistAuthorization: true
      });
    };
  </script>
</body>
</html>`

const openAPISpec = `{
  "openapi": "3.0.3",
  "info": {
    "title": "Multi-Branch Inventory & POS API",
    "version": "1.0.0",
    "description": "Backend API documentation for authentication, products, categories, inventory, transfers, sales, users, branches, and dashboard reports."
  },
  "servers": [
    { "url": "http://localhost:8080/api/v1", "description": "Local backend" }
  ],
  "tags": [
    { "name": "Auth" },
    { "name": "Products" },
    { "name": "Categories" },
    { "name": "Inventory" },
    { "name": "Sales" },
    { "name": "Users" },
    { "name": "Branches" },
    { "name": "Dashboard" }
  ],
  "components": {
    "securitySchemes": {
      "bearerAuth": { "type": "http", "scheme": "bearer", "bearerFormat": "JWT" }
    },
    "schemas": {
      "ApiResponse": {
        "type": "object",
        "properties": {
          "success": { "type": "boolean" },
          "message": { "type": "string" },
          "data": {}
        }
      },
      "ErrorResponse": {
        "type": "object",
        "properties": {
          "success": { "type": "boolean", "example": false },
          "message": { "type": "string" }
        }
      },
      "LoginRequest": {
        "type": "object",
        "required": ["email", "password"],
        "properties": {
          "email": { "type": "string", "format": "email", "example": "owner@example.com" },
          "password": { "type": "string", "example": "password123" }
        }
      },
      "User": {
        "type": "object",
        "properties": {
          "id": { "type": "string", "format": "uuid" },
          "branch_id": { "type": "string", "format": "uuid", "nullable": true },
          "role": { "type": "string", "enum": ["OWNER", "MANAGER", "EMPLOYEE"] },
          "name": { "type": "string" },
          "email": { "type": "string", "format": "email" },
          "status": { "type": "string" },
          "created_at": { "type": "string", "format": "date-time" }
        }
      },
      "Branch": {
        "type": "object",
        "required": ["code", "name"],
        "properties": {
          "id": { "type": "string", "format": "uuid", "readOnly": true },
          "code": { "type": "string", "example": "BR01" },
          "name": { "type": "string", "example": "Branch 01" },
          "address": { "type": "string" },
          "phone": { "type": "string" },
          "status": { "type": "string", "example": "ACTIVE" },
          "created_at": { "type": "string", "format": "date-time", "readOnly": true }
        }
      },
      "Category": {
        "type": "object",
        "required": ["name"],
        "properties": {
          "id": { "type": "string", "format": "uuid", "readOnly": true },
          "name": { "type": "string", "example": "General" },
          "description": { "type": "string" }
        }
      },
      "Product": {
        "type": "object",
        "required": ["sku", "barcode", "name", "category_id"],
        "properties": {
          "id": { "type": "string", "format": "uuid", "readOnly": true },
          "sku": { "type": "string", "example": "SKU-001" },
          "barcode": { "type": "string", "example": "885000000001" },
          "qr_code": { "type": "string" },
          "name": { "type": "string", "example": "Sample Product" },
          "description": { "type": "string" },
          "category_id": { "type": "string", "format": "uuid" },
          "image_url": { "type": "string" },
          "cost_price": { "type": "integer", "format": "int64", "description": "Amount in cents/satang" },
          "sell_price": { "type": "integer", "format": "int64", "description": "Amount in cents/satang" },
          "status": { "type": "string", "example": "ACTIVE" }
        }
      },
      "Inventory": {
        "type": "object",
        "properties": {
          "id": { "type": "string", "format": "uuid" },
          "branch_id": { "type": "string", "format": "uuid" },
          "product_id": { "type": "string", "format": "uuid" },
          "quantity": { "type": "integer", "format": "int64" },
          "reserved_quantity": { "type": "integer", "format": "int64" },
          "updated_at": { "type": "string", "format": "date-time" }
        }
      },
      "InventoryAdjustRequest": {
        "type": "object",
        "required": ["branch_id", "product_id", "quantity_delta"],
        "properties": {
          "branch_id": { "type": "string", "format": "uuid" },
          "product_id": { "type": "string", "format": "uuid" },
          "quantity_delta": { "type": "integer", "format": "int64", "example": 10 },
          "reason": { "type": "string" }
        }
      },
      "InventoryTransferRequest": {
        "type": "object",
        "required": ["from_branch_id", "to_branch_id", "product_id", "quantity"],
        "properties": {
          "from_branch_id": { "type": "string", "format": "uuid" },
          "to_branch_id": { "type": "string", "format": "uuid" },
          "product_id": { "type": "string", "format": "uuid" },
          "quantity": { "type": "integer", "format": "int64", "minimum": 1 }
        }
      },
      "Transfer": {
        "type": "object",
        "properties": {
          "id": { "type": "string", "format": "uuid" },
          "from_branch_id": { "type": "string", "format": "uuid" },
          "from_branch_code": { "type": "string" },
          "from_branch_name": { "type": "string" },
          "to_branch_id": { "type": "string", "format": "uuid" },
          "to_branch_code": { "type": "string" },
          "to_branch_name": { "type": "string" },
          "status": { "type": "string", "enum": ["PENDING", "APPROVED", "REJECTED", "COMPLETED"] },
          "requested_by": { "type": "string", "format": "uuid" },
          "requested_by_name": { "type": "string" },
          "approved_by": { "type": "string", "format": "uuid", "nullable": true },
          "approved_by_name": { "type": "string" },
          "product_id": { "type": "string", "format": "uuid" },
          "product_sku": { "type": "string" },
          "product_barcode": { "type": "string" },
          "product_name": { "type": "string" },
          "quantity": { "type": "integer", "format": "int64" },
          "created_at": { "type": "string", "format": "date-time" },
          "updated_at": { "type": "string", "format": "date-time" }
        }
      },
      "InventoryMovementDetail": {
        "type": "object",
        "properties": {
          "id": { "type": "string", "format": "uuid" },
          "branch_id": { "type": "string", "format": "uuid" },
          "branch_code": { "type": "string" },
          "branch_name": { "type": "string" },
          "product_id": { "type": "string", "format": "uuid" },
          "sku": { "type": "string" },
          "barcode": { "type": "string" },
          "product_name": { "type": "string" },
          "movement_type": { "type": "string", "enum": ["RECEIVE", "SALE", "RETURN", "ADJUSTMENT", "TRANSFER_IN", "TRANSFER_OUT"] },
          "quantity": { "type": "integer", "format": "int64" },
          "reference_id": { "type": "string" },
          "created_by": { "type": "string", "format": "uuid" },
          "created_by_name": { "type": "string" },
          "created_at": { "type": "string", "format": "date-time" }
        }
      },
      "ProductStockSummary": {
        "type": "object",
        "properties": {
          "product_id": { "type": "string", "format": "uuid" },
          "sku": { "type": "string" },
          "barcode": { "type": "string" },
          "product_name": { "type": "string" },
          "total_quantity": { "type": "integer", "format": "int64" },
          "total_reserved": { "type": "integer", "format": "int64" },
          "branches": { "type": "array", "items": { "type": "object" } }
        }
      },
      "CartItemInput": {
        "type": "object",
        "required": ["product_id", "quantity", "final_price"],
        "properties": {
          "product_id": { "type": "string", "format": "uuid" },
          "quantity": { "type": "integer", "format": "int64", "minimum": 1 },
          "final_price": { "type": "integer", "format": "int64" },
          "discount_amount": { "type": "integer", "format": "int64" },
          "discount_reason": { "type": "string" }
        }
      },
      "CreateSaleRequest": {
        "type": "object",
        "required": ["branch_id", "items", "payments"],
        "properties": {
          "branch_id": { "type": "string", "format": "uuid" },
          "discount": { "type": "integer", "format": "int64" },
          "tax": { "type": "integer", "format": "int64" },
          "items": { "type": "array", "items": { "$ref": "#/components/schemas/CartItemInput" } },
          "payments": {
            "type": "array",
            "items": {
              "type": "object",
              "required": ["payment_method", "amount"],
              "properties": {
                "payment_method": { "type": "string", "enum": ["CASH", "PROMPTPAY", "BANK_TRANSFER", "CREDIT_CARD"] },
                "amount": { "type": "integer", "format": "int64" }
              }
            }
          }
        }
      },
      "Sale": {
        "type": "object",
        "properties": {
          "id": { "type": "string", "format": "uuid" },
          "receipt_number": { "type": "string" },
          "branch_id": { "type": "string", "format": "uuid" },
          "employee_id": { "type": "string", "format": "uuid" },
          "subtotal": { "type": "integer", "format": "int64" },
          "discount": { "type": "integer", "format": "int64" },
          "tax": { "type": "integer", "format": "int64" },
          "total": { "type": "integer", "format": "int64" },
          "payment_status": { "type": "string" },
          "created_at": { "type": "string", "format": "date-time" }
        }
      },
      "RefundRequest": {
        "type": "object",
        "required": ["sale_id", "items"],
        "properties": {
          "sale_id": { "type": "string", "format": "uuid" },
          "items": { "type": "array", "items": { "$ref": "#/components/schemas/CartItemInput" } }
        }
      }
    },
    "responses": {
      "Unauthorized": { "description": "Unauthorized", "content": { "application/json": { "schema": { "$ref": "#/components/schemas/ErrorResponse" } } } },
      "Forbidden": { "description": "Forbidden", "content": { "application/json": { "schema": { "$ref": "#/components/schemas/ErrorResponse" } } } },
      "BadRequest": { "description": "Bad request", "content": { "application/json": { "schema": { "$ref": "#/components/schemas/ErrorResponse" } } } }
    }
  },
  "paths": {
    "/auth/login": {
      "post": {
        "tags": ["Auth"],
        "summary": "Login",
        "requestBody": { "required": true, "content": { "application/json": { "schema": { "$ref": "#/components/schemas/LoginRequest" } } } },
        "responses": { "200": { "description": "Login success" }, "401": { "$ref": "#/components/responses/Unauthorized" } }
      }
    },
    "/auth/logout": {
      "post": {
        "tags": ["Auth"],
        "summary": "Logout and clear auth cookies",
        "responses": { "200": { "description": "Logout success" } }
      }
    },
    "/auth/me": {
      "get": {
        "tags": ["Auth"],
        "summary": "Current authenticated user",
        "security": [{ "bearerAuth": [] }],
        "responses": { "200": { "description": "Current user" }, "401": { "$ref": "#/components/responses/Unauthorized" } }
      }
    },
    "/products": {
      "get": {
        "tags": ["Products"],
        "summary": "List products",
        "security": [{ "bearerAuth": [] }],
        "parameters": [
          { "name": "q", "in": "query", "schema": { "type": "string" } },
          { "name": "limit", "in": "query", "schema": { "type": "integer", "default": 50 } },
          { "name": "offset", "in": "query", "schema": { "type": "integer", "default": 0 } }
        ],
        "responses": { "200": { "description": "Products" } }
      },
      "post": {
        "tags": ["Products"],
        "summary": "Create product (OWNER)",
        "security": [{ "bearerAuth": [] }],
        "requestBody": { "required": true, "content": { "application/json": { "schema": { "$ref": "#/components/schemas/Product" } } } },
        "responses": { "201": { "description": "Product created" }, "403": { "$ref": "#/components/responses/Forbidden" } }
      }
    },
    "/products/barcode/{barcode}": {
      "get": {
        "tags": ["Products"],
        "summary": "Find product by barcode",
        "security": [{ "bearerAuth": [] }],
        "parameters": [{ "name": "barcode", "in": "path", "required": true, "schema": { "type": "string" } }],
        "responses": { "200": { "description": "Product" }, "404": { "description": "Not found" } }
      }
    },
    "/products/{id}": {
      "get": {
        "tags": ["Products"],
        "summary": "Find product by ID",
        "security": [{ "bearerAuth": [] }],
        "parameters": [{ "name": "id", "in": "path", "required": true, "schema": { "type": "string", "format": "uuid" } }],
        "responses": { "200": { "description": "Product" } }
      },
      "put": {
        "tags": ["Products"],
        "summary": "Update product (OWNER)",
        "security": [{ "bearerAuth": [] }],
        "parameters": [{ "name": "id", "in": "path", "required": true, "schema": { "type": "string", "format": "uuid" } }],
        "requestBody": { "required": true, "content": { "application/json": { "schema": { "$ref": "#/components/schemas/Product" } } } },
        "responses": { "200": { "description": "Product updated" } }
      },
      "delete": {
        "tags": ["Products"],
        "summary": "Soft delete product (OWNER)",
        "security": [{ "bearerAuth": [] }],
        "parameters": [{ "name": "id", "in": "path", "required": true, "schema": { "type": "string", "format": "uuid" } }],
        "responses": { "200": { "description": "Product deleted" } }
      }
    },
    "/categories": {
      "get": { "tags": ["Categories"], "summary": "List categories", "security": [{ "bearerAuth": [] }], "responses": { "200": { "description": "Categories" } } },
      "post": {
        "tags": ["Categories"],
        "summary": "Create category (OWNER)",
        "security": [{ "bearerAuth": [] }],
        "requestBody": { "required": true, "content": { "application/json": { "schema": { "$ref": "#/components/schemas/Category" } } } },
        "responses": { "201": { "description": "Category created" } }
      }
    },
    "/categories/{id}": {
      "put": {
        "tags": ["Categories"],
        "summary": "Update category (OWNER)",
        "security": [{ "bearerAuth": [] }],
        "parameters": [{ "name": "id", "in": "path", "required": true, "schema": { "type": "string", "format": "uuid" } }],
        "requestBody": { "required": true, "content": { "application/json": { "schema": { "$ref": "#/components/schemas/Category" } } } },
        "responses": { "200": { "description": "Category updated" } }
      },
      "delete": {
        "tags": ["Categories"],
        "summary": "Soft delete category (OWNER)",
        "security": [{ "bearerAuth": [] }],
        "parameters": [{ "name": "id", "in": "path", "required": true, "schema": { "type": "string", "format": "uuid" } }],
        "responses": { "200": { "description": "Category deleted" } }
      }
    },
    "/inventories": {
      "get": {
        "tags": ["Inventory"],
        "summary": "List inventories",
        "security": [{ "bearerAuth": [] }],
        "parameters": [
          { "name": "q", "in": "query", "schema": { "type": "string" } },
          { "name": "branch_id", "in": "query", "schema": { "type": "string", "format": "uuid" } }
        ],
        "responses": { "200": { "description": "Inventories" } }
      }
    },
    "/inventories/all-stock": {
      "get": {
        "tags": ["Inventory"],
        "summary": "All product stock across all branches (OWNER)",
        "security": [{ "bearerAuth": [] }],
        "parameters": [{ "name": "q", "in": "query", "schema": { "type": "string" } }],
        "responses": { "200": { "description": "Stock summary", "content": { "application/json": { "schema": { "$ref": "#/components/schemas/ApiResponse" } } } } }
      }
    },
    "/inventories/movements": {
      "get": {
        "tags": ["Inventory"],
        "summary": "Inventory movement history (OWNER)",
        "security": [{ "bearerAuth": [] }],
        "parameters": [
          { "name": "q", "in": "query", "schema": { "type": "string" } },
          { "name": "branch_id", "in": "query", "schema": { "type": "string", "format": "uuid" } },
          { "name": "limit", "in": "query", "schema": { "type": "integer", "default": 150 } }
        ],
        "responses": { "200": { "description": "Movement history" } }
      }
    },
    "/inventories/adjust": {
      "post": {
        "tags": ["Inventory"],
        "summary": "Adjust stock (OWNER)",
        "security": [{ "bearerAuth": [] }],
        "requestBody": { "required": true, "content": { "application/json": { "schema": { "$ref": "#/components/schemas/InventoryAdjustRequest" } } } },
        "responses": { "200": { "description": "Stock adjusted" } }
      }
    },
    "/inventories/receive": {
      "post": {
        "tags": ["Inventory"],
        "summary": "Receive stock (OWNER or MANAGER)",
        "security": [{ "bearerAuth": [] }],
        "requestBody": { "required": true, "content": { "application/json": { "schema": { "$ref": "#/components/schemas/InventoryAdjustRequest" } } } },
        "responses": { "200": { "description": "Stock received" } }
      }
    },
    "/inventories/transfers": {
      "get": {
        "tags": ["Inventory"],
        "summary": "List stock transfer requests (OWNER)",
        "security": [{ "bearerAuth": [] }],
        "parameters": [
          { "name": "status", "in": "query", "schema": { "type": "string", "enum": ["PENDING", "APPROVED", "REJECTED", "COMPLETED"] } },
          { "name": "limit", "in": "query", "schema": { "type": "integer", "default": 150 } }
        ],
        "responses": { "200": { "description": "Transfer requests" } }
      }
    },
    "/inventories/transfer": {
      "post": {
        "tags": ["Inventory"],
        "summary": "Request a branch stock transfer (OWNER)",
        "security": [{ "bearerAuth": [] }],
        "requestBody": { "required": true, "content": { "application/json": { "schema": { "$ref": "#/components/schemas/InventoryTransferRequest" } } } },
        "responses": { "201": { "description": "Transfer requested" } }
      }
    },
    "/inventories/transfers/{id}/approve": {
      "post": {
        "tags": ["Inventory"],
        "summary": "Approve a pending transfer (OWNER)",
        "security": [{ "bearerAuth": [] }],
        "parameters": [{ "name": "id", "in": "path", "required": true, "schema": { "type": "string", "format": "uuid" } }],
        "responses": { "200": { "description": "Transfer approved" } }
      }
    },
    "/inventories/transfers/{id}/reject": {
      "post": {
        "tags": ["Inventory"],
        "summary": "Reject a pending transfer (OWNER)",
        "security": [{ "bearerAuth": [] }],
        "parameters": [{ "name": "id", "in": "path", "required": true, "schema": { "type": "string", "format": "uuid" } }],
        "responses": { "200": { "description": "Transfer rejected" } }
      }
    },
    "/inventories/transfers/{id}/complete": {
      "post": {
        "tags": ["Inventory"],
        "summary": "Complete an approved transfer and move stock (OWNER)",
        "security": [{ "bearerAuth": [] }],
        "parameters": [{ "name": "id", "in": "path", "required": true, "schema": { "type": "string", "format": "uuid" } }],
        "responses": { "200": { "description": "Transfer completed" } }
      }
    },
    "/sales": {
      "get": {
        "tags": ["Sales"],
        "summary": "List sales visible to current user",
        "security": [{ "bearerAuth": [] }],
        "parameters": [{ "name": "branch_id", "in": "query", "schema": { "type": "string", "format": "uuid" } }],
        "responses": { "200": { "description": "Sales" } }
      },
      "post": {
        "tags": ["Sales"],
        "summary": "Create POS sale",
        "security": [{ "bearerAuth": [] }],
        "requestBody": { "required": true, "content": { "application/json": { "schema": { "$ref": "#/components/schemas/CreateSaleRequest" } } } },
        "responses": { "201": { "description": "Sale created" } }
      }
    },
    "/sales/branch": {
      "get": {
        "tags": ["Sales"],
        "summary": "List branch sales (OWNER or MANAGER)",
        "security": [{ "bearerAuth": [] }],
        "parameters": [{ "name": "branch_id", "in": "query", "schema": { "type": "string", "format": "uuid" } }],
        "responses": { "200": { "description": "Branch sales" } }
      }
    },
    "/sales/{id}": {
      "get": {
        "tags": ["Sales"],
        "summary": "Sale detail",
        "security": [{ "bearerAuth": [] }],
        "parameters": [{ "name": "id", "in": "path", "required": true, "schema": { "type": "string", "format": "uuid" } }],
        "responses": { "200": { "description": "Sale detail" } }
      }
    },
    "/sales/refund": {
      "post": {
        "tags": ["Sales"],
        "summary": "Refund sale items",
        "security": [{ "bearerAuth": [] }],
        "requestBody": { "required": true, "content": { "application/json": { "schema": { "$ref": "#/components/schemas/RefundRequest" } } } },
        "responses": { "200": { "description": "Refund completed" } }
      }
    },
    "/users": {
      "get": { "tags": ["Users"], "summary": "List users (OWNER or MANAGER scoped)", "security": [{ "bearerAuth": [] }], "responses": { "200": { "description": "Users" } } },
      "post": {
        "tags": ["Users"],
        "summary": "Create user",
        "security": [{ "bearerAuth": [] }],
        "requestBody": { "required": true, "content": { "application/json": { "schema": { "allOf": [{ "$ref": "#/components/schemas/User" }, { "type": "object", "properties": { "password": { "type": "string", "example": "password123" } } }] } } } },
        "responses": { "201": { "description": "User created" } }
      }
    },
    "/users/{id}": {
      "put": {
        "tags": ["Users"],
        "summary": "Update user",
        "security": [{ "bearerAuth": [] }],
        "parameters": [{ "name": "id", "in": "path", "required": true, "schema": { "type": "string", "format": "uuid" } }],
        "requestBody": { "required": true, "content": { "application/json": { "schema": { "$ref": "#/components/schemas/User" } } } },
        "responses": { "200": { "description": "User updated" } }
      }
    },
    "/users/sales-summary": {
      "get": { "tags": ["Users"], "summary": "Sales by employee", "security": [{ "bearerAuth": [] }], "responses": { "200": { "description": "Employee sales summary" } } }
    },
    "/branches/my": {
      "get": { "tags": ["Branches"], "summary": "List accessible branches", "security": [{ "bearerAuth": [] }], "responses": { "200": { "description": "Branches" } } }
    },
    "/branches": {
      "post": {
        "tags": ["Branches"],
        "summary": "Create branch (OWNER)",
        "security": [{ "bearerAuth": [] }],
        "requestBody": { "required": true, "content": { "application/json": { "schema": { "$ref": "#/components/schemas/Branch" } } } },
        "responses": { "201": { "description": "Branch created" } }
      }
    },
    "/branches/{id}": {
      "put": {
        "tags": ["Branches"],
        "summary": "Update branch (OWNER)",
        "security": [{ "bearerAuth": [] }],
        "parameters": [{ "name": "id", "in": "path", "required": true, "schema": { "type": "string", "format": "uuid" } }],
        "requestBody": { "required": true, "content": { "application/json": { "schema": { "$ref": "#/components/schemas/Branch" } } } },
        "responses": { "200": { "description": "Branch updated" } }
      }
    },
    "/dashboard/summary": {
      "get": {
        "tags": ["Dashboard"],
        "summary": "Dashboard summary (OWNER or MANAGER)",
        "security": [{ "bearerAuth": [] }],
        "parameters": [{ "name": "branch_id", "in": "query", "schema": { "type": "string", "format": "uuid" } }],
        "responses": { "200": { "description": "Dashboard summary" } }
      }
    }
  }
}`
