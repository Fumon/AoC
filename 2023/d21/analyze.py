import csv

def is_multiple(value, previous_values):
    for prev_value in previous_values:
        if value % prev_value == 0 or prev_value % value == 0:
            return True
    return False

def analyze_csv(file_path):
    with open(file_path, 'r') as file:
        reader = csv.reader(file)
        previous_values = []

        for row in reader:
            if row:  # Ensure row is not empty
                try:
                    current_value = int(row[1])  # Assuming the second column has index 1
                    if is_multiple(current_value, previous_values):
                        print(f"The value {current_value} in row {reader.line_num} is a multiple of a previous value.")
                    sidelength = 2*(int(row[0]) + 1) - 1
                    print(f"Adding {sidelength ** 2}")
                    previous_values.append(current_value / (sidelength ** 2))
                except ValueError:
                    print(f"Non-integer value in row {reader.line_num}, skipping.")

        print(previous_values)

# Example usage
file_path = '131steps.csv'
analyze_csv(file_path)
