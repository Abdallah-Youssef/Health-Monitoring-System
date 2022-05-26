import json



f = open(f'messages.txt', 'r')
o = open(f'messages.csv', 'w')

for line in f:
  j = json.loads(line)
  o.write(f'{j["serviceName"].split("-")[1]},{j["Timestamp"]},{j["CPU"]},{j["RAM"]["Total"]},{j["RAM"]["Free"]},{j["Disk"]["Total"]},{j["Disk"]["Free"]}\n')

o.close()
f.close()
