# Nginx GeoLite2 Region Blocking for Docker Compose

This deployment keeps all API gateway endpoints available from every region, but redirects mainland China visitors away from the web UI to `/unsupported-region`.

The frontend only displays the blocked-region page. Nginx is the authoritative layer for page access.

## Request Flow

```text
Visitor -> Cloudflare -> Host Nginx -> 127.0.0.1:8080 -> sub2api Docker Compose service
```

## 1. Install Dependencies

Ubuntu / Debian:

```bash
sudo apt update
sudo apt install -y nginx libnginx-mod-http-geoip2 geoipupdate mmdb-bin
```

If the module is not auto-loaded, add this at the top level of `/etc/nginx/nginx.conf`:

```nginx
load_module modules/ngx_http_geoip2_module.so;
```

## 2. Download GeoLite2 Country

Edit `/etc/GeoIP.conf`:

```ini
AccountID YOUR_ACCOUNT_ID
LicenseKey YOUR_LICENSE_KEY
EditionIDs GeoLite2-Country
DatabaseDirectory /var/lib/GeoIP
```

Download and verify:

```bash
sudo mkdir -p /var/lib/GeoIP
sudo geoipupdate -v
ls -lh /var/lib/GeoIP/GeoLite2-Country.mmdb
```

Only `GeoLite2-Country.mmdb` is required for CN blocking. Do not configure `GeoLite2-City.mmdb` unless the file exists.

## 3. Restore Real Visitor IP from Cloudflare

Create `/etc/nginx/conf.d/cloudflare-real-ip.conf`:

```nginx
# Cloudflare IPv4
set_real_ip_from 173.245.48.0/20;
set_real_ip_from 103.21.244.0/22;
set_real_ip_from 103.22.200.0/22;
set_real_ip_from 103.31.4.0/22;
set_real_ip_from 141.101.64.0/18;
set_real_ip_from 108.162.192.0/18;
set_real_ip_from 190.93.240.0/20;
set_real_ip_from 188.114.96.0/20;
set_real_ip_from 197.234.240.0/22;
set_real_ip_from 198.41.128.0/17;
set_real_ip_from 162.158.0.0/15;
set_real_ip_from 104.16.0.0/13;
set_real_ip_from 104.24.0.0/14;
set_real_ip_from 172.64.0.0/13;
set_real_ip_from 131.0.72.0/22;

# Cloudflare IPv6
set_real_ip_from 2400:cb00::/32;
set_real_ip_from 2606:4700::/32;
set_real_ip_from 2803:f800::/32;
set_real_ip_from 2405:b500::/32;
set_real_ip_from 2405:8100::/32;
set_real_ip_from 2a06:98c0::/29;
set_real_ip_from 2c0f:f248::/32;

real_ip_header CF-Connecting-IP;
real_ip_recursive on;
```

After this, `$remote_addr` is the real visitor IP. The origin server IP, such as `13.231.40.7`, must not be used in `set_real_ip_from`.

## 4. Configure GeoLite2 and Page-Only Blocking

Create `/etc/nginx/conf.d/geoip2-region-block.conf`:

```nginx
geoip2 /var/lib/GeoIP/GeoLite2-Country.mmdb {
    auto_reload 1d;
    $geoip2_country_code country iso_code;
    $geoip2_country_name country names en;
}

map $geoip2_country_code $region_blocked {
    default 0;
    CN 1;
}

# Paths that must never be region-blocked.
map $uri $region_block_exempt_path {
    default 0;

    # Block page and browser assets
    /unsupported-region 1;
    /region-check 1;
    /health 1;
    /favicon.ico 1;
    /robots.txt 1;
    ~^/assets/ 1;

    # API gateway forwarding endpoints. These remain available from every region.
    ~^/v1/ 1;
    /v1 1;
    ~^/v1beta/ 1;
    /v1beta 1;
    ~^/openai/v1/ 1;
    /openai/v1 1;
    ~^/antigravity/ 1;
    ~^/backend-api/codex/ 1;
    ~^/responses(?:/.*)?$ 1;
    ~^/chat/completions$ 1;
    ~^/embeddings$ 1;
    ~^/images/generations$ 1;
    ~^/images/edits$ 1;
}

map "$region_blocked:$region_block_exempt_path" $redirect_unsupported_region {
    default 0;
    "1:0" 1;
}
```

## 5. Configure the Site Server

Keep your existing Docker Compose upstream at `127.0.0.1:8080`. Add `/region-check` and the redirect check before `location /`.

```nginx
server {
    listen 443 ssl http2;
    server_name your-domain.com;

    ssl_certificate /path/fullchain.pem;
    ssl_certificate_key /path/privkey.pem;

    # Used by the unsupported-region page retry button. This checks only page access.
    location = /region-check {
        add_header Cache-Control "no-store";

        if ($region_blocked) {
            return 403;
        }

        return 204;
    }

    # Web page access is blocked for mainland China IPs. API gateway paths are exempt by map.
    if ($redirect_unsupported_region) {
        return 302 /unsupported-region?ip=$remote_addr&country=$geoip2_country_name&path=$uri;
    }

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;

        # WebSocket support
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_cache_bypass $http_upgrade;

        # Long-running gateway requests
        proxy_connect_timeout 600s;
        proxy_send_timeout 600s;
        proxy_read_timeout 600s;

        # Streaming / SSE
        proxy_buffering off;
        proxy_request_buffering off;
        proxy_cache off;
        gzip off;

        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

If you also listen on port 80, redirect HTTP to HTTPS in a separate server block.

## 6. Validate and Reload

```bash
sudo nginx -t
sudo systemctl reload nginx
```

Check Docker Compose service:

```bash
curl -I http://127.0.0.1:8080/health
```

## 7. Test

Temporarily force page blocking:

```nginx
map $geoip2_country_code $region_blocked {
    default 1;
}
```

Reload Nginx, then verify:

```bash
sudo nginx -t
sudo systemctl reload nginx
```

- Visiting `/home` or `/status` should redirect to `/unsupported-region`.
- Calling gateway endpoints such as `/v1/messages`, `/v1/chat/completions`, `/v1beta/models`, `/responses`, or `/antigravity/v1/messages` must not redirect.
- Clicking retry on `/unsupported-region` calls `/region-check`; it returns `403` while blocked and `204` when page access is allowed.

Restore production blocking after the test:

```nginx
map $geoip2_country_code $region_blocked {
    default 0;
    CN 1;
}
```

Then reload Nginx again.
