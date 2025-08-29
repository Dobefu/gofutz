// @ts-check

/** @type {string} */
const WS_URL = `ws://${window.location.host}/ws`;

(() => {
  "use strict";

  /** @type {WebSocket | null} */
  let wss;

  window.addEventListener("beforeunload", () => {
    if (!wss || wss.readyState !== WebSocket.OPEN) {
      return;
    }

    wss.close();
    wss = null;
  });

  /**
   * @param {string} url
   */
  const initWebSocket = (url) => {
    // If a prior websocket still exists, clear all event listeners and close it.
    if (wss) {
      wss.onopen = null;
      wss.onmessage = null;
      wss.onclose = null;
      wss.onerror = null;
      wss.close();

      wss = null;
    }

    wss = new WebSocket(url);

    /**
     * @param {Event} _e
     */
    wss.onopen = (_e) => {};

    /**
     * @param {MessageEvent} _e
     */
    wss.onmessage = (_e) => {};

    /**
     * @param {CloseEvent} e
     */
    wss.onclose = (e) => {
      if (e.wasClean) {
        return;
      }

      setTimeout(() => {
        console.error("Websocket closed unexpectedly, reconnecting...");
        initWebSocket(url);
      }, 1000);
    };

    /**
     * @param {ErrorEvent} e
     */
    wss.onerror = (e) => {
      console.error("Websocket error", e);

      setTimeout(() => {
        console.error("Websocket closed unexpectedly, reconnecting...");
        initWebSocket(url);
      }, 1000);
    };
  };

  initWebSocket(WS_URL);
})();
