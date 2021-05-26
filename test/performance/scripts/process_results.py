import csv, json, os

# get current working directory
cwd = os.getcwd()

# csv data as parsed from the raw results
raw_results = {}

# relative path of the results produced by the perf tests
csv_file_name = f'{cwd}/test/performance/reports/perf_test_stats.csv'

# template with empty results (to which the actual results be injected)
json_template_file_name = f'{cwd}/test/performance/templates/results.json'

# processed data filename relative path
processed_json_data = f'{cwd}/test/performance/reports/perf_test_stats.json'

# convert csv to lower case
f = open(csv_file_name, 'r')
text = f.read()

lines = text.lower()
with open(csv_file_name, 'w') as out:
  out.writelines(lines)

with open(json_template_file_name) as json_data:
  schema_data = json.load(json_data)

# read csv 
with open(csv_file_name) as csvRaw:
  csvReader = csv.DictReader(csvRaw)
  for result in csvReader:
    # populate aggregated stats
    if result['type'] == '' and result['name'] == 'aggregated':
      for k, v in result.items():
        if k != 'name' and k != 'type':
          for schema_endpoint in schema_data['endpoints']:
            if 'aggregated' in schema_endpoint:
              schema_endpoint['aggregated'][k] = float(v)
  
    # populate endpoints stats
    for schema_endpoint in schema_data['endpoints']:
      if result['name'] in schema_endpoint and result['type'] in schema_endpoint[result['name']]:
        for k,v in result.items():
          if k != 'name' and k != 'type':
            schema_endpoint[result['name']][result['type']][k] = float(v)

# persist processed JSON results
with open(processed_json_data, 'w') as outfile:
  json.dump(schema_data, outfile)
