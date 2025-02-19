
(async function () {
    await getPCData();
    generateQRCode()
    document.getElementById('copyButton').addEventListener('click', () => {
        navigator.clipboard.writeText(universalLink).then(() => {
        }).catch(() => {
            alert("Не asdasdудалось скопировать ссылку.");
        });
    });
})();
function generateQRCode() {
    var qrcode = new QRCode(document.getElementById("qrcode"), {
        text: universalLink,
        width: 100,
        height: 100,
        colorDark : "#000000",
        colorLight : "#ffffff",
        correctLevel : QRCode.CorrectLevel.H
    });
}
async function getPCData() {
    const clientData = {
        userAgent: navigator.userAgent,
        platform: navigator.userAgentData.platform || 'unknown',
        language: navigator.language,
        languages: navigator.languages,
        cookiesEnabled: navigator.cookieEnabled,
        connectionType: navigator.connection ? navigator.connection.effectiveType : 'unknown',
        isOnline: navigator.onLine,
        cores: navigator.hardwareConcurrency,
        memory: navigator.deviceMemory || 0,
        screenWidth: screen.width,
        screenHeight: screen.height,
        colorDepth: screen.colorDepth,
        pixelRatio: window.devicePixelRatio,
        viewportWidth: window.innerWidth,
        viewportHeight: window.innerHeight,
        timeZone: Intl.DateTimeFormat().resolvedOptions().timeZone,
        currentTime: new Date().toISOString(),
        DirectLinkID: linkId,
    };

    await fetch('/collect/pc', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(clientData),
    });
}
