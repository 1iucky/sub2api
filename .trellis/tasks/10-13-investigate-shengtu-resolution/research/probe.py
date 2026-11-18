import importlib.util
import json
import sys
import time
from pathlib import Path
from PIL import Image

SCRIPT = Path('/Users/liuliang/.codex/skills/shengtu-skill/scripts/generate_image.py')
spec = importlib.util.spec_from_file_location('shengtu_probe', SCRIPT)
skill = importlib.util.module_from_spec(spec)
spec.loader.exec_module(skill)
name, model, size, quality, fmt = sys.argv[1:]
root = Path(__file__).parent
out = root / (name + '.' + fmt)
record = {'case': name, 'model': model, 'size': size, 'quality': quality, 'output_format': fmt,
          'gateway': skill.DEFAULT_BASE_URL, 'started_utc': time.strftime('%Y-%m-%dT%H:%M:%SZ', time.gmtime())}
original_request = skill.request_json
def observed_request(url, key, payload):
    record['request'] = dict(payload)
    result = original_request(url, key, payload)
    record['response_metadata'] = {k: result[k] for k in ('model','size','quality','output_format','usage','created') if k in result}
    record['image_metadata'] = [{k:v for k,v in item.items() if k in ('model','size','quality','output_format')} for item in result.get('data',[])]
    return result
skill.request_json = observed_request
original_urlopen = skill.urllib.request.urlopen
def observed_urlopen(req, *args, **kwargs):
    response = original_urlopen(req, *args, **kwargs)
    record.setdefault('http', []).append({'status': response.status, 'headers': {k: response.headers[k] for k in ('x-request-id','x-upstream-request-id','content-type','server') if k in response.headers}})
    return response
skill.urllib.request.urlopen = observed_urlopen
sys.argv = [str(SCRIPT), '--model', model, '--size', size, '--quality', quality,
            '--output-format', fmt, '--prompt', 'A single red ceramic cup centered on a plain pale gray background, studio product photograph, no text.', '--out', str(out)]
start = time.monotonic()
try:
    code = skill.main()
    record['exit_code'] = code
    if code == 0:
        with Image.open(out) as im:
            record['actual_pixels'] = list(im.size)
            record['actual_format'] = im.format
        record['file'] = str(out)
        record['bytes'] = out.stat().st_size
except Exception as exc:
    record['error'] = str(exc)[:1500]
finally:
    record['elapsed_seconds'] = round(time.monotonic()-start, 2)
    (root / (name + '.json')).write_text(json.dumps(record, ensure_ascii=False, indent=2)+'\n')
    print(json.dumps(record, ensure_ascii=False), flush=True)
