-- //{
-- //id
-- //country
-- //city
-- //street
-- //}


create table if not exists address (
                                       id uuid PRIMARY KEY,
                                       country varchar(10),
                                       city varchar(30),
                                       street varchar(100)
);


--     images
-- {
--     id : UUID
--     image: bytea
-- }


create table if not exists images
(
    id uuid primary key,
    image bytea
);

--     supplier
-- {
--     id
--     name
--     address_id
--     phone_number
-- }


create table if not exists supplier
(
    id uuid primary key,
    name varchar(300),
    address_id uuid,
    phone_number varchar(50),
    foreign key (address_id) references address(id)
);



create table if not exists product
(
    id uuid primary key ,
    name varchar(100),
    category varchar(100),
    price float,
    available_stock int, --// число закупленных экземпляров товара
    last_update_date timestamp, --// число последней закупки
    supplier_id uuid,
    foreign key (supplier_id) references supplier(id),
    image_id UUID,
    foreign key (image_id) references images(id)
);


create table if not exists client
(
    id                uuid primary key,
    client_name       varchar(50),
    client_surname    varchar(50),
    birthday          timestamp,
    gender            varchar(10),
    registration_date timestamp,
    address_id        uuid,
    foreign key (address_id) references address(id)
);