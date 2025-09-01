"use strict";

/** @type {string} */
const WS_URL = `ws://${window.location.host}/ws`;

/** @type {WebSocket | null} */
let wss;
/** @type {AbortController | null} */
let abortController;

/**
 * @param {string} url
 */
function initWebSocket(url) {
  // If a websocket still exists, clear all event listeners and close it.
  if (wss) {
    wss.onopen = null;
    wss.onmessage = null;
    wss.onclose = null;
    wss.onerror = null;
    wss.close();

    wss = null;
  }

  if (abortController) {
    abortController.abort();
    abortController = null;
  }

  abortController = new AbortController();
  wss = new WebSocket(url);

  wss.onopen = () => {
    if (!wss) {
      return;
    }

    wss.send(JSON.stringify({ method: "gofutz:init" }));
  };

  /**
   * @param {MessageEvent} e
   */
  wss.onmessage = (e) => {
    /** @type {Message} */
    const msg = JSON.parse(e.data);

    switch (msg.method) {
      case "gofutz:init":
        window.dispatchEvent(
          new CustomEvent("gofutz:init", { detail: msg }),
        );
        break;

      default:
        console.log({ "Unknown event": msg });
        break;
    }
  };

  /**
   * @param {CloseEvent} e
   */
  wss.onclose = (e) => {
    if (e.wasClean) {
      return;
    }

    if (abortController?.signal.aborted) {
      return;
    }

    setTimeout(() => {
      if (abortController?.signal.aborted) {
        return;
      }

      console.error("Websocket closed unexpectedly, reconnecting...");
      initWebSocket(url);
    }, 1000);
  };

  /**
   * @param {ErrorEvent} e
   */
  wss.onerror = (e) => {
    console.error("Websocket error", e);

    if (abortController?.signal.aborted) {
      return;
    }

    setTimeout(() => {
      if (abortController?.signal.aborted) {
        return;
      }

      console.error("Websocket closed unexpectedly, reconnecting...");
      initWebSocket(url);
    }, 1000);
  };
}

function closeWebSocket() {
  if (!wss || wss.readyState !== WebSocket.OPEN) {
    return;
  }

  if (abortController) {
    abortController.abort();
    abortController = null;
  }

  wss.close();
  wss = null;
}

(() => {
  window.addEventListener("beforeunload", closeWebSocket);
  initWebSocket(WS_URL);
})();
