# ============================================
# Create .env File Script
# ============================================

Write-Host "`n╔════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║   Create .env File                     ║" -ForegroundColor Cyan
Write-Host "╚════════════════════════════════════════╝`n" -ForegroundColor Cyan

$envPath = ".env"

# Check if .env already exists
if (Test-Path $envPath) {
    Write-Host "⚠️  .env file already exists!" -ForegroundColor Yellow
    $overwrite = Read-Host "Do you want to overwrite it? (yes/no)"
    if ($overwrite -ne "yes") {
        Write-Host "`n❌ Cancelled. Keeping existing .env file." -ForegroundColor Red
        exit 0
    }
}

# Create .env content
$envContent = @"
# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=50051
HTTP_HOST=0.0.0.0
HTTP_PORT=8080

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=iam_db
DB_SSL_MODE=disable

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT Configuration
JWT_SECRET=your-secret-key-change-this-in-production
JWT_EXPIRATION_HOURS=24
JWT_REFRESH_EXPIRATION_HOURS=168

# Log Configuration
LOG_LEVEL=info
LOG_ENCODING=json

# Swagger Configuration
SWAGGER_ENABLED=true
SWAGGER_BASE_PATH=/swagger/
SWAGGER_SPEC_PATH=/swagger.json
SWAGGER_TITLE=IAM Service API Documentation
SWAGGER_AUTH_USERNAME=admin
SWAGGER_AUTH_PASSWORD=changeme
SWAGGER_AUTH_REALM=IAM Service API Documentation
"@

# Write to file
Write-Host "📝 Creating .env file..." -ForegroundColor Cyan
$envContent | Out-File -FilePath $envPath -Encoding ASCII -NoNewline

# Verify
if (Test-Path $envPath) {
    Write-Host "✅ .env file created successfully!`n" -ForegroundColor Green
    
    Write-Host "📄 Content:" -ForegroundColor Cyan
    Write-Host "─────────────────────────────────────────" -ForegroundColor Gray
    Get-Content $envPath | ForEach-Object {
        if ($_ -match "^#") {
            Write-Host $_ -ForegroundColor Yellow
        } elseif ($_ -match "=") {
            $parts = $_ -split "=", 2
            Write-Host "$($parts[0])=" -NoNewline -ForegroundColor White
            Write-Host $parts[1] -ForegroundColor Green
        } else {
            Write-Host $_
        }
    }
    Write-Host "─────────────────────────────────────────`n" -ForegroundColor Gray
    
    Write-Host "✅ You can now start the service with:" -ForegroundColor Green
    Write-Host "   go run cmd/server/main.go`n" -ForegroundColor Yellow
    
} else {
    Write-Host "❌ Failed to create .env file!" -ForegroundColor Red
    exit 1
}

