import csv
import random

num_rows = 10000
output_file = "servers_mixed_input.csv"

def random_ipv4(valid=True):
    if valid:
        return f"{random.randint(1, 255)}.{random.randint(0, 255)}.{random.randint(0, 255)}.{random.randint(1, 254)}"
    else:
        return random.choice([
            "999.999.999.999",
            "invalid_ip",
            "",
            "abc.def.ghi.jkl"
        ])

def random_name(valid=True, i=0):
    if valid:
        return f"server{i}"
    else:
        return random.choice([
            "",
            "!!!@@@",
            "server"*50
        ])

with open(output_file, mode="w", newline="") as file:
    writer = csv.writer(file)

    writer.writerow(["name", "ipv4_address"])

    for i in range(1, num_rows + 1):
        is_valid = random.random() > 0.3  # 70% valid

        name = random_name(is_valid, i)
        ipv4 = random_ipv4(is_valid)

        writer.writerow([name, ipv4])

print(f"{output_file} created with mixed valid/invalid rows")