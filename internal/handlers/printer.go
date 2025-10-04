package handlers

import (
	"fmt"
	"net/http"
)

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<title>OAuth 1.0 Example - Jane's Photo Service</title>
			<script src="https://unpkg.com/vue@3/dist/vue.global.js"></script>
			<style>
				* { margin: 0; padding: 0; box-sizing: border-box; }
				body {
					font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
					background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
					min-height: 100vh;
					padding: 20px;
				}
				#app {
					max-width: 900px;
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
					font-size: 28px;
				}
				h2 {
					color: #555;
					margin-bottom: 20px;
					font-size: 20px;
					border-bottom: 2px solid #667eea;
					padding-bottom: 10px;
				}
				.subtitle {
					color: #666;
					margin-bottom: 30px;
					font-size: 14px;
				}
				.step {
					background: #f8f9fa;
					padding: 20px;
					border-radius: 8px;
					margin-bottom: 15px;
					border-left: 4px solid #667eea;
				}
				.step h3 {
					color: #667eea;
					margin-bottom: 10px;
					font-size: 16px;
				}
				button {
					background: #667eea;
					color: white;
					border: none;
					padding: 12px 24px;
					border-radius: 6px;
					cursor: pointer;
					font-size: 14px;
					font-weight: 600;
					transition: all 0.3s;
					margin-right: 10px;
					margin-top: 10px;
				}
				button:hover {
					background: #5568d3;
					transform: translateY(-2px);
					box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
				}
				button:disabled {
					background: #ccc;
					cursor: not-allowed;
					transform: none;
				}
				.status {
					padding: 10px 15px;
					border-radius: 6px;
					margin-top: 15px;
					font-size: 14px;
				}
				.status.success {
					background: #d4edda;
					color: #155724;
					border: 1px solid #c3e6cb;
				}
				.status.error {
					background: #f8d7da;
					color: #721c24;
					border: 1px solid #f5c6cb;
				}
				.status.info {
					background: #d1ecf1;
					color: #0c5460;
					border: 1px solid #bee5eb;
				}
				pre {
					background: #2d3748;
					color: #e2e8f0;
					padding: 15px;
					border-radius: 6px;
					overflow-x: auto;
					font-size: 12px;
					margin-top: 10px;
				}
				.photo-grid {
					display: grid;
					grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
					gap: 15px;
					margin-top: 20px;
				}
				.photo-item {
					background: #f8f9fa;
					padding: 15px;
					border-radius: 8px;
					text-align: center;
				}
				.photo-item img {
					width: 100%%;
					height: 150px;
					object-fit: cover;
					border-radius: 6px;
					margin-bottom: 10px;
				}
				input[type="text"], input[type="password"] {
					width: 100%%;
					padding: 10px;
					border: 1px solid #ddd;
					border-radius: 6px;
					font-size: 14px;
					margin-bottom: 10px;
				}
				.credentials {
					background: #fff3cd;
					border: 1px solid #ffc107;
					padding: 15px;
					border-radius: 6px;
					margin-bottom: 20px;
				}
				.credentials code {
					background: #fff;
					padding: 2px 6px;
					border-radius: 3px;
					font-family: monospace;
					font-size: 12px;
				}
			</style>
		</head>
		<body>
			<div id="app">
				<div class="card">
					<h1>OAuth 1.0 Example</h1>
					<p class="subtitle">Jane prints her private vacation photos using OAuth 1.0 (RFC 5849)</p>
					
					<div class="credentials">
						<strong>Client Credentials (Printer Service):</strong><br>
						Client Key: <code>dpf43f3p2l4k3l03</code><br>
						Client Secret: <code>kd94hf93k423kf44</code>
					</div>

					<div class="step">
						<h3>Step 1: Get Temporary Credentials</h3>
						<p>The printer service requests temporary credentials from Jane's photo service.</p>
						<button @click="getTemporaryCredentials" :disabled="step > 1">Request Temporary Credentials</button>
						<div v-if="tempToken" class="status info">
							<strong>Temporary Token:</strong> {{ tempToken }}<br>
							<strong>Temporary Secret:</strong> {{ tempSecret }}
						</div>
					</div>

					<div class="step" v-if="step >= 2">
						<h3>Step 2: User Authorization</h3>
						<p>Jane needs to authorize the printer service to access her photos.</p>
						<button @click="authorizeToken" :disabled="step > 2">Go to Authorization Page</button>
						
						<div v-if="showAuthForm" style="margin-top: 15px;">
							<h4 style="margin-bottom: 10px;">Authorization Required</h4>
							<p style="margin-bottom: 10px;">The printer service wants to access your photos. Please log in to authorize.</p>
							<input v-model="username" type="text" placeholder="Username (e.g., jane)">
							<input v-model="password" type="password" placeholder="Password">
							<button @click="submitAuthorization">Authorize Application</button>
							<button @click="denyAuthorization" style="background: #dc3545;">Deny</button>
						</div>

						<div v-if="verifier" class="status success">
							<strong>Authorization Successful!</strong><br>
							Verifier: {{ verifier }}
						</div>
					</div>

					<div class="step" v-if="step >= 3">
						<h3>Step 3: Exchange for Access Token</h3>
						<p>Exchange the temporary credentials and verifier for a long-term access token.</p>
						<button @click="getAccessToken" :disabled="step > 3">Get Access Token</button>
						<div v-if="accessToken" class="status success">
							<strong>Access Token:</strong> {{ accessToken }}<br>
							<strong>Access Secret:</strong> {{ accessSecret }}
						</div>
					</div>

					<div class="step" v-if="step >= 4">
						<h3>Step 4: Access Protected Resources</h3>
						<p>Now the printer can access Jane's private vacation photos.</p>
						<button @click="fetchPhotos">Fetch Jane's Photos</button>
						
						<div v-if="photos.length > 0" class="photo-grid">
							<div v-for="photo in photos" :key="photo.id" class="photo-item">
								<img :src="photo.url" :alt="photo.title">
								<div><strong>{{ photo.title }}</strong></div>
								<button @click="printPhoto(photo)">🖨️ Print</button>
							</div>
						</div>
					</div>

					<div v-if="error" class="status error">
						{{ error }}
					</div>

					<div v-if="logs.length > 0" style="margin-top: 30px;">
						<h2>Request/Response Log</h2>
						<pre>{{ logs.join('\n\n') }}</pre>
					</div>
				</div>
			</div>

			<script>
				const { createApp } = Vue;

				createApp({
					data() {
						return {
							step: 1,
							tempToken: '',
							tempSecret: '',
							verifier: '',
							accessToken: '',
							accessSecret: '',
							username: 'jane',
							password: 'password',
							showAuthForm: false,
							photos: [],
							error: '',
							logs: []
						}
					},
					methods: {
						async getTemporaryCredentials() {
							this.error = '';
							this.logs.push('=== STEP 1: REQUEST TEMPORARY CREDENTIALS ===');
							
							try {
								const response = await fetch('/initiate', {
									method: 'POST',
									headers: { 'Content-Type': 'application/json' },
									body: JSON.stringify({
										oauth_callback: 'http://localhost:8080/callback'
									})
								});
								
								const data = await response.json();
								this.logs.push('Response: ' + JSON.stringify(data, null, 2));
								
								if (data.oauth_token) {
									this.tempToken = data.oauth_token;
									this.tempSecret = data.oauth_token_secret;
									this.step = 2;
								} else {
									this.error = 'Failed to get temporary credentials';
								}
							} catch (e) {
								this.error = e.message;
							}
						},
						
						authorizeToken() {
							this.showAuthForm = true;
							this.logs.push('=== STEP 2: USER AUTHORIZATION ===');
							this.logs.push('Opening authorization form...');
						},
						
						async submitAuthorization() {
							try {
								const response = await fetch('/authorize-submit', {
									method: 'POST',
									headers: { 'Content-Type': 'application/json' },
									body: JSON.stringify({
										oauth_token: this.tempToken,
										username: this.username,
										password: this.password,
										authorize: true
									})
								});
								
								const data = await response.json();
								this.logs.push('Authorization Response: ' + JSON.stringify(data, null, 2));
								
								if (data.oauth_verifier) {
									this.verifier = data.oauth_verifier;
									this.showAuthForm = false;
									this.step = 3;
								} else {
									this.error = data.error || 'Authorization failed';
								}
							} catch (e) {
								this.error = e.message;
							}
						},
						
						denyAuthorization() {
							this.showAuthForm = false;
							this.error = 'Authorization denied by user';
						},
						
						async getAccessToken() {
							this.error = '';
							this.logs.push('=== STEP 3: EXCHANGE FOR ACCESS TOKEN ===');
							
							try {
								const response = await fetch('/token', {
									method: 'POST',
									headers: { 'Content-Type': 'application/json' },
									body: JSON.stringify({
										oauth_token: this.tempToken,
										oauth_verifier: this.verifier
									})
								});
								
								const data = await response.json();
								this.logs.push('Response: ' + JSON.stringify(data, null, 2));
								
								if (data.oauth_token) {
									this.accessToken = data.oauth_token;
									this.accessSecret = data.oauth_token_secret;
									this.step = 4;
								} else {
									this.error = 'Failed to get access token';
								}
							} catch (e) {
								this.error = e.message;
							}
						},
						
						async fetchPhotos() {
							this.error = '';
							this.logs.push('=== STEP 4: ACCESS PROTECTED RESOURCES ===');
							
							try {
								const response = await fetch('/photos', {
									headers: {
										'Authorization': 'Bearer ' + this.accessToken
									}
								});
								
								const data = await response.json();
								this.logs.push('Response: ' + JSON.stringify(data, null, 2));
								
								if (data.photos) {
									this.photos = data.photos;
								} else {
									this.error = 'Failed to fetch photos';
								}
							} catch (e) {
								this.error = e.message;
							}
						},
						
						printPhoto(photo) {
							alert('Printing: ' + photo.title);
							this.logs.push('🖨️ PRINTED: ' + photo.title);
						}
					}
				}).mount('#app');
			</script>
		</body>
		</html>
	`)
}
