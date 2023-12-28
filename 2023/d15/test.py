import re

tokens = set()
longest = 0
longest_token = ''
with open('input', 'r') as f:
    for line in f:
        for cmd in line.split(','):
            code = re.split('[-=]', cmd)[0]
            if len(code) > longest:
                longest = len(code)
                longest_token = code
            tokens.add(code)

print(tokens)
print(f"longest code: {longest_token} @ {longest}")
