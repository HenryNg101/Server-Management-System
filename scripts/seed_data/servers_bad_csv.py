import csv
import random

# Configuration
num_rows = 1000000
output_file = "servers.csv"
protocols = ["tcp", "udp", "http", "https", "ftp", "ssh"]

with open(output_file, mode="w", newline="") as file:
    writer = csv.writer(file)
    
    # Write header
    writer.writerow(["name", "status", "ipv4_address", "port", "protocol"])
    
    for i in range(1, num_rows + 1):
        name = f"server{i}"
        status = random.choice([True, False])
        ipv4 = "127.0.0.1"
        # port = 9999 + i  # 10000 to 19999
        port = "Random invalid port"
        protocol = random.choice(protocols)
        
        writer.writerow([name, status, ipv4, port, protocol])

print(f"CSV file '{output_file}' created with {num_rows} rows.")