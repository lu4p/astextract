// -----------------------------------------------------------------------------
// Init service worker
// -----------------------------------------------------------------------------
if ("serviceWorker" in navigator) {
  navigator.serviceWorker
    .register("/astextract/app-worker.js")
    .then(reg => {
      console.log("registering app service worker");
    })
    .catch(err => {
      console.error("offline service worker registration failed", err);
    });
}

// -----------------------------------------------------------------------------
// Init progressive app
// -----------------------------------------------------------------------------
const goappEnv = {"GOAPP_ROOT_PREFIX":"/astextract","GOAPP_STATIC_RESOURCES_URL":"/astextract","GOAPP_VERSION":"91e7989a230f05103fe4e93dc2e57405536cff69"};

function goappGetenv(k) {
  return goappEnv[k];
}

let deferredPrompt;

window.addEventListener("beforeinstallprompt", e => {
  e.preventDefault();
  deferredPrompt = e;
});

// -----------------------------------------------------------------------------
// Init Web Assembly
// -----------------------------------------------------------------------------
if (!WebAssembly.instantiateStreaming) {
  WebAssembly.instantiateStreaming = async (resp, importObject) => {
    const source = await (await resp).arrayBuffer();
    return await WebAssembly.instantiate(source, importObject);
  };
}

const go = new Go();

WebAssembly.instantiateStreaming(fetch("/astextract/web/app.wasm"), go.importObject)
  .then(result => {
    const loaderIcon = document.getElementById("app-wasm-loader-icon");
    loaderIcon.className = "goapp-logo";

    go.run(result.instance);
  })
  .catch(err => {
    const loaderIcon = document.getElementById("app-wasm-loader-icon");
    loaderIcon.className = "goapp-logo";

    const loaderLabel = document.getElementById("app-wasm-loader-label");
    loaderLabel.innerText = err;

    console.error("loading wasm failed: " + err);
  });

// -----------------------------------------------------------------------------
// Keep body clean
// -----------------------------------------------------------------------------
function goappKeepBodyClean() {
  const body = document.body;
  const bodyChildrenCount = body.children.length;

  const mutationObserver = new MutationObserver(function (mutationList) {
    mutationList.forEach((mutation) => {
      switch (mutation.type) {
        case 'childList':
          while (body.children.length > bodyChildrenCount) {
            body.removeChild(body.lastChild);
          }
          break;
      }
    });
  });

  mutationObserver.observe(document.body, {
    childList: true,
  });

  return () => mutationObserver.disconnect();
}
