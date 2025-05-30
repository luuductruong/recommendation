-- Postgres
DO
$$
    DECLARE
        user_ids     TEXT[] := ARRAY ['User_id_001', 'User_id_002', 'User_id_003', 'User_id_004', 'User_id_005'];
        category_ids TEXT[] := ARRAY ['Cat_001', 'Cat_002', 'Cat_003', 'Cat_004', 'Cat_005'];
        num_product  INT    := 10;
        existing_product_count INT;
        uid          TEXT;
        cid          TEXT;
BEGIN
    -- Đếm số sản phẩm. Do tạo id tăng dần, nên có thể đếm số lượng để biết id đang tới bao nhiêu
    SELECT COUNT(*) INTO existing_product_count FROM product;
        -- Tạo sản phẩm
    IF existing_product_count+1 < num_product THEN
        FOR i IN existing_product_count..num_product
            LOOP cid := category_ids[(trunc(random() * array_length(category_ids, 1)) + 1)::int];
                INSERT INTO product (product_id, name, price, category_id)
                VALUES (i,
                        'Product ' || i,
                        (random() * 100)::int + 10,
                        cid);
            END LOOP;
    END IF;

    -- Tạo lịch sử xem sản phẩm
    FOR i IN 1..100
        LOOP uid := user_ids[(trunc(random() * array_length(user_ids, 1)) + 1)::int];
            INSERT INTO user_view_history (id, product_id, user_id, view_at)
            VALUES (gen_random_uuid(),
                    (random() * num_product)::int + 1,
                    uid,
                    NOW() - (INTERVAL '3 minutes 5 seconds' * ((random() * 10)::int + 1)));

            -- Ngủ 0.091 giây sau mỗi 10 insert
            IF i % 10 = 0 THEN
                PERFORM pg_sleep(0.091);
            END IF;
        END LOOP;
END;
$$;
