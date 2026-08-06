import csv

num_rows = 10000
output_file = "servers_input.csv"

with open(output_file, mode="w", newline="") as file:
    writer = csv.writer(file)

    writer.writerow(["name", "ipv4_address"])

    for i in range(1, num_rows + 1):
        name = f"server{i}"
        ipv4 = "127.0.0.1"

        writer.writerow([name, ipv4])

print(f"{output_file} created with {num_rows} rows")