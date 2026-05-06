[CmdletBinding()]
param(
    [string]$HandoffsRoot = (Join-Path $PSScriptRoot "..\handoffs"),
    [string]$Prefix = "REQ"
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

function Resolve-FullPath {
    param(
        [Parameter(Mandatory = $true)]
        [string]$Path
    )

    return [System.IO.Path]::GetFullPath($Path)
}

$handoffsRootPath = Resolve-FullPath -Path $HandoffsRoot
if (-not (Test-Path -LiteralPath $handoffsRootPath)) {
    New-Item -ItemType Directory -Path $handoffsRootPath -Force | Out-Null
}

$maxAttempts = 20
$escapedPrefix = [Regex]::Escape($Prefix)

for ($attempt = 1; $attempt -le $maxAttempts; $attempt++) {
    $timestamp = Get-Date -Format "yyyyMMdd-HHmmss"
    $pattern = "^{0}-{1}-(\d{{3}})$" -f $escapedPrefix, [Regex]::Escape($timestamp)

    $existingSequenceNumbers = @(
        Get-ChildItem -LiteralPath $handoffsRootPath -Directory -ErrorAction SilentlyContinue |
            ForEach-Object {
                $match = [Regex]::Match($_.Name, $pattern)
                if ($match.Success) {
                    [int]$match.Groups[1].Value
                }
            }
    )

    $sequence = if ($existingSequenceNumbers.Count -gt 0) {
        (($existingSequenceNumbers | Measure-Object -Maximum).Maximum) + 1
    }
    else {
        1
    }

    if ($sequence -gt 999) {
        throw "同一秒内の request_id 連番が 999 を超えたため採番できませんでした。1 秒待って再実行してください。"
    }

    $requestId = "{0}-{1}-{2:000}" -f $Prefix, $timestamp, $sequence
    $handoffDirectory = Join-Path $handoffsRootPath $requestId

    try {
        New-Item -ItemType Directory -Path $handoffDirectory -ErrorAction Stop | Out-Null

        $result = [PSCustomObject]@{
            request_id        = $requestId
            handoff_directory = $handoffDirectory
        }

        $result | ConvertTo-Json -Compress

        return
    }
    catch {
        if (Test-Path -LiteralPath $handoffDirectory) {
            continue
        }

        throw
    }
}

throw "request_id の採番に繰り返し失敗しました。時間を置いて再実行してください。"