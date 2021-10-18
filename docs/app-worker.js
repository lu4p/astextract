const cacheName = "app-" + "5e127a19e718aafbeb227fcf4ad575f36e3e8316";

self.addEventListener("install", event => {
  console.log("installing app worker 5e127a19e718aafbeb227fcf4ad575f36e3e8316");

  event.waitUntil(
    caches.open(cacheName).
      then(cache => {
        return cache.addAll([
          "/astextract",
          "/astextract/app.css",
          "/astextract/app.js",
          "/astextract/manifest.webmanifest",
          "/astextract/wasm_exec.js",
          "/astextract/web/app.wasm",
          "https://storage.googleapis.com/murlok-github/icon-192.png",
          "https://storage.googleapis.com/murlok-github/icon-512.png",
          "https://unpkg.com/spectre.css/dist/spectre.min.css",
          
        ]);
      }).
      then(() => {
        self.skipWaiting();
      })
  );
});

self.addEventListener("activate", event => {
  event.waitUntil(
    caches.keys().then(keyList => {
      return Promise.all(
        keyList.map(key => {
          if (key !== cacheName) {
            return caches.delete(key);
          }
        })
      );
    })
  );
  console.log("app worker 5e127a19e718aafbeb227fcf4ad575f36e3e8316 is activated");
});

self.addEventListener("fetch", event => {
  event.respondWith(
    caches.match(event.request).then(response => {
      return response || fetch(event.request);
    })
  );
});
