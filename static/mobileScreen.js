//import copy from 'copy-to-clipboard';

const appStoreLinkIOS = "https://apps.apple.com/app/maestri/id6469101735";
const appStoreLinkAndroid = "https://play.google.com/apps/internaltest/4701369389039828090";
var appStoreLink = appStoreLinkIOS;

(async function () {
    await getData();
})();

async function getData() {
    
    var version = getIOSVersion();
    if (version == 'unknown') {
        version =  getAndroidVersion();
        appStoreLink = appStoreLinkAndroid
    }

    const clientData = {
        userAgent: navigator.userAgent,
        platform: navigator.platform || 'unknown',
        version: version,
        language: navigator.language,
        languages: navigator.languages,
        cookiesEnabled: navigator.cookieEnabled,
        connectionType: navigator.connection ? navigator.connection.effectiveType : 'unknown',
        isOnline: navigator.onLine,
        cores: navigator.hardwareConcurrency || 0,
        memory: navigator.deviceMemory || 0,
        screenWidth: screen.width,
        screenHeight: screen.height,
        colorDepth: screen.colorDepth,
        pixelRatio: window.devicePixelRatio,
        viewportWidth: window.innerWidth,
        viewportHeight: window.innerHeight,
        timeZone: Intl.DateTimeFormat().resolvedOptions().timeZone,
        currentTime: new Date().toISOString(),
		DirectURLID: linkId,
    };

    await fetch('/Mobile', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(clientData),
    });

    document.getElementById('copyButton').addEventListener('click', () => {
        navigator.clipboard.writeText(universalLink).then(() => {
            window.location.href = appStoreLink;
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

function getAndroidVersion() {
    var userAgent = navigator.userAgent;
    var match = userAgent.match(/Android (\d+\.\d+)/);

    if (match) {
        return match[1]; 
    }
    return "unknown"; 
}