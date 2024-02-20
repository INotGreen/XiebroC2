with open('TestDLL.cs', 'r', encoding='utf-8') as file:
    content = file.read()
start = content.find('byte[] rawData = {')
end = content.find('};', start) + 1
byte_array_string = content[start:end]

cleaned_byte_array = (
    byte_array_string.replace('\n', '').replace(' ', '').replace('\t', '')
)
print(cleaned_byte_array)
