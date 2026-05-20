(function () {
  "use strict";

  function ready(fn) {
    if (document.readyState === "loading") {
      document.addEventListener("DOMContentLoaded", fn, { once: true });
      return;
    }
    fn();
  }

  function text(node, value) {
    if (node) node.textContent = value;
  }

  function compact(value) {
    return String(value || "").replace(/\s+/g, " ").trim();
  }

  function storageKey(scope) {
    return "gosx-studio-flow:" + compact(scope.getAttribute("data-studio-flow-storage-key") || scope.getAttribute("data-studio-shell") || window.location.pathname || "studio");
  }

  function storedFlowKey(scope) {
    try {
      return window.sessionStorage ? window.sessionStorage.getItem(storageKey(scope)) || "" : "";
    } catch (error) {
      return "";
    }
  }

  function storeFlowKey(scope, key) {
    try {
      if (window.sessionStorage) window.sessionStorage.setItem(storageKey(scope), key);
    } catch (error) {
      return;
    }
  }

  function hashFlowKey(flowByKey) {
    if (!window.location.hash) return "";
    var value = "";
    try {
      value = decodeURIComponent(window.location.hash.slice(1));
    } catch (error) {
      value = window.location.hash.slice(1);
    }
    if (value.indexOf("flow=") === 0) value = value.slice("flow=".length);
    else if (value.indexOf("flow-") === 0) value = value.slice("flow-".length);
    return flowByKey.has(value) ? value : "";
  }

  function queryAll(root, selector) {
    return Array.prototype.slice.call(root.querySelectorAll(selector));
  }

  function emit(scope, name, detail) {
    scope.dispatchEvent(new CustomEvent(name, { bubbles: true, detail: detail || {} }));
  }

  function initFlowEditor(scope) {
    if (!scope || scope.dataset.gosxStudioFlowRuntimeBound === "true") return;
    var cards = queryAll(scope, "[data-studio-flow-card]");
    var editors = queryAll(scope, "[data-studio-flow-editor]");
    if (!cards.length || !editors.length) return;
    scope.dataset.gosxStudioFlowRuntimeBound = "true";

    var preview = scope.querySelector("[data-studio-preview-frame]");
    var selectedLabel = scope.querySelector("[data-studio-selected-flow-label]");
    var selectedRoute = scope.querySelector("[data-studio-selected-flow-route]");
    var previewStatus = scope.querySelector("[data-studio-preview-status]");
    var openPreview = scope.querySelector("[data-studio-open-preview]");
    var editorByKey = new Map();
    editors.forEach(function (editor) {
      editorByKey.set(editor.getAttribute("data-studio-flow-editor") || "", editor);
    });

    var flows = cards.map(function (card, index) {
      var key = card.getAttribute("data-studio-flow-card") || "";
      return {
        key: key,
        index: index,
        card: card,
        editor: editorByKey.get(key),
        route: card.getAttribute("data-studio-flow-route") || "",
        label: card.getAttribute("data-studio-flow-label") || compact(card.textContent),
        status: card.getAttribute("data-studio-flow-status") || ""
      };
    }).filter(function (flow) {
      return !!flow.key;
    });
    if (!flows.length) return;

    var flowByKey = new Map();
    flows.forEach(function (flow) { flowByKey.set(flow.key, flow); });

    function setPreviewState(state) {
      var label = "Ready";
      if (state === "loading") label = "Loading";
      else if (state === "error") label = "Preview failed";
      if (preview) preview.setAttribute("data-studio-preview-state", state);
      text(previewStatus, label);
    }

    function setPreviewMeta(flow) {
      text(selectedLabel, flow.label);
      text(selectedRoute, flow.route || "No public route");
      if (!openPreview) return;
      if (flow.route) {
        openPreview.setAttribute("href", flow.route);
        openPreview.removeAttribute("aria-disabled");
        openPreview.classList.remove("is-disabled");
      } else {
        openPreview.setAttribute("href", "#");
        openPreview.setAttribute("aria-disabled", "true");
        openPreview.classList.add("is-disabled");
      }
    }

    function setPreview(route) {
      if (!preview || !route) return;
      if (preview.getAttribute("src") === route) {
        setPreviewState("ready");
        return;
      }
      setPreviewState("loading");
      preview.setAttribute("src", route);
    }

    function selectFlow(key, route, options) {
      if (!key || !flowByKey.has(key)) return;
      options = options || {};
      var flow = flowByKey.get(key);
      flows.forEach(function (item) {
        var active = item.key === key;
        item.card.classList.toggle("is-selected", active);
        item.card.setAttribute("aria-selected", active ? "true" : "false");
        item.card.tabIndex = active ? 0 : -1;
        if (item.editor) {
          item.editor.classList.toggle("is-selected", active);
          item.editor.hidden = !active;
          item.editor.setAttribute("aria-hidden", active ? "false" : "true");
        }
      });
      scope.setAttribute("data-studio-selected-flow", key);
      setPreviewMeta(flow);
      if (options.persist !== false) storeFlowKey(scope, key);
      if (options.focus && flow.card.focus) flow.card.focus();
      if (options.preview !== false) setPreview(route || flow.route);
      emit(scope, "gosxstudio:flow-select", {
        key: key,
        label: flow.label,
        route: route || flow.route,
        status: flow.status,
        reason: options.reason || ""
      });
    }

    function selectByIndex(index) {
      var flow = flows[index];
      if (flow) selectFlow(flow.key, flow.route, { focus: true, reason: "keyboard" });
    }

    function findFlowKeyByRoute(route) {
      for (var index = 0; index < flows.length; index += 1) {
        if (flows[index].route === route) return flows[index].key;
      }
      return "";
    }

    function updateDirtyState(editor) {
      var flow = editor && flowByKey.get(editor.getAttribute("data-studio-flow-editor") || "");
      if (!flow) return;
      var dirtyCount = 0;
      queryAll(editor, "[data-studio-initial-value]").forEach(function (input) {
        if (input.value !== input.getAttribute("data-studio-initial-value")) dirtyCount += 1;
      });
      var dirtyLabel = dirtyCount === 1 ? "1 unsaved" : String(dirtyCount) + " unsaved";
      editor.classList.toggle("is-dirty", dirtyCount > 0);
      flow.card.classList.toggle("is-dirty", dirtyCount > 0);
      flow.card.setAttribute("data-studio-dirty-count", String(dirtyCount));
      queryAll(flow.card, "[data-studio-flow-dirty-badge]").forEach(function (badge) {
        badge.hidden = dirtyCount === 0;
        text(badge, dirtyLabel);
      });
      queryAll(editor, "[data-studio-flow-editor-dirty]").forEach(function (badge) {
        badge.hidden = dirtyCount === 0;
        text(badge, dirtyLabel);
      });
      emit(scope, "gosxstudio:flow-dirty", {
        key: flow.key,
        label: flow.label,
        dirtyCount: dirtyCount
      });
    }

    function syncSavedFlowState() {
      editors.forEach(function (editor) {
        queryAll(editor, "input, textarea, select").forEach(function (input) {
          input.setAttribute("data-studio-initial-value", input.value || "");
        });
        updateDirtyState(editor);
      });
    }

    flows.forEach(function (flow) {
      flow.card.addEventListener("click", function () {
        selectFlow(flow.key, flow.route, { reason: "click" });
      });
      flow.card.addEventListener("keydown", function (event) {
        if (event.key === "ArrowRight" || event.key === "ArrowDown") {
          event.preventDefault();
          selectByIndex((flow.index + 1) % flows.length);
        } else if (event.key === "ArrowLeft" || event.key === "ArrowUp") {
          event.preventDefault();
          selectByIndex((flow.index - 1 + flows.length) % flows.length);
        } else if (event.key === "Home") {
          event.preventDefault();
          selectByIndex(0);
        } else if (event.key === "End") {
          event.preventDefault();
          selectByIndex(flows.length - 1);
        }
      });
    });

    scope.addEventListener("click", function (event) {
      if (openPreview && event.target && event.target.closest && event.target.closest("[data-studio-open-preview]") === openPreview && openPreview.getAttribute("aria-disabled") === "true") {
        event.preventDefault();
        return;
      }
      var button = event.target && event.target.closest && event.target.closest("[data-studio-preview-flow]");
      if (!button || !scope.contains(button)) return;
      var route = button.getAttribute("data-studio-preview-flow") || "";
      var owningEditor = button.closest("[data-studio-flow-editor]");
      var flowKey = owningEditor ? owningEditor.getAttribute("data-studio-flow-editor") : findFlowKeyByRoute(route);
      if (flowKey) selectFlow(flowKey, route, { preview: false, reason: "preview-button" });
      setPreview(route);
      emit(scope, "gosxstudio:flow-preview", { key: flowKey || "", route: route });
    });

    scope.addEventListener("input", function (event) {
      var editor = event.target && event.target.closest && event.target.closest("[data-studio-flow-editor]");
      if (editor && scope.contains(editor)) updateDirtyState(editor);
    });

    scope.addEventListener("gosxstudio:command", function (event) {
      var detail = event.detail || {};
      if (detail.href && detail.href.indexOf("#flow=") === 0) {
        var key = detail.href.slice("#flow=".length);
        if (flowByKey.has(key)) selectFlow(key, flowByKey.get(key).route, { preview: true, persist: true, reason: "command" });
      }
    });
    scope.addEventListener("gosxstudio:save-state", function (event) {
      if (event.detail && event.detail.state === "saved") syncSavedFlowState();
    });

    if (preview) {
      preview.addEventListener("load", function () { setPreviewState("ready"); });
      preview.addEventListener("error", function () { setPreviewState("error"); });
    }

    editors.forEach(function (editor) {
      queryAll(editor, "input, textarea, select").forEach(function (input) {
        if (!input.hasAttribute("data-studio-initial-value")) {
          input.setAttribute("data-studio-initial-value", input.value || "");
        }
      });
      updateDirtyState(editor);
    });

    var initialKey = hashFlowKey(flowByKey) || storedFlowKey(scope);
    if (!flowByKey.has(initialKey)) initialKey = flows[0].key;
    selectFlow(initialKey, flowByKey.get(initialKey).route, { preview: true, persist: true, reason: "init" });
    window.addEventListener("hashchange", function () {
      var key = hashFlowKey(flowByKey);
      if (key) selectFlow(key, flowByKey.get(key).route, { preview: true, persist: true, reason: "hash" });
    });
  }

  function initAll(root) {
    var scope = root && root.querySelectorAll ? root : document;
    queryAll(scope, "form[data-studio-workbench], [data-studio-flow-runtime]").forEach(initFlowEditor);
  }

  ready(function () { initAll(document); });
  document.addEventListener("gosx:navigate", function () { initAll(document); });
  document.addEventListener("gosx:render", function () { initAll(document); });
})();
