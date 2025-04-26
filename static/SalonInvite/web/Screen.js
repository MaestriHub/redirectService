function generateQRCode() {
    var qrcode = new QRCode(document.getElementById("qrcode"), {
        text: universalLink,
    });
}

async function getPCData() {
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
    await getPCData();
    generateQRCode()

    document.getElementById("copyButton").addEventListener("click", () => {
        const tempTextArea = document.createElement("textarea");
        tempTextArea.value = universalLink;
        tempTextArea.style.position = "absolute";
        tempTextArea.style.left = "-9999px";
        document.body.appendChild(tempTextArea);

        tempTextArea.select();
        tempTextArea.setSelectionRange(0, 99999);

        try {
            document.execCommand("copy");
        } catch (error) {
            console.error("Ошибка при копировании:", error);
        } finally {
            document.body.removeChild(tempTextArea);
        }
    });
})();
