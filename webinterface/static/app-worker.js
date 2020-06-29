const cacheName = "app-" + "dcf1a24a37101967a97d885ea128b33d9d54dafb";

self.addEventListener("install", event => {
  console.log("installing app worker dcf1a24a37101967a97d885ea128b33d9d54dafb");
  self.skipWaiting();

  event.waitUntil(
    caches.open(cacheName).then(cache => {
      return cache.addAll([
        "/astextract",
        "/astextract/app.css",
        "/astextract/app.js",
        "/astextract/manifest.json",
        "/astextract/wasm_exec.js",
        "/astextract/web/app.wasm",
        "https://storage.googleapis.com/murlok-github/icon-192.png",
        "https://storage.googleapis.com/murlok-github/icon-512.png",
        "https://unpkg.com/spectre.css/dist/spectre.min.css",
        
      ]);
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
  console.log("app worker dcf1a24a37101967a97d885ea128b33d9d54dafb is activated");
});

self.addEventListener("fetch", event => {
  event.respondWith(
    caches.match(event.request).then(response => {
      return response || fetch(event.request);
    })
  );
});
