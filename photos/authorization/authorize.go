package authorization

import (
	"fmt"
	"html/template"
	"net/http"
)

func HandleAuthorize(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("oauth_token")

	fmt.Printf("[AUTHORIZE] User authorization request for token: %s\n", token)

	store.mu.RLock()
	_, exists := store.tempCredentials[token]
	store.mu.RUnlock()

	if !exists {
		http.Error(w, "Invalid or expired token", http.StatusBadRequest)
		fmt.Printf("[AUTHORIZE] Token not found: %s\n", token)
		return
	}

	tmpl := template.Must(template.New("authorize").Parse(`
<!DOCTYPE html>
<html>
<head>
	<title>Authorize Printer Service</title>
	<style>
		body {
			font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
			background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
			min-height: 100vh;
			display: flex;
			align-items: center;
			justify-content: center;
			padding: 20px;
		}
		.auth-card {
			background: white;
			border-radius: 12px;
			padding: 40px;
			box-shadow: 0 10px 30px rgba(0,0,0,0.3);
			max-width: 450px;
			width: 100%;
		}
		h1 { color: #333; margin-bottom: 10px; font-size: 24px; }
		.subtitle { color: #666; margin-bottom: 30px; font-size: 14px; }
		.warning {
			background: #fff3cd;
			border: 1px solid #ffc107;
			padding: 15px;
			border-radius: 8px;
			margin-bottom: 25px;
			font-size: 14px;
		}
		.form-group {
			margin-bottom: 20px;
		}
		label {
			display: block;
			margin-bottom: 5px;
			font-weight: 600;
			color: #333;
			font-size: 14px;
		}
		input[type="text"], input[type="password"] {
			width: 100%;
			padding: 12px;
			border: 1px solid #ddd;
			border-radius: 6px;
			font-size: 14px;
			box-sizing: border-box;
		}
		input[type="text"]:focus, input[type="password"]:focus {
			outline: none;
			border-color: #667eea;
			box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
		}
		.button-group {
			display: flex;
			gap: 10px;
			margin-top: 25px;
		}
		button {
			flex: 1;
			padding: 12px;
			border: none;
			border-radius: 6px;
			font-size: 14px;
			font-weight: 600;
			cursor: pointer;
			transition: all 0.3s;
		}
		.btn-authorize {
			background: #667eea;
			color: white;
		}
		.btn-authorize:hover {
			background: #5568d3;
			transform: translateY(-2px);
			box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
		}
		.btn-deny {
			background: #dc3545;
			color: white;
		}
		.btn-deny:hover {
			background: #c82333;
		}
		.hint {
			background: #e7f3ff;
			border: 1px solid #b3d9ff;
			padding: 10px;
			border-radius: 6px;
			font-size: 12px;
			color: #004085;
			margin-top: 15px;
		}
	</style>
</head>
<body>
	<div class="auth-card">
		<h1>🔐 Authorization Required</h1>
		<p class="subtitle">Jane's Photo Service</p>
		
		<div class="warning">
			<strong>⚠️ Printer Service</strong> wants to access your private vacation photos.
		</div>

		<form method="POST" action="/oauth/authorize-submit">
			<input type="hidden" name="oauth_token" value="{{.Token}}">
			
			<div class="form-group">
				<label for="username">Username</label>
				<input type="text" id="username" name="username" placeholder="Enter your username" required autofocus>
			</div>
			
			<div class="form-group">
				<label for="password">Password</label>
				<input type="password" id="password" name="password" placeholder="Enter your password" required>
			</div>

			<div class="button-group">
				<button type="submit" name="authorize" value="true" class="btn-authorize">✓ Authorize</button>
				<button type="submit" name="authorize" value="false" class="btn-deny">✗ Deny</button>
			</div>

			<div class="hint">
				💡 Test credentials: username=<strong>jane</strong>, password=<strong>password123</strong>
			</div>
		</form>
	</div>
</body>
</html>
	`))

	tmpl.Execute(w, map[string]string{"Token": token})
}
