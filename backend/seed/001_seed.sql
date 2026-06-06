-- Demo seed data for end-to-end testing.
-- Password for every demo user: password123

INSERT INTO branches (code, name, address, phone, status)
VALUES
  ('HQ','Head Office','Silom Road, Bangkok','020000000','ACTIVE'),
  ('BR01','Siam Branch','Siam Square, Bangkok','020000001','ACTIVE'),
  ('BR02','Chiang Mai Branch','Nimman Road, Chiang Mai','053000002','ACTIVE'),
  ('BR03','Phuket Branch','Patong, Phuket','076000003','ACTIVE')
ON CONFLICT (code) DO UPDATE SET
  name=EXCLUDED.name,
  address=EXCLUDED.address,
  phone=EXCLUDED.phone,
  status=EXCLUDED.status,
  deleted_at=NULL,
  updated_at=now();

INSERT INTO categories (id, name, description)
VALUES
  ('00000000-0000-0000-0000-000000000201','General','Default category'),
  ('00000000-0000-0000-0000-000000000202','Beverages','Drinks, bottled water, tea, and coffee'),
  ('00000000-0000-0000-0000-000000000203','Snacks','Ready-to-sell snacks and sweets'),
  ('00000000-0000-0000-0000-000000000204','Personal Care','Daily personal care products'),
  ('00000000-0000-0000-0000-000000000205','Household','Home and cleaning essentials')
ON CONFLICT (name) DO UPDATE SET
  description=EXCLUDED.description,
  deleted_at=NULL,
  updated_at=now();

INSERT INTO users (id, branch_id, role, name, email, password_hash, status)
VALUES
  ('00000000-0000-0000-0000-000000000301',NULL,'OWNER','Owner','owner@example.com','$2a$10$18RSjOfjiptcqK9EMw9ODOEZ56wu6AN0lO20EuHlVdWEVdYwEzwVy','ACTIVE'),
  ('00000000-0000-0000-0000-000000000303',(SELECT id FROM branches WHERE code='BR01'),'MANAGER','Manager Siam','manager@example.com','$2a$10$18RSjOfjiptcqK9EMw9ODOEZ56wu6AN0lO20EuHlVdWEVdYwEzwVy','ACTIVE'),
  ('00000000-0000-0000-0000-000000000307',(SELECT id FROM branches WHERE code='BR02'),'MANAGER','Manager North','manager-north@example.com','$2a$10$18RSjOfjiptcqK9EMw9ODOEZ56wu6AN0lO20EuHlVdWEVdYwEzwVy','ACTIVE'),
  ('00000000-0000-0000-0000-000000000302',(SELECT id FROM branches WHERE code='BR01'),'EMPLOYEE','Employee Siam','employee@example.com','$2a$10$18RSjOfjiptcqK9EMw9ODOEZ56wu6AN0lO20EuHlVdWEVdYwEzwVy','ACTIVE'),
  ('00000000-0000-0000-0000-000000000305',(SELECT id FROM branches WHERE code='BR02'),'EMPLOYEE','Employee North','employee-north@example.com','$2a$10$18RSjOfjiptcqK9EMw9ODOEZ56wu6AN0lO20EuHlVdWEVdYwEzwVy','ACTIVE'),
  ('00000000-0000-0000-0000-000000000306',(SELECT id FROM branches WHERE code='BR03'),'EMPLOYEE','Employee Phuket','employee-phuket@example.com','$2a$10$18RSjOfjiptcqK9EMw9ODOEZ56wu6AN0lO20EuHlVdWEVdYwEzwVy','ACTIVE')
ON CONFLICT (id) DO UPDATE SET
  branch_id=EXCLUDED.branch_id,
  role=EXCLUDED.role,
  name=EXCLUDED.name,
  email=EXCLUDED.email,
  password_hash=EXCLUDED.password_hash,
  status=EXCLUDED.status,
  deleted_at=NULL,
  updated_at=now();

INSERT INTO user_branches (user_id, branch_id)
VALUES
  ('00000000-0000-0000-0000-000000000303',(SELECT id FROM branches WHERE code='BR01')),
  ('00000000-0000-0000-0000-000000000303',(SELECT id FROM branches WHERE code='HQ')),
  ('00000000-0000-0000-0000-000000000307',(SELECT id FROM branches WHERE code='BR02')),
  ('00000000-0000-0000-0000-000000000307',(SELECT id FROM branches WHERE code='BR03'))
ON CONFLICT (user_id, branch_id) DO UPDATE SET
  deleted_at=NULL,
  updated_at=now();

INSERT INTO products (id, sku, barcode, qr_code, name, description, category_id, image_url, cost_price, sell_price, status)
VALUES
  ('00000000-0000-0000-0000-000000000401','SKU-001','885000000001','QR-885000000001','Arabica Coffee 250g','Roasted coffee beans for POS demo',(SELECT id FROM categories WHERE name='Beverages'),'https://images.unsplash.com/photo-1447933601403-0c6688de566e?auto=format&fit=crop&w=600&q=80',12000,18900,'ACTIVE'),
  ('00000000-0000-0000-0000-000000000402','SKU-002','885000000002','QR-885000000002','Green Tea Bottle','Ready-to-drink green tea',(SELECT id FROM categories WHERE name='Beverages'),'https://images.unsplash.com/photo-1556679343-c7306c1976bc?auto=format&fit=crop&w=600&q=80',1800,2900,'ACTIVE'),
  ('00000000-0000-0000-0000-000000000403','SKU-003','885000000003','QR-885000000003','Mineral Water 600ml','Bottled mineral water',(SELECT id FROM categories WHERE name='Beverages'),'https://images.unsplash.com/photo-1548839140-29a749e1cf4d?auto=format&fit=crop&w=600&q=80',700,1200,'ACTIVE'),
  ('00000000-0000-0000-0000-000000000404','SKU-004','885000000004','QR-885000000004','Chocolate Cookies','Box of chocolate cookies',(SELECT id FROM categories WHERE name='Snacks'),'https://images.unsplash.com/photo-1499636136210-6f4ee915583e?auto=format&fit=crop&w=600&q=80',3200,5900,'ACTIVE'),
  ('00000000-0000-0000-0000-000000000405','SKU-005','885000000005','QR-885000000005','Potato Chips Classic','Classic salted potato chips',(SELECT id FROM categories WHERE name='Snacks'),'https://images.unsplash.com/photo-1566478989037-eec170784d0b?auto=format&fit=crop&w=600&q=80',2100,3900,'ACTIVE'),
  ('00000000-0000-0000-0000-000000000406','SKU-006','885000000006','QR-885000000006','Almond Granola Bar','Healthy snack bar',(SELECT id FROM categories WHERE name='Snacks'),'https://images.unsplash.com/photo-1622484212850-eb596d769edc?auto=format&fit=crop&w=600&q=80',2500,4500,'ACTIVE'),
  ('00000000-0000-0000-0000-000000000407','SKU-007','885000000007','QR-885000000007','Herbal Shampoo','Daily care herbal shampoo',(SELECT id FROM categories WHERE name='Personal Care'),'https://images.unsplash.com/photo-1620916566398-39f1143ab7be?auto=format&fit=crop&w=600&q=80',7800,12900,'ACTIVE'),
  ('00000000-0000-0000-0000-000000000408','SKU-008','885000000008','QR-885000000008','Toothpaste Fresh Mint','Fresh mint toothpaste',(SELECT id FROM categories WHERE name='Personal Care'),'https://images.unsplash.com/photo-1607613009820-a29f7bb81c04?auto=format&fit=crop&w=600&q=80',2400,4900,'ACTIVE'),
  ('00000000-0000-0000-0000-000000000409','SKU-009','885000000009','QR-885000000009','Hand Soap Refill','Liquid hand soap refill pack',(SELECT id FROM categories WHERE name='Household'),'https://images.unsplash.com/photo-1588086390101-53c1a8b1a4f1?auto=format&fit=crop&w=600&q=80',3900,6900,'ACTIVE'),
  ('00000000-0000-0000-0000-000000000410','SKU-010','885000000010','QR-885000000010','Laundry Detergent','Concentrated laundry detergent',(SELECT id FROM categories WHERE name='Household'),'https://images.unsplash.com/photo-1624372635288-475b2fbddff8?auto=format&fit=crop&w=600&q=80',9500,15900,'ACTIVE'),
  ('00000000-0000-0000-0000-000000000411','SKU-011','885000000011','QR-885000000011','Low Stock Test Item','Use this item to test low-stock widgets',(SELECT id FROM categories WHERE name='General'),'https://images.unsplash.com/photo-1516321497487-e288fb19713f?auto=format&fit=crop&w=600&q=80',5000,8900,'ACTIVE'),
  ('00000000-0000-0000-0000-000000000412','SKU-012','885000000012','QR-885000000012','Price Override Demo','Use this item for discount and override testing',(SELECT id FROM categories WHERE name='General'),'https://images.unsplash.com/photo-1581090464777-f3220bbe1b8b?auto=format&fit=crop&w=600&q=80',10000,19900,'ACTIVE')
ON CONFLICT (barcode) DO UPDATE SET
  sku=EXCLUDED.sku,
  barcode=EXCLUDED.barcode,
  qr_code=EXCLUDED.qr_code,
  name=EXCLUDED.name,
  description=EXCLUDED.description,
  category_id=EXCLUDED.category_id,
  image_url=EXCLUDED.image_url,
  cost_price=EXCLUDED.cost_price,
  sell_price=EXCLUDED.sell_price,
  status=EXCLUDED.status,
  deleted_at=NULL,
  updated_at=now();

INSERT INTO inventories (branch_id, product_id, quantity, reserved_quantity, reorder_threshold)
VALUES
  ((SELECT id FROM branches WHERE code='HQ'),(SELECT id FROM products WHERE sku='SKU-001'),60,0,12),
  ((SELECT id FROM branches WHERE code='HQ'),(SELECT id FROM products WHERE sku='SKU-002'),130,0,25),
  ((SELECT id FROM branches WHERE code='HQ'),(SELECT id FROM products WHERE sku='SKU-003'),220,0,30),
  ((SELECT id FROM branches WHERE code='HQ'),(SELECT id FROM products WHERE sku='SKU-004'),70,0,15),
  ((SELECT id FROM branches WHERE code='HQ'),(SELECT id FROM products WHERE sku='SKU-010'),24,0,10),
  ((SELECT id FROM branches WHERE code='BR01'),(SELECT id FROM products WHERE sku='SKU-001'),34,0,10),
  ((SELECT id FROM branches WHERE code='BR01'),(SELECT id FROM products WHERE sku='SKU-002'),95,0,20),
  ((SELECT id FROM branches WHERE code='BR01'),(SELECT id FROM products WHERE sku='SKU-003'),140,0,25),
  ((SELECT id FROM branches WHERE code='BR01'),(SELECT id FROM products WHERE sku='SKU-004'),42,0,12),
  ((SELECT id FROM branches WHERE code='BR01'),(SELECT id FROM products WHERE sku='SKU-005'),85,0,20),
  ((SELECT id FROM branches WHERE code='BR01'),(SELECT id FROM products WHERE sku='SKU-006'),18,0,15),
  ((SELECT id FROM branches WHERE code='BR01'),(SELECT id FROM products WHERE sku='SKU-007'),16,0,8),
  ((SELECT id FROM branches WHERE code='BR01'),(SELECT id FROM products WHERE sku='SKU-008'),50,0,12),
  ((SELECT id FROM branches WHERE code='BR01'),(SELECT id FROM products WHERE sku='SKU-011'),3,0,12),
  ((SELECT id FROM branches WHERE code='BR02'),(SELECT id FROM products WHERE sku='SKU-001'),28,0,10),
  ((SELECT id FROM branches WHERE code='BR02'),(SELECT id FROM products WHERE sku='SKU-002'),82,0,20),
  ((SELECT id FROM branches WHERE code='BR02'),(SELECT id FROM products WHERE sku='SKU-003'),110,0,25),
  ((SELECT id FROM branches WHERE code='BR02'),(SELECT id FROM products WHERE sku='SKU-005'),65,0,20),
  ((SELECT id FROM branches WHERE code='BR02'),(SELECT id FROM products WHERE sku='SKU-007'),9,0,10),
  ((SELECT id FROM branches WHERE code='BR02'),(SELECT id FROM products WHERE sku='SKU-009'),22,0,10),
  ((SELECT id FROM branches WHERE code='BR02'),(SELECT id FROM products WHERE sku='SKU-012'),7,0,10),
  ((SELECT id FROM branches WHERE code='BR03'),(SELECT id FROM products WHERE sku='SKU-002'),45,0,18),
  ((SELECT id FROM branches WHERE code='BR03'),(SELECT id FROM products WHERE sku='SKU-003'),90,0,25),
  ((SELECT id FROM branches WHERE code='BR03'),(SELECT id FROM products WHERE sku='SKU-004'),20,0,12),
  ((SELECT id FROM branches WHERE code='BR03'),(SELECT id FROM products WHERE sku='SKU-006'),8,0,12),
  ((SELECT id FROM branches WHERE code='BR03'),(SELECT id FROM products WHERE sku='SKU-008'),4,0,12),
  ((SELECT id FROM branches WHERE code='BR03'),(SELECT id FROM products WHERE sku='SKU-010'),12,0,8)
ON CONFLICT (branch_id, product_id) DO UPDATE SET
  quantity=EXCLUDED.quantity,
  reserved_quantity=EXCLUDED.reserved_quantity,
  reorder_threshold=EXCLUDED.reorder_threshold,
  deleted_at=NULL,
  updated_at=now();

INSERT INTO sales (id, receipt_number, branch_id, employee_id, subtotal, discount, tax, total, payment_status, created_at, updated_at)
VALUES
  ('00000000-0000-0000-0000-000000000501','RCPT-DEMO-001',(SELECT id FROM branches WHERE code='BR01'),'00000000-0000-0000-0000-000000000302',20100,1000,0,19100,'PAID',now() - interval '6 days',now() - interval '6 days'),
  ('00000000-0000-0000-0000-000000000502','RCPT-DEMO-002',(SELECT id FROM branches WHERE code='BR01'),'00000000-0000-0000-0000-000000000302',17700,0,0,17700,'PAID',now() - interval '5 days',now() - interval '5 days'),
  ('00000000-0000-0000-0000-000000000503','RCPT-DEMO-003',(SELECT id FROM branches WHERE code='BR02'),'00000000-0000-0000-0000-000000000305',22600,1500,0,21100,'PAID',now() - interval '4 days',now() - interval '4 days'),
  ('00000000-0000-0000-0000-000000000504','RCPT-DEMO-004',(SELECT id FROM branches WHERE code='BR03'),'00000000-0000-0000-0000-000000000306',34800,0,0,34800,'PAID',now() - interval '3 days',now() - interval '3 days'),
  ('00000000-0000-0000-0000-000000000505','RCPT-DEMO-005',(SELECT id FROM branches WHERE code='HQ'),'00000000-0000-0000-0000-000000000301',29900,2000,0,27900,'PAID',now() - interval '2 days',now() - interval '2 days'),
  ('00000000-0000-0000-0000-000000000506','RCPT-DEMO-006',(SELECT id FROM branches WHERE code='BR02'),'00000000-0000-0000-0000-000000000305',53100,5000,0,48100,'PAID',now() - interval '1 day',now() - interval '1 day'),
  ('00000000-0000-0000-0000-000000000507','RCPT-DEMO-007',(SELECT id FROM branches WHERE code='BR01'),'00000000-0000-0000-0000-000000000302',24700,0,0,24700,'PAID',now() - interval '8 hours',now() - interval '8 hours'),
  ('00000000-0000-0000-0000-000000000508','RCPT-DEMO-008',(SELECT id FROM branches WHERE code='BR03'),'00000000-0000-0000-0000-000000000306',14700,0,0,14700,'PAID',now() - interval '2 hours',now() - interval '2 hours')
ON CONFLICT (receipt_number) DO UPDATE SET
  branch_id=EXCLUDED.branch_id,
  employee_id=EXCLUDED.employee_id,
  subtotal=EXCLUDED.subtotal,
  discount=EXCLUDED.discount,
  tax=EXCLUDED.tax,
  total=EXCLUDED.total,
  payment_status=EXCLUDED.payment_status,
  created_at=EXCLUDED.created_at,
  updated_at=EXCLUDED.updated_at,
  deleted_at=NULL;

INSERT INTO sale_items (id, sale_id, product_id, quantity, original_price, final_price, discount_amount, discount_reason, employee_id, created_at, updated_at)
VALUES
  ('00000000-0000-0000-0000-000000000601','00000000-0000-0000-0000-000000000501',(SELECT id FROM products WHERE sku='SKU-001'),1,18900,17900,1000,'Loyalty discount','00000000-0000-0000-0000-000000000302',now() - interval '6 days',now() - interval '6 days'),
  ('00000000-0000-0000-0000-000000000602','00000000-0000-0000-0000-000000000501',(SELECT id FROM products WHERE sku='SKU-003'),1,1200,1200,0,'','00000000-0000-0000-0000-000000000302',now() - interval '6 days',now() - interval '6 days'),
  ('00000000-0000-0000-0000-000000000603','00000000-0000-0000-0000-000000000502',(SELECT id FROM products WHERE sku='SKU-004'),3,5900,5900,0,'','00000000-0000-0000-0000-000000000302',now() - interval '5 days',now() - interval '5 days'),
  ('00000000-0000-0000-0000-000000000604','00000000-0000-0000-0000-000000000503',(SELECT id FROM products WHERE sku='SKU-002'),4,2900,2900,0,'','00000000-0000-0000-0000-000000000305',now() - interval '4 days',now() - interval '4 days'),
  ('00000000-0000-0000-0000-000000000605','00000000-0000-0000-0000-000000000503',(SELECT id FROM products WHERE sku='SKU-007'),1,12900,11500,1400,'Manager override','00000000-0000-0000-0000-000000000305',now() - interval '4 days',now() - interval '4 days'),
  ('00000000-0000-0000-0000-000000000606','00000000-0000-0000-0000-000000000504',(SELECT id FROM products WHERE sku='SKU-010'),2,15900,15900,0,'','00000000-0000-0000-0000-000000000306',now() - interval '3 days',now() - interval '3 days'),
  ('00000000-0000-0000-0000-000000000607','00000000-0000-0000-0000-000000000504',(SELECT id FROM products WHERE sku='SKU-005'),2,3900,3900,0,'','00000000-0000-0000-0000-000000000306',now() - interval '3 days',now() - interval '3 days'),
  ('00000000-0000-0000-0000-000000000608','00000000-0000-0000-0000-000000000505',(SELECT id FROM products WHERE sku='SKU-012'),1,19900,17900,2000,'Price override demo','00000000-0000-0000-0000-000000000301',now() - interval '2 days',now() - interval '2 days'),
  ('00000000-0000-0000-0000-000000000609','00000000-0000-0000-0000-000000000505',(SELECT id FROM products WHERE sku='SKU-009'),1,6900,6900,0,'','00000000-0000-0000-0000-000000000301',now() - interval '2 days',now() - interval '2 days'),
  ('00000000-0000-0000-0000-000000000610','00000000-0000-0000-0000-000000000506',(SELECT id FROM products WHERE sku='SKU-001'),2,18900,17900,1000,'Bulk discount','00000000-0000-0000-0000-000000000305',now() - interval '1 day',now() - interval '1 day'),
  ('00000000-0000-0000-0000-000000000611','00000000-0000-0000-0000-000000000506',(SELECT id FROM products WHERE sku='SKU-006'),4,4500,4200,300,'Promo','00000000-0000-0000-0000-000000000305',now() - interval '1 day',now() - interval '1 day'),
  ('00000000-0000-0000-0000-000000000612','00000000-0000-0000-0000-000000000507',(SELECT id FROM products WHERE sku='SKU-008'),2,4900,4900,0,'','00000000-0000-0000-0000-000000000302',now() - interval '8 hours',now() - interval '8 hours'),
  ('00000000-0000-0000-0000-000000000613','00000000-0000-0000-0000-000000000507',(SELECT id FROM products WHERE sku='SKU-005'),3,3900,3900,0,'','00000000-0000-0000-0000-000000000302',now() - interval '8 hours',now() - interval '8 hours'),
  ('00000000-0000-0000-0000-000000000614','00000000-0000-0000-0000-000000000507',(SELECT id FROM products WHERE sku='SKU-002'),3,2900,2900,0,'','00000000-0000-0000-0000-000000000302',now() - interval '8 hours',now() - interval '8 hours'),
  ('00000000-0000-0000-0000-000000000615','00000000-0000-0000-0000-000000000508',(SELECT id FROM products WHERE sku='SKU-004'),1,5900,5900,0,'','00000000-0000-0000-0000-000000000306',now() - interval '2 hours',now() - interval '2 hours'),
  ('00000000-0000-0000-0000-000000000616','00000000-0000-0000-0000-000000000508',(SELECT id FROM products WHERE sku='SKU-006'),2,4500,4500,0,'','00000000-0000-0000-0000-000000000306',now() - interval '2 hours',now() - interval '2 hours')
ON CONFLICT (id) DO UPDATE SET
  sale_id=EXCLUDED.sale_id,
  product_id=EXCLUDED.product_id,
  quantity=EXCLUDED.quantity,
  original_price=EXCLUDED.original_price,
  final_price=EXCLUDED.final_price,
  discount_amount=EXCLUDED.discount_amount,
  discount_reason=EXCLUDED.discount_reason,
  employee_id=EXCLUDED.employee_id,
  created_at=EXCLUDED.created_at,
  updated_at=EXCLUDED.updated_at,
  deleted_at=NULL;

INSERT INTO payments (id, sale_id, payment_method, amount, created_at, updated_at)
VALUES
  ('00000000-0000-0000-0000-000000000701','00000000-0000-0000-0000-000000000501','CASH',19100,now() - interval '6 days',now() - interval '6 days'),
  ('00000000-0000-0000-0000-000000000702','00000000-0000-0000-0000-000000000502','PROMPTPAY',17700,now() - interval '5 days',now() - interval '5 days'),
  ('00000000-0000-0000-0000-000000000703','00000000-0000-0000-0000-000000000503','BANK_TRANSFER',21100,now() - interval '4 days',now() - interval '4 days'),
  ('00000000-0000-0000-0000-000000000704','00000000-0000-0000-0000-000000000504','CREDIT_CARD',34800,now() - interval '3 days',now() - interval '3 days'),
  ('00000000-0000-0000-0000-000000000705','00000000-0000-0000-0000-000000000505','CASH',27900,now() - interval '2 days',now() - interval '2 days'),
  ('00000000-0000-0000-0000-000000000706','00000000-0000-0000-0000-000000000506','PROMPTPAY',48100,now() - interval '1 day',now() - interval '1 day'),
  ('00000000-0000-0000-0000-000000000707','00000000-0000-0000-0000-000000000507','CASH',24700,now() - interval '8 hours',now() - interval '8 hours'),
  ('00000000-0000-0000-0000-000000000708','00000000-0000-0000-0000-000000000508','CREDIT_CARD',14700,now() - interval '2 hours',now() - interval '2 hours')
ON CONFLICT (id) DO UPDATE SET
  sale_id=EXCLUDED.sale_id,
  payment_method=EXCLUDED.payment_method,
  amount=EXCLUDED.amount,
  created_at=EXCLUDED.created_at,
  updated_at=EXCLUDED.updated_at,
  deleted_at=NULL;

INSERT INTO inventory_movements (id, branch_id, product_id, movement_type, quantity, reference_id, created_by, created_at, updated_at)
VALUES
  ('00000000-0000-0000-0000-000000000801',(SELECT id FROM branches WHERE code='BR01'),(SELECT id FROM products WHERE sku='SKU-001'),'SALE',-1,'00000000-0000-0000-0000-000000000501','00000000-0000-0000-0000-000000000302',now() - interval '6 days',now() - interval '6 days'),
  ('00000000-0000-0000-0000-000000000802',(SELECT id FROM branches WHERE code='BR01'),(SELECT id FROM products WHERE sku='SKU-003'),'SALE',-1,'00000000-0000-0000-0000-000000000501','00000000-0000-0000-0000-000000000302',now() - interval '6 days',now() - interval '6 days'),
  ('00000000-0000-0000-0000-000000000803',(SELECT id FROM branches WHERE code='BR01'),(SELECT id FROM products WHERE sku='SKU-004'),'SALE',-3,'00000000-0000-0000-0000-000000000502','00000000-0000-0000-0000-000000000302',now() - interval '5 days',now() - interval '5 days'),
  ('00000000-0000-0000-0000-000000000804',(SELECT id FROM branches WHERE code='BR02'),(SELECT id FROM products WHERE sku='SKU-002'),'SALE',-4,'00000000-0000-0000-0000-000000000503','00000000-0000-0000-0000-000000000305',now() - interval '4 days',now() - interval '4 days'),
  ('00000000-0000-0000-0000-000000000805',(SELECT id FROM branches WHERE code='BR02'),(SELECT id FROM products WHERE sku='SKU-007'),'SALE',-1,'00000000-0000-0000-0000-000000000503','00000000-0000-0000-0000-000000000305',now() - interval '4 days',now() - interval '4 days'),
  ('00000000-0000-0000-0000-000000000806',(SELECT id FROM branches WHERE code='BR03'),(SELECT id FROM products WHERE sku='SKU-010'),'SALE',-2,'00000000-0000-0000-0000-000000000504','00000000-0000-0000-0000-000000000306',now() - interval '3 days',now() - interval '3 days'),
  ('00000000-0000-0000-0000-000000000807',(SELECT id FROM branches WHERE code='BR03'),(SELECT id FROM products WHERE sku='SKU-005'),'SALE',-2,'00000000-0000-0000-0000-000000000504','00000000-0000-0000-0000-000000000306',now() - interval '3 days',now() - interval '3 days'),
  ('00000000-0000-0000-0000-000000000808',(SELECT id FROM branches WHERE code='HQ'),(SELECT id FROM products WHERE sku='SKU-012'),'SALE',-1,'00000000-0000-0000-0000-000000000505','00000000-0000-0000-0000-000000000301',now() - interval '2 days',now() - interval '2 days'),
  ('00000000-0000-0000-0000-000000000809',(SELECT id FROM branches WHERE code='HQ'),(SELECT id FROM products WHERE sku='SKU-009'),'SALE',-1,'00000000-0000-0000-0000-000000000505','00000000-0000-0000-0000-000000000301',now() - interval '2 days',now() - interval '2 days'),
  ('00000000-0000-0000-0000-000000000810',(SELECT id FROM branches WHERE code='BR02'),(SELECT id FROM products WHERE sku='SKU-001'),'SALE',-2,'00000000-0000-0000-0000-000000000506','00000000-0000-0000-0000-000000000305',now() - interval '1 day',now() - interval '1 day'),
  ('00000000-0000-0000-0000-000000000811',(SELECT id FROM branches WHERE code='BR02'),(SELECT id FROM products WHERE sku='SKU-006'),'SALE',-4,'00000000-0000-0000-0000-000000000506','00000000-0000-0000-0000-000000000305',now() - interval '1 day',now() - interval '1 day'),
  ('00000000-0000-0000-0000-000000000812',(SELECT id FROM branches WHERE code='BR01'),(SELECT id FROM products WHERE sku='SKU-008'),'SALE',-2,'00000000-0000-0000-0000-000000000507','00000000-0000-0000-0000-000000000302',now() - interval '8 hours',now() - interval '8 hours'),
  ('00000000-0000-0000-0000-000000000813',(SELECT id FROM branches WHERE code='BR01'),(SELECT id FROM products WHERE sku='SKU-005'),'SALE',-3,'00000000-0000-0000-0000-000000000507','00000000-0000-0000-0000-000000000302',now() - interval '8 hours',now() - interval '8 hours'),
  ('00000000-0000-0000-0000-000000000814',(SELECT id FROM branches WHERE code='BR01'),(SELECT id FROM products WHERE sku='SKU-002'),'SALE',-3,'00000000-0000-0000-0000-000000000507','00000000-0000-0000-0000-000000000302',now() - interval '8 hours',now() - interval '8 hours'),
  ('00000000-0000-0000-0000-000000000815',(SELECT id FROM branches WHERE code='BR03'),(SELECT id FROM products WHERE sku='SKU-004'),'SALE',-1,'00000000-0000-0000-0000-000000000508','00000000-0000-0000-0000-000000000306',now() - interval '2 hours',now() - interval '2 hours'),
  ('00000000-0000-0000-0000-000000000816',(SELECT id FROM branches WHERE code='BR03'),(SELECT id FROM products WHERE sku='SKU-006'),'SALE',-2,'00000000-0000-0000-0000-000000000508','00000000-0000-0000-0000-000000000306',now() - interval '2 hours',now() - interval '2 hours'),
  ('00000000-0000-0000-0000-000000000817',(SELECT id FROM branches WHERE code='BR02'),(SELECT id FROM products WHERE sku='SKU-006'),'RETURN',2,'00000000-0000-0000-0000-000000000506','00000000-0000-0000-0000-000000000307',now() - interval '18 hours',now() - interval '18 hours'),
  ('00000000-0000-0000-0000-000000000818',(SELECT id FROM branches WHERE code='BR01'),(SELECT id FROM products WHERE sku='SKU-011'),'RECEIVE',12,'Demo receive stock','00000000-0000-0000-0000-000000000301',now() - interval '7 days',now() - interval '7 days'),
  ('00000000-0000-0000-0000-000000000819',(SELECT id FROM branches WHERE code='BR03'),(SELECT id FROM products WHERE sku='SKU-008'),'ADJUSTMENT',-3,'Cycle count adjustment','00000000-0000-0000-0000-000000000301',now() - interval '12 hours',now() - interval '12 hours'),
  ('00000000-0000-0000-0000-000000000820',(SELECT id FROM branches WHERE code='HQ'),(SELECT id FROM products WHERE sku='SKU-002'),'TRANSFER_OUT',-10,'00000000-0000-0000-0000-000000000901','00000000-0000-0000-0000-000000000301',now() - interval '20 hours',now() - interval '20 hours'),
  ('00000000-0000-0000-0000-000000000821',(SELECT id FROM branches WHERE code='BR01'),(SELECT id FROM products WHERE sku='SKU-002'),'TRANSFER_IN',10,'00000000-0000-0000-0000-000000000901','00000000-0000-0000-0000-000000000301',now() - interval '20 hours',now() - interval '20 hours')
ON CONFLICT (id) DO UPDATE SET
  branch_id=EXCLUDED.branch_id,
  product_id=EXCLUDED.product_id,
  movement_type=EXCLUDED.movement_type,
  quantity=EXCLUDED.quantity,
  reference_id=EXCLUDED.reference_id,
  created_by=EXCLUDED.created_by,
  created_at=EXCLUDED.created_at,
  updated_at=EXCLUDED.updated_at,
  deleted_at=NULL;

INSERT INTO transfers (id, from_branch_id, to_branch_id, status, requested_by, approved_by, created_at, updated_at)
VALUES
  ('00000000-0000-0000-0000-000000000901',(SELECT id FROM branches WHERE code='HQ'),(SELECT id FROM branches WHERE code='BR01'),'COMPLETED','00000000-0000-0000-0000-000000000301','00000000-0000-0000-0000-000000000301',now() - interval '20 hours',now() - interval '20 hours'),
  ('00000000-0000-0000-0000-000000000902',(SELECT id FROM branches WHERE code='BR01'),(SELECT id FROM branches WHERE code='BR03'),'PENDING','00000000-0000-0000-0000-000000000301',NULL,now() - interval '1 hour',now() - interval '1 hour'),
  ('00000000-0000-0000-0000-000000000903',(SELECT id FROM branches WHERE code='BR02'),(SELECT id FROM branches WHERE code='HQ'),'APPROVED','00000000-0000-0000-0000-000000000301','00000000-0000-0000-0000-000000000301',now() - interval '30 minutes',now() - interval '30 minutes'),
  ('00000000-0000-0000-0000-000000000904',(SELECT id FROM branches WHERE code='BR03'),(SELECT id FROM branches WHERE code='BR02'),'REJECTED','00000000-0000-0000-0000-000000000301','00000000-0000-0000-0000-000000000301',now() - interval '2 days',now() - interval '2 days')
ON CONFLICT (id) DO UPDATE SET
  from_branch_id=EXCLUDED.from_branch_id,
  to_branch_id=EXCLUDED.to_branch_id,
  status=EXCLUDED.status,
  requested_by=EXCLUDED.requested_by,
  approved_by=EXCLUDED.approved_by,
  created_at=EXCLUDED.created_at,
  updated_at=EXCLUDED.updated_at,
  deleted_at=NULL;

INSERT INTO transfer_items (id, transfer_id, product_id, quantity, created_at, updated_at)
VALUES
  ('00000000-0000-0000-0000-000000000911','00000000-0000-0000-0000-000000000901',(SELECT id FROM products WHERE sku='SKU-002'),10,now() - interval '20 hours',now() - interval '20 hours'),
  ('00000000-0000-0000-0000-000000000912','00000000-0000-0000-0000-000000000902',(SELECT id FROM products WHERE sku='SKU-011'),2,now() - interval '1 hour',now() - interval '1 hour'),
  ('00000000-0000-0000-0000-000000000913','00000000-0000-0000-0000-000000000903',(SELECT id FROM products WHERE sku='SKU-007'),3,now() - interval '30 minutes',now() - interval '30 minutes'),
  ('00000000-0000-0000-0000-000000000914','00000000-0000-0000-0000-000000000904',(SELECT id FROM products WHERE sku='SKU-008'),4,now() - interval '2 days',now() - interval '2 days')
ON CONFLICT (id) DO UPDATE SET
  transfer_id=EXCLUDED.transfer_id,
  product_id=EXCLUDED.product_id,
  quantity=EXCLUDED.quantity,
  created_at=EXCLUDED.created_at,
  updated_at=EXCLUDED.updated_at,
  deleted_at=NULL;

INSERT INTO audit_logs (id, user_id, action, entity_type, entity_id, old_data, new_data, ip_address, created_at, updated_at)
VALUES
  ('00000000-0000-0000-0000-000000000a01','00000000-0000-0000-0000-000000000301','CREATE','product',(SELECT id FROM products WHERE sku='SKU-012'),NULL,'{"name":"Price Override Demo"}','127.0.0.1',now() - interval '7 days',now() - interval '7 days'),
  ('00000000-0000-0000-0000-000000000a02','00000000-0000-0000-0000-000000000301','UPDATE','inventory',(SELECT id FROM products WHERE sku='SKU-011'),'{"quantity":0}','{"quantity":3,"reorder_threshold":12}','127.0.0.1',now() - interval '6 days',now() - interval '6 days'),
  ('00000000-0000-0000-0000-000000000a03','00000000-0000-0000-0000-000000000307','REFUND','sale','00000000-0000-0000-0000-000000000506',NULL,'{"items":1,"quantity":2}','127.0.0.1',now() - interval '18 hours',now() - interval '18 hours'),
  ('00000000-0000-0000-0000-000000000a04','00000000-0000-0000-0000-000000000301','TRANSFER','inventory','00000000-0000-0000-0000-000000000901',NULL,'{"status":"COMPLETED","quantity":10}','127.0.0.1',now() - interval '20 hours',now() - interval '20 hours')
ON CONFLICT (id) DO UPDATE SET
  user_id=EXCLUDED.user_id,
  action=EXCLUDED.action,
  entity_type=EXCLUDED.entity_type,
  entity_id=EXCLUDED.entity_id,
  old_data=EXCLUDED.old_data,
  new_data=EXCLUDED.new_data,
  ip_address=EXCLUDED.ip_address,
  created_at=EXCLUDED.created_at,
  updated_at=EXCLUDED.updated_at,
  deleted_at=NULL;
