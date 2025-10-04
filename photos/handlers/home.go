package handlers

import (
	"fmt"
	"net/http"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	html := `
<!DOCTYPE html>
<html>
<head>
	<title>Jane's Photo Service</title>
	<style>
		body {
			font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
			max-width: 800px;
			margin: 50px auto;
			padding: 20px;
			background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
			min-height: 100vh;
		}
		.card {
			background: white;
			border-radius: 12px;
			padding: 30px;
			box-shadow: 0 10px 30px rgba(0,0,0,0.2);
		}
		h1 { color: #333; margin-bottom: 10px; }
		.subtitle { color: #666; margin-bottom: 30px; }
		.info { background: #f8f9fa; padding: 15px; border-radius: 8px; margin: 15px 0; border-left: 4px solid #667eea; }
		code { background: #e9ecef; padding: 2px 6px; border-radius: 3px; font-family: monospace; }
		a { color: #667eea; text-decoration: none; font-weight: 600; }
		a:hover { text-decoration: underline; }
	</style>
</head>
<body>
	<div class="card">
		<h1>Jane's Photo Service</h1>
		<p class="subtitle">OAuth 1.0 Provider (Resource Server)</p>
		
		<div class="info">
			<strong>Service Running:</strong> http://localhost:8081<br>
			<strong>Role:</strong> OAuth Provider / Resource Server
		</div>

		<h2>OAuth Endpoints:</h2>
		<ul>
			<li><code>POST /oauth/initiate</code> - Request temporary credentials</li>
			<li><code>GET /oauth/authorize</code> - User authorization page</li>
			<li><code>POST /oauth/token</code> - Exchange for access token</li>
			<li><code>GET /api/photos</code> - Protected resource (requires OAuth)</li>
		</ul>

		<h2>Registered Client:</h2>
		<div class="info">
			<strong>Client Key:</strong> <code>dpf43f3p2l4k3l03</code><br>
			<strong>Client Secret:</strong> <code>kd94hf93k423kf44</code>
		</div>

		<h2>Test User:</h2>
		<div class="info">
			<strong>Username:</strong> <code>jane</code><br>
			<strong>Password:</strong> <code>password123</code>
		</div>

		<p style="margin-top: 30px;">
			To start the OAuth flow, go to the 
			<a href="http://localhost:8080" target="_blank">Printer Service →</a>
		</p>
	</div>
</body>
</html>
	`
	fmt.Fprint(w, html)
}
