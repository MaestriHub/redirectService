
const appStoreLinkIOS = "https://apps.apple.com/app/maestri/id6469101735";

(async function () {
    await getData();
})();

async function getData() {
    var userAgent = navigator.userAgent
    const version = getIOSVersion();
    appStoreLink = appStoreLinkIOS;
    
    const canvas = document.getElementById("canvas");
    const gl = canvas.getContext("webgl");

    const debugInfo = gl.getExtension("WEBGL_debug_renderer_info");
    const vendor = gl.getParameter(debugInfo.UNMASKED_VENDOR_WEBGL);
    const renderer = gl.getParameter(debugInfo.UNMASKED_RENDERER_WEBGL);
    const clientData = {
        userAgent: userAgent,
        platform: navigator.platform || 'unknown',
        version: version,
        language: navigator.language,
        languages: navigator.languages,
        cores: navigator.hardwareConcurrency || 0,
        memory: navigator.deviceMemory || 0,
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

    await fetch('/collect/mobile', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(clientData),
    });

    document.getElementById('copyButton').addEventListener('click', () => {
        navigator.clipboard.writeText(universalLink).then(() => {
            window.location.href = appStoreLinkIOS;
        }).catch(() => {
            alert("Не asdasdудалось скопировать ссылку.");
        });
    });
}

function getIOSVersion() {
    var userAgent = navigator.userAgent;
    var match = userAgent.match(/(iPhone|iPad).*OS (\d+_\d+_\d+)/);

    if (match) {
        return match[2].replace(/_/g, '.'); 
    }
    return "unknown"; 
}