import json
from datetime import datetime

import bson
from faker import Faker
from random import choice, randint

fake = Faker()


def generate_users_data(dishes_data, num_users=100):
    users_data = []
    for _ in range(num_users):
        first_name = fake.first_name()
        last_name = fake.last_name()
        f1 = choice(dishes_data)
        f2 = choice(dishes_data)
        f3 = choice(dishes_data)
        sum = f1["price"] + f2["price"] + f3["price"]

        user = {
            "first_name": first_name,
            "last_name": last_name,
            "email": fake.email(),
            "password": "$2a$12$5qg8degfG15UmApJ5U7NO.5QiWoXVDKnmOzZevHltgjleEd54ZC9q",  # Abcd_1234
            "phone_number": fake.phone_number(),
            "verification_token": fake.sha256(raw_output=False),
            "email_verified": True,
            "orders": [
                {
                    "_id": {"$oid": str(bson.ObjectId())},
                    "dishes": [f1, f2, f3],
                    "total_price": float(sum),
                    "ordered_date": {"$date": datetime.utcnow().isoformat()},

                }
            ],
            "cart": [choice(dishes_data), choice(dishes_data), choice(dishes_data)],
            "delivery":
                {
                    "full_name": f"{first_name} {last_name}",
                    "address": fake.street_address(),
                    "city": fake.city(),
                    "zip_code": fake.zipcode(),
                    "phone_number": fake.phone_number(),
                },
        }
        users_data.append(user)

    return users_data


with open('dishes_data.json', 'r', encoding='utf-8') as file:
    dishes_data = json.load(file)

users_data = generate_users_data(dishes_data, 100)
with open('users_data.json', 'w', encoding='utf-8') as file:
    json.dump(users_data, file, indent=4, ensure_ascii=False)

print("Файл 'users_data.json' успешно создан с фейковой информацией о пользователях.")
