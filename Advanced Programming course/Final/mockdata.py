import json
import random
import faker

fake = faker.Faker()


def generate_real_dish_data(n):
    dishes_info = [
        {"name": "Pasta Carbonara",
         "img_URL": "https://upload.wikimedia.org/wikipedia/commons/3/32/Spaghetti_alla_Carbonara_%282%29.jpg"},
        {"name": "Margherita Pizza",
         "img_URL": "https://upload.wikimedia.org/wikipedia/commons/6/66/Pizza_Margherita_stu_spivack.jpg"},
        {"name": "Caesar Salad",
         "img_URL": "https://upload.wikimedia.org/wikipedia/commons/1/15/Caesar_salad_%281%29.jpg"},
        {"name": "Beef Bourguignon",
         "img_URL": "https://upload.wikimedia.org/wikipedia/commons/6/6e/Boeuf_bourguignon.jpg"},
        {"name": "Fish and Chips",
         "img_URL": "https://upload.wikimedia.org/wikipedia/commons/f/ff/Fish_and_chips_blackpool.jpg"},
        {"name": "Chicken Curry",
         "img_URL": "https://upload.wikimedia.org/wikipedia/commons/4/4b/Chicken_makhani_curry.jpg"},
        {"name": "Vegetable Stir Fry",
         "img_URL": "https://upload.wikimedia.org/wikipedia/commons/a/ac/Stir_Fried_Veg.JPG"},
        {"name": "Lamb Gyro", "img_URL": "https://upload.wikimedia.org/wikipedia/commons/1/1b/Gyros2.jpg"},
        {"name": "Quinoa Salad", "img_URL": "https://upload.wikimedia.org/wikipedia/commons/3/3c/Quinoa_Salad.jpg"},
        {"name": "Tofu Scramble", "img_URL": "https://upload.wikimedia.org/wikipedia/commons/a/a3/Scrambled_Tofu.png"},
        {"name": "Mushroom Risotto",
         "img_URL": "https://upload.wikimedia.org/wikipedia/commons/2/2c/Risotto_funghi_pancetta.jpg"},
        {"name": "Tomato Soup", "img_URL": "https://upload.wikimedia.org/wikipedia/commons/5/58/Tomato_soup.jpg"},
        {"name": "BBQ Ribs", "img_URL": "https://upload.wikimedia.org/wikipedia/commons/3/32/BBQ_Ribs.jpg"},
        {"name": "Falafel Wrap", "img_URL": "https://upload.wikimedia.org/wikipedia/commons/6/67/Falafel_wrap.jpg"},
        {"name": "Spinach Quiche", "img_URL": "https://upload.wikimedia.org/wikipedia/commons/9/94/Quiche.jpg"},
        {"name": "Duck Confit", "img_URL": "https://upload.wikimedia.org/wikipedia/commons/b/b8/Confit_de_Canard.jpg"},
        {"name": "Vegan Burger", "img_URL": "https://upload.wikimedia.org/wikipedia/commons/0/0b/Vegan_burger.jpg"},
        {"name": "Shrimp Paella", "img_URL": "https://upload.wikimedia.org/wikipedia/commons/e/ed/Paella_mixta.jpg"},
        {"name": "Ratatouille", "img_URL": "https://upload.wikimedia.org/wikipedia/commons/2/29/Ratatouille_dish.jpg"},
        {"name": "Sushi Roll",
         "img_URL": "https://upload.wikimedia.org/wikipedia/commons/0/0b/Inside-Out_California_Roll.jpg"}
    ]

    data = []
    for _ in range(n):
        dish = random.choice(dishes_info)
        name = dish["name"]
        img_url = dish["img_URL"]
        description = fake.text(max_nb_chars=200)
        price = round(random.uniform(5, 30), 2)
        weight = round(random.uniform(100, 500), 2)
        proteins = round(random.uniform(5, 30), 2)
        fats = round(random.uniform(5, 30), 2)
        carbohydrates = round(random.uniform(5, 30), 2)

        data.append({
            "name": name,
            "description": description,
            "price": price,
            "weight": weight,
            "protein": proteins,
            "fats": fats,
            "carbohydrates": carbohydrates,
            "img_URL": img_url
        })

    return data

dish_data = generate_real_dish_data(50)
file_path = 'dishes_data.json'
with open(file_path, '+w') as file:
    json.dump(dish_data, file, indent=2)
