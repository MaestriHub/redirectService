async function getData() {
    const clientData = {
        userAgent: navigator.userAgent,
        language: navigator.language,
        languages: navigator.languages,
        cookiesEnabled: navigator.cookieEnabled,
        connectionType: navigator.connection,
        isOnline: navigator.onLine,
        cores: navigator.hardwareConcurrency || 0,
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

    await fetch('/fingerprint/' + linkId, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(clientData),
    });
}

(async function () {
    await getData();

    document.getElementById('copyButton').addEventListener('click', () => {
        window.location.href = appStoreLink;
    });
})();
