[CmdletBinding()]
param(
    [Parameter(Mandatory = $true)]
    [string]$RequestId,

    [Parameter(Mandatory = $true)]
    [Alias("AgentName", "AgentSlug")]
    [string]$Agent,

    [string]$HandoffsRoot = (Join-Path $PSScriptRoot "..\handoffs"),
    [string]$RelativeHandoffsRoot = ".github/agents/handoffs"
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

function ConvertTo-AgentSlug {
    param(
        [Parameter(Mandatory = $true)]
        [string]$Name
    )

    $slug = [Regex]::Replace($Name.ToLowerInvariant(), "[^a-z0-9]+", "-").Trim('-')
    if ([string]::IsNullOrWhiteSpace($slug)) {
        throw "エージェント名 '$Name' から有効な slug を生成できませんでした。英数字を含む名前を指定してください。"
    }

    return $slug
}

function Join-ForwardSlashPath {
    param(
        [Parameter(Mandatory = $true)]
        [string[]]$Segments
    )

    return (($Segments | Where-Object { -not [string]::IsNullOrWhiteSpace($_) }) -join "/") -replace "\\", "/"
}

$handoffsRootPath = Resolve-FullPath -Path $HandoffsRoot
$handoffDirectory = Join-Path $handoffsRootPath $RequestId
if (-not (Test-Path -LiteralPath $handoffDirectory -PathType Container)) {
    throw "request_id '$RequestId' に対応する handoff ディレクトリが見つかりません。先に .github/agents/scripts/new-request-id.ps1 を実行してください。"
}

$agentSlug = ConvertTo-AgentSlug -Name $Agent
$relativeHandoffsRootNormalized = (($RelativeHandoffsRoot -replace "\\", "/").TrimEnd('/'))
$maxAttempts = 20

for ($attempt = 1; $attempt -le $maxAttempts; $attempt++) {
    $timestamp = Get-Date -Format "yyyyMMdd-HHmmss-fff"
    $baseFileName = "{0}-{1}" -f $timestamp, $agentSlug
    $requestPath = Join-Path $handoffDirectory ("{0}.request.yaml" -f $baseFileName)
    $responsePath = Join-Path $handoffDirectory ("{0}.response.yaml" -f $baseFileName)
    $requestCreated = $false

    try {
        New-Item -ItemType File -Path $requestPath -ErrorAction Stop | Out-Null
        $requestCreated = $true
        New-Item -ItemType File -Path $responsePath -ErrorAction Stop | Out-Null

        $requestFile = Join-ForwardSlashPath -Segments @($relativeHandoffsRootNormalized, $RequestId, ("{0}.request.yaml" -f $baseFileName))
        $responseFile = Join-ForwardSlashPath -Segments @($relativeHandoffsRootNormalized, $RequestId, ("{0}.response.yaml" -f $baseFileName))

        $result = [PSCustomObject]@{
            request_id        = $RequestId
            handoff_directory = $handoffDirectory
            agent             = $Agent
            agent_slug        = $agentSlug
            timestamp         = $timestamp
            request_file      = $requestFile
            response_file     = $responseFile
            request_path      = $requestPath
            response_path     = $responsePath
        }

        $result | ConvertTo-Json -Compress

        return
    }
    catch {
        if ($requestCreated -and (Test-Path -LiteralPath $requestPath -PathType Leaf)) {
            Remove-Item -LiteralPath $requestPath -ErrorAction SilentlyContinue
        }

        if ((Test-Path -LiteralPath $requestPath -PathType Leaf) -or (Test-Path -LiteralPath $responsePath -PathType Leaf)) {
            continue
        }

        throw
    }
}

throw "handoff ファイル名の採番に繰り返し失敗しました。時間を置いて再実行してください。"
