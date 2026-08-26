$crt = Get-Content "E:\server\alnitak\server\conf\public.crt" -Raw
$key = Get-Content "E:\server\alnitak\server\conf\private.key" -Raw
$env:NITRO_SSL_CERT = $crt
$env:NITRO_SSL_KEY = $key
node .output/server/index.mjs
