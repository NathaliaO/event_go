# Script PowerShell simples para Windows
# Execute com: .\build.ps1 [comando]

param([string]$Command = "clean")

$DockerCompose = "docker compose"

switch ($Command.ToLower()) {
    "build" {
        Write-Host "Construindo aplicação..." -ForegroundColor Yellow
        & $DockerCompose build
    }
    
    "run" {
        Write-Host "Iniciando aplicação..." -ForegroundColor Green
        & $DockerCompose up -d
    }
    
    "stop" {
        Write-Host "Parando aplicação..." -ForegroundColor Red
        & $DockerCompose down
    }
    
    "logs" {
        Write-Host "Mostrando logs..." -ForegroundColor Cyan
        & $DockerCompose logs -f
    }
    
    "dev" {
        Write-Host "Modo desenvolvimento (logs em tempo real)..." -ForegroundColor Yellow
        & $DockerCompose up
    }
    
    default {
        Write-Host "Event Go - Scripts para Windows" -ForegroundColor Green
        Write-Host "Comandos: build, run, stop, logs, dev" -ForegroundColor Cyan
        Write-Host "Exemplo: .\build.ps1 build" -ForegroundColor Yellow
    }
}