# Test script for Load Balancer
Write-Host "Starting Load Balancer Tests..." -ForegroundColor Green

Write-Host "`nTest 1: No Authentication" -ForegroundColor Yellow
Write-Host "Expected: 401 Unauthorized"
try {
    Invoke-WebRequest -Uri "http://localhost:8090/loadBalancer" -Method GET -ErrorAction Stop
} catch {
    Write-Host "Response: $($_.Exception.Response.StatusCode.Value__) $($_.Exception.Response.StatusDescription)"
}

Write-Host "`nTest 2: Valid Authentication" -ForegroundColor Yellow
Write-Host "Expected: 200 OK"
try {
    $headers = @{
        "X-API-Key" = "test-key"
    }
    $response = Invoke-WebRequest -Uri "http://localhost:8090/loadBalancer" -Headers $headers -Method GET -ErrorAction Stop
    Write-Host "Response: $($response.StatusCode) $($response.StatusDescription)"
} catch {
    Write-Host "Response: $($_.Exception.Response.StatusCode.Value__) $($_.Exception.Response.StatusDescription)"
}

Write-Host "`nTest 3: Invalid Authentication" -ForegroundColor Yellow
Write-Host "Expected: 401 Unauthorized"
try {
    $headers = @{
        "X-API-Key" = "wrong-key"
    }
    Invoke-WebRequest -Uri "http://localhost:8090/loadBalancer" -Headers $headers -Method GET -ErrorAction Stop
} catch {
    Write-Host "Response: $($_.Exception.Response.StatusCode.Value__) $($_.Exception.Response.StatusDescription)"
}

Write-Host "`nTest 4: Round Robin (10 requests)" -ForegroundColor Yellow
Write-Host "Expected: Different server responses in round-robin order"
$headers = @{
    "X-API-Key" = "test-key"
}

$responses = @{}
for ($i = 1; $i -le 10; $i++) {
    Write-Host "`nRequest ${i}"
    try {
        $response = Invoke-WebRequest -Uri "http://localhost:8090/loadBalancer" -Headers $headers -Method GET -ErrorAction Stop
        Write-Host "Response: $($response.Content)"
        
        # Track response patterns to verify round-robin
        if ($responses.ContainsKey($response.Content)) {
            $responses[$response.Content]++
        } else {
            $responses[$response.Content] = 1
        }
        
        # Small delay to make output readable
        Start-Sleep -Milliseconds 100
    } catch {
        Write-Host "Error: $($_.Exception.Message)"
    }
}

# Verify round-robin distribution
Write-Host "`nResponse Distribution:" -ForegroundColor Yellow
$responses.GetEnumerator() | ForEach-Object {
    Write-Host "$($_.Key): $($_.Value) times"
}

Write-Host "`nTests Complete!" -ForegroundColor Green