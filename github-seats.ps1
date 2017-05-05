Param(
    [string]$HostAddress,
    [string]$Token,
    [string]$Ssl
)

if($ssl){
    $prefix = "https://"
} else {
    $prefix = "http://"
}

$fqdn = $prefix + $HostAddress
$header = @{"Authorization" = "token $Token"}

$res = iwr -Uri $fqdn/api/v3/enterprise/settings/license -Headers $header -UseBasicParsing|ConvertFrom-Json

Write-Host @"
<?xml version="1.0" encoding="UTF-8" ?>
<prtg>
    <result>
        <channel>Seats Available</channel>
        <value>$($res.seats_available)</value>
    </result>
    <result>
        <channel>Seats Licensed</channel>
        <value>$($res.seats)</value>
    </result>
    <result>
        <channel>Seats Used</channel>
        <value>$($res.seats_used)</value>
    </result>
    <result>
        <channel>Days Until Expiration</channel>
        <value>$($res.days_until_expiration)</value>
    </result>
    <text>$($res.seats_available) seats remaining</text>
</prtg>
"@
