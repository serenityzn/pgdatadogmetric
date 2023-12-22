CREATE TABLE products (
                          product_no integer,
                          name text,
                          price numeric
);

CREATE TABLE tpart_charge_20231220w3 (
    end_time timestamp,
    product_no integer,
    name text
);
INSERT INTO tpart_charge_20231220w3 VALUES ('2023-12-20 00:00:00', 1, 'test1');
INSERT INTO tpart_charge_20231220w3 VALUES ('2023-12-20 01:10:00', 2, 'test2');
INSERT INTO tpart_charge_20231220w3 VALUES ('2023-12-20 03:23:00', 3, 'test3');
INSERT INTO tpart_charge_20231220w3 VALUES ('2023-12-20 14:00:00', 4, 'test4');
INSERT INTO tpart_charge_20231220w3 VALUES ('2023-12-20 14:31:00', 1, 'test');
INSERT INTO tpart_charge_20231220w3 VALUES ('2023-12-20 14:32:00', 1, 'test');
INSERT INTO tpart_charge_20231220w3 VALUES ('2023-12-20 14:33:00', 1, 'test');
INSERT INTO tpart_charge_20231220w3 VALUES ('2023-12-20 14:34:00', 1, 'test');
INSERT INTO tpart_charge_20231220w3 VALUES ('2023-12-20 14:35:00', 1, 'test');
INSERT INTO tpart_charge_20231220w3 VALUES ('2023-12-21 01:30:00', 1, 'test');
