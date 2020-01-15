-- 创建数据库
CREATE DATABASE crawler DEFAULT CHARSET = utf8mb4;

-- 创建商品表
CREATE TABLE items
(
    id          bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    website     varchar(20)         NULL     DEFAULT NULL,
    product_id  varchar(30)         NULL     DEFAULT NULL,
    internal_id varchar(30)         NULL     DEFAULT NULL,
    title       varchar(255)        NOT NULL,
    image       varchar(255)        NULL     DEFAULT NULL,
    created_at  timestamp           NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  timestamp           NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at  timestamp           NULL     DEFAULT NULL,
    CONSTRAINT item_pk
        PRIMARY KEY (id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;
CREATE UNIQUE INDEX item_production_id ON items (website, product_id);

-- 创建价格记录表
CREATE TABLE prices
(
    id         bigint(20) unsigned not null auto_increment,
    item_id    bigint(20) unsigned,
    price      decimal(10, 2)      not null default 0,
    created_at timestamp           not null default current_timestamp,
    constraint item_price_pk
        primary key (id)
)ENGINE = InnoDB
 DEFAULT CHARSET = utf8mb4
 COLLATE = utf8mb4_unicode_ci;
ALTER TABLE prices
    ADD CONSTRAINT price_item_id_fk
        FOREIGN KEY (item_id) REFERENCES items (id)
            ON UPDATE CASCADE ON DELETE CASCADE;


