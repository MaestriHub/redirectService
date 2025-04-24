(async function () {
    await getData();
})();

async function getData() {
    const canvas = document.getElementById("canvas");
    const gl = canvas.getContext("webgl");

    const debugInfo = gl.getExtension("WEBGL_debug_renderer_info");
    const vendor = gl.getParameter(debugInfo.UNMASKED_VENDOR_WEBGL);
    const renderer = gl.getParameter(debugInfo.UNMASKED_RENDERER_WEBGL);
    const clientData = {
        userAgent: navigator.userAgent,
        language: navigator.language,
        languages: navigator.languages,
        cores: navigator.hardwareConcurrency || 0,
        screenWidth: screen.width,
        screenHeight: screen.height,
        colorDepth: screen.colorDepth,
        pixelRatio: window.devicePixelRatio,
        viewportWidth: window.innerWidth,
        viewportHeight: window.innerHeight,
        renderer: renderer,
        vendorRender: vendor,
        timeZone: Intl.DateTimeFormat().resolvedOptions().timeZone,
		DirectLinkID: linkId,
    };

    await fetch('/fingerprint/' + linkId, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(clientData),
    });

    document.getElementById('copyButton').addEventListener('click', () => {
        window.location.href = appStoreLink;
    });
}