async function getData() {
    const clientData = {
        language: navigator.language,
        languages: navigator.languages,
        cores: navigator.hardwareConcurrency || 0,
        screenWidth: screen.width,
        screenHeight: screen.height,
        colorDepth: screen.colorDepth,
        pixelRatio: window.devicePixelRatio,
        viewportWidth: window.innerWidth,
        viewportHeight: window.innerHeight,
        timeZone: Intl.DateTimeFormat().resolvedOptions().timeZone,
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

    document.getElementById("copyButton").addEventListener("click", () => {
        try {
            await navigator.clipboard.writeText(universalLink);
            console.log("Текст успешно скопирован в буфер обмена:", universalLink);
            window.location.href = appStoreLink;
        } catch (error) {
            console.error("Ошибка при копировании:", error);
        }
    });
})();