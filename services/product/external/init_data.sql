-- PostgreSQL 15.13 (Debian 15.13-1.pgdg120+1)
-- Insert product data
DO
$$
    DECLARE
        category_ids TEXT[] := ARRAY ['Cat_001', 'Cat_002', 'Cat_003', 'Cat_004', 'Cat_005'];
        num_product  INT    := 10;
        existing_product_count INT;
        cid          TEXT;
    BEGIN
        SELECT COUNT(*) INTO existing_product_count FROM product;
        IF existing_product_count < num_product THEN
            FOR i IN existing_product_count+1..num_product LOOP
                    cid := category_ids[(trunc(random() * array_length(category_ids, 1)) + 1)::int];

                    INSERT INTO product (product_id, name, price, category_id)
                    VALUES (i,
                            'Product ' || i,
                            (random() * 100)::int + 10,
                            cid);
                END LOOP;
        END IF;
    END;
$$;


-- Insert history data
DO
$$
    DECLARE
        user_ids     TEXT[] := ARRAY ['User_id_001', 'User_id_002', 'User_id_003', 'User_id_004', 'User_id_005'];
        uid          TEXT;
        pid          INT;
        view_time    TIMESTAMP;
        cat_id       TEXT;
        num_product  INT := 10;
    BEGIN
        FOR i IN 1..100 LOOP
                uid := user_ids[(trunc(random() * array_length(user_ids, 1)) + 1)::int];
                pid := (random() * num_product)::int + 1;
                view_time := NOW() - (INTERVAL '1 minutes 26 seconds' * ((random() * 10)::int + 1));

                -- Insert view history
                INSERT INTO user_view_history (id, product_id, user_id, view_at)
                VALUES (
                           gen_random_uuid(), -- For postgres version >= 13. Older version (>=8.3) use uuid_generate_v4() instead
                           pid,
                           uid,
                           view_time
                       );

                -- get category_id from product
                SELECT category_id INTO cat_id FROM product WHERE product_id = pid;

                -- insert / update category_view_history
                INSERT INTO category_view_history (id, category_id, total_view, last_view_at)
                VALUES (
                           gen_random_uuid(),
                           cat_id,
                           1,
                           view_time
                       )
                ON CONFLICT (category_id) DO UPDATE
                    SET
                        total_view = category_view_history.total_view + 1,
                        last_view_at = GREATEST(category_view_history.last_view_at, EXCLUDED.last_view_at);

                -- Sleep for a while cuz I don't want to see the same time
                IF i % 7 = 0 THEN
                    PERFORM pg_sleep(0.261);
                END IF;
            END LOOP;
    END;
$$;
