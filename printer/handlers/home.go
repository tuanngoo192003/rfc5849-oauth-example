package handlers

import (
	"go-oauth1/printer/authorization"
	"html/template"
	"net/http"
)

func HandlerHome(w http.ResponseWriter, r *http.Request) {
	authorization.GetSession().GetMu().RLock()
	data := struct {
		TempToken   string
		AccessToken string
		Photos      []authorization.Photo
		Logs        []string
		HasAccess   bool
	}{
		TempToken:   authorization.GetSession().GetTempToken(),
		AccessToken: authorization.GetSession().GetAccessToken(),
		Photos:      authorization.GetSession().GetPhotos(),
		Logs:        authorization.GetSession().GetLogs(),
		HasAccess:   authorization.GetSession().GetAccessToken() != "",
	}
	authorization.GetSession().GetMu().RUnlock()

	tmpl := template.Must(template.New("home").Parse(`
<!DOCTYPE html>
<html>
<head>
	<title>Printer Service - OAuth Client</title>
	<style>
		* { margin: 0; padding: 0; box-sizing: border-box; }
		body {
			font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
			background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
			min-height: 100vh;
			padding: 20px;
		}
		.container {
			max-width: 1200px;
			margin: 0 auto;
		}
		.card {
			background: white;
			border-radius: 12px;
			padding: 30px;
			margin-bottom: 20px;
			box-shadow: 0 10px 30px rgba(0,0,0,0.2);
		}
		h1 {
			color: #333;
			margin-bottom: 10px;
			font-size: 32px;
		}
		.subtitle {
			color: #666;
			margin-bottom: 30px;
			font-size: 16px;
		}
		.hero {
			text-align: center;
			padding: 40px 20px;
			background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
			border-radius: 12px;
			color: white;
			margin-bottom: 30px;
		}
		.hero h2 {
			font-size: 28px;
			margin-bottom: 15px;
		}
		.hero p {
			font-size: 16px;
			margin-bottom: 25px;
			opacity: 0.9;
		}
		button {
			background: white;
			color: #667eea;
			border: none;
			padding: 15px 30px;
			border-radius: 8px;
			cursor: pointer;
			font-size: 16px;
			font-weight: 600;
			transition: all 0.3s;
		}
		button:hover {
			transform: translateY(-2px);
			box-shadow: 0 6px 20px rgba(255,255,255,0.3);
		}
		button:disabled {
			background: #ccc;
			color: #666;
			cursor: not-allowed;
			transform: none;
		}
		.step {
			background: #f8f9fa;
			padding: 20px;
			border-radius: 8px;
			margin-bottom: 20px;
			border-left: 4px solid #667eea;
		}
		.step h3 {
			color: #667eea;
			margin-bottom: 10px;
			font-size: 18px;
		}
		.step p {
			color: #555;
			line-height: 1.6;
			margin-bottom: 15px;
		}
		.status {
			display: inline-block;
			padding: 8px 15px;
			border-radius: 20px;
			font-size: 13px;
			font-weight: 600;
			margin-top: 10px;
		}
		.status.success {
			background: #d4edda;
			color: #155724;
		}
		.status.pending {
			background: #fff3cd;
			color: #856404;
		}
		.photo-grid {
			display: grid;
			grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
			gap: 20px;
			margin-top: 20px;
		}
		.photo-card {
			background: white;
			border-radius: 8px;
			overflow: hidden;
			box-shadow: 0 2px 8px rgba(0,0,0,0.1);
			transition: transform 0.3s;
		}
		.photo-card:hover {
			transform: translateY(-5px);
			box-shadow: 0 6px 20px rgba(0,0,0,0.15);
		}
		.photo-card img {
			width: 100%;
			height: 200px;
			object-fit: cover;
		}
		.photo-info {
			padding: 15px;
		}
		.photo-info h4 {
			color: #333;
			margin-bottom: 10px;
			font-size: 16px;
		}
		.photo-info button {
			width: 100%;
			background: #667eea;
			color: white;
			padding: 10px;
			font-size: 14px;
		}
		.logs {
			background: #2d3748;
			color: #e2e8f0;
			padding: 20px;
			border-radius: 8px;
			font-family: 'Courier New', monospace;
			font-size: 13px;
			max-height: 400px;
			overflow-y: auto;
			line-height: 1.6;
		}
		.log-entry {
			margin-bottom: 8px;
			padding: 4px 0;
		}
		.info-box {
			background: #e7f3ff;
			border: 1px solid #b3d9ff;
			padding: 15px;
			border-radius: 8px;
			margin-bottom: 20px;
		}
		.info-box strong {
			color: #004085;
		}
		code {
			background: #f4f4f4;
			padding: 2px 6px;
			border-radius: 3px;
			font-family: monospace;
			font-size: 13px;
		}
		.empty-state {
			text-align: center;
			padding: 60px 20px;
			color: #999;
		}
		.empty-state svg {
			width: 80px;
			height: 80px;
			margin-bottom: 20px;
			opacity: 0.3;
		}
	</style>
</head>
<body>
	<div class="container">
		<div class="card">
			<h1>🖨️ Printer Service</h1>
			<p class="subtitle">OAuth 1.0 Client Implementation (RFC 5849)</p>

			<div class="info-box">
				<strong>📋 Scenario:</strong> Jane wants to print her private vacation photos from her Photo Service account.
				The Printer Service needs OAuth authorization to access her photos.
			</div>

			{{if not .HasAccess}}
			<div class="hero">
				<h2>🖼️ Print Your Vacation Photos</h2>
				<p>Connect your Photo Service account to print your private photos</p>
				<form action="/start-oauth" method="POST">
					<button type="submit">🔐 Connect to Photo Service</button>
				</form>
			</div>
			{{end}}

			<div class="step">
				<h3>Step 1: Request Temporary Credentials</h3>
				<p>The printer requests temporary credentials from the Photo Service.</p>
				{{if .TempToken}}
					<div class="status success">✓ Temporary Token Received</div>
					<p style="margin-top: 10px;"><code>{{.TempToken}}</code></p>
				{{else}}
					<div class="status pending">⏳ Waiting...</div>
				{{end}}
			</div>

			<div class="step">
				<h3>Step 2: User Authorization</h3>
				<p>User is redirected to Photo Service to authorize access.</p>
				{{if .TempToken}}
					{{if .AccessToken}}
						<div class="status success">✓ Authorization Complete</div>
					{{else}}
						<div class="status pending">⏳ Redirecting to authorization...</div>
					{{end}}
				{{else}}
					<div class="status pending">⏳ Waiting...</div>
				{{end}}
			</div>

			<div class="step">
				<h3>Step 3: Exchange for Access Token</h3>
				<p>Exchange temporary credentials and verifier for access token.</p>
				{{if .AccessToken}}
					<div class="status success">✓ Access Token Received</div>
					<p style="margin-top: 10px;"><code>{{.AccessToken}}</code></p>
				{{else}}
					<div class="status pending">⏳ Waiting...</div>
				{{end}}
			</div>

			<div class="step">
				<h3>Step 4: Access Protected Resources</h3>
				<p>Use access token to fetch Jane's private photos.</p>
				{{if .AccessToken}}
					<form action="/fetch-photos" method="POST">
						<button type="submit">📸 Fetch Photos</button>
					</form>
				{{else}}
					<button disabled>📸 Fetch Photos</button>
				{{end}}
			</div>
		</div>

		{{if .Photos}}
		<div class="card">
			<h2 style="margin-bottom: 20px;">📸 Jane's Vacation Photos</h2>
			<div class="photo-grid">
				{{range .Photos}}
				<div class="photo-card">
					<img src="{{.URL}}" alt="{{.Title}}">
					<div class="photo-info">
						<h4>{{.Title}}</h4>
						<form action="/print" method="POST">
							<input type="hidden" name="photo_id" value="{{.ID}}">
							<input type="hidden" name="photo_title" value="{{.Title}}">
							<button type="submit">🖨️ Print This Photo</button>
						</form>
					</div>
				</div>
				{{end}}
			</div>
		</div>
		{{end}}

		{{if .Logs}}
		<div class="card">
			<h2 style="margin-bottom: 20px;">📋 OAuth Flow Logs</h2>
			<div class="logs">
				{{range .Logs}}
				<div class="log-entry">{{.}}</div>
				{{end}}
			</div>
		</div>
		{{end}}
	</div>
</body>
</html>
	`))

	tmpl.Execute(w, data)
}
