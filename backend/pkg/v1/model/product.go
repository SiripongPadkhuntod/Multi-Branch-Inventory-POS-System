package model

type Product struct {
	ID          string `json:"id"`
	SKU         string `json:"sku"`
	Barcode     string `json:"barcode"`
	QRCode      string `json:"qr_code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CategoryID  string `json:"category_id"`
	ImageURL    string `json:"image_url"`
	CostPrice   int64  `json:"cost_price"`
	SellPrice   int64  `json:"sell_price"`
	Status      string `json:"status"`
}
