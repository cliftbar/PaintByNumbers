
from http.server import HTTPServer, BaseHTTPRequestHandler, SimpleHTTPRequestHandler
import socketserver


PORT = 8000


class CORSRequestHandler (SimpleHTTPRequestHandler):
    def end_headers (self):
        self.send_header('Access-Control-Allow-Origin', '*')
        SimpleHTTPRequestHandler.end_headers(self)

Handler = CORSRequestHandler

Handler.extensions_map={
    '.manifest': 'text/cache-manifest',
    '.html': 'text/html',
    '.png': 'image/png',
    '.jpg': 'image/jpg',
    '.svg': 'image/svg+xml',
    '.css': 'text/css',
    '.wasm': 'application/wasm',
    '.js': 'application/x-javascript',
    '': 'application/octet-stream', # Default
}

httpd = socketserver.TCPServer(("0.0.0.0", PORT), Handler)

print("serving at port", PORT)
httpd.serve_forever()
