"use strict";

(() => {
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
      if (!wss || !abortController) {
        return;
      }

      wss.send(JSON.stringify({ method: "gofutz:init" }));

      window.addEventListener(
        "gofutz:run-all-tests",
        () => {
          sendMessage("gofutz:run-all-tests");
        },
        { signal: abortController.signal },
      );

      window.addEventListener(
        "gofutz:stop-tests",
        () => {
          sendMessage("gofutz:stop-tests");
        },
        { signal: abortController.signal },
      );
    };

    /**
     * @param {MessageEvent} e
     */
    wss.onmessage = (e) => {
      /** @type {InitMessage | UpdateMessage | OutputMessage} */
      const msg = JSON.parse(e.data);

      switch (msg.method) {
        case "gofutz:init":
          if (!("files" in msg.params) || !("output" in msg.params)) {
            return;
          }

          globalThis.testData.files = msg.params.files;
          globalThis.testData.coverage = msg.params.coverage;
          globalThis.testData.isRunning = msg.params.isRunning;
          globalThis.testData.output = msg.params.output;

          window.dispatchEvent(new CustomEvent("gofutz:init"));

          break;

        case "gofutz:update":
          if ("output" in msg.params) {
            return;
          }

          globalThis.testData.coverage = msg.params.coverage;
          globalThis.testData.isRunning = msg.params.isRunning;

          if (msg.params.files) {
            for (const file of Object.values(msg.params.files)) {
              globalThis.testData.files[file.name] = file;
            }
          }

          window.dispatchEvent(new CustomEvent("gofutz:update"));

          break;

        case "gofutz:output":
          if (!("output" in msg.params)) {
            return;
          }

          globalThis.testData.output = msg.params.output;

          window.dispatchEvent(new CustomEvent("gofutz:output"));

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

  /**
   * @param {string} method
   */
  function sendMessage(method) {
    if (!wss) {
      console.error("Websocket not found, could not send message");

      return;
    }

    wss.send(JSON.stringify({ method }));
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

  window.addEventListener("beforeunload", closeWebSocket);
  initWebSocket(WS_URL);
})();
