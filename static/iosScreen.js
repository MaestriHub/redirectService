//import copy from 'copy-to-clipboard';

const appStoreLink = "https://apps.apple.com/ru/app/youtube/id544007664";


(async function () {
    await getIosData();
    openApp();
})();


const openApp = () => {
    
    const copyButton = document.getElementById('copyButton');
};
copyButton.addEventListener('click', () => {
    navigator.clipboard.writeText(universalLink).then(() => {
        window.location.href = appStoreLink;
    }).catch(() => {
        alert("Не asdasdудалось скопировать ссылку.");
    });
});
window.addEventListener('load', () => {
    //copyToClipboard(universalLink);
    document.getElementById('copyButton').click();
    window.location.href = appStoreLink;
});
// function copyToClipboard(str) {
//     var el = document.createElement('textarea');  // Create a <textarea> element
//     el.value = str;                                 // Set its value to the string that you want copied
//     el.setAttribute('readonly', '');                // Make it readonly to be tamper-proof
//     el.style.position = 'absolute';                 
//     el.style.left = '-9999px';                      // Move outside the screen to make it invisible
//     document.body.appendChild(el);                  // Append the <textarea> element to the HTML document
//     var selected =            
//       document.getSelection().rangeCount > 0        // Check if there is any content selected previously
//         ? document.getSelection().getRangeAt(0)     // Store selection if found
//         : false;                                    // Mark as false to know no selection existed before
//     el.select();                                    // Select the <textarea> content
//     document.execCommand('copy');                   // Copy - only works as a result of a user action (e.g. click events)
//     document.body.removeChild(el);                  // Remove the <textarea> element
//     if (selected) {                                 // If a selection existed before copying
//       document.getSelection().removeAllRanges();    // Unselect everything on the HTML document
//       document.getSelection().addRange(selected);   // Restore the original selection
//     }
//   };
async function getIosData() {
    document.querySelector('#copyButton').addEventListener('click', async () => {
        try {
          await navigator.clipboard.writeText('Текст для копированияфыаЙУПЬ');
        } catch (err) {
          console.error('Ошибка при копировании: ', err);
        }
      });
    
        
    const clientData = {
        userAgent: navigator.userAgent,
        platform: navigator.platform || 'unknown',
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
    };

    await fetch('/IOS', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(clientData),
    });
}