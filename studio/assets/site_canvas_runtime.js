(function () {
  "use strict";

  function ready(fn) {
    if (document.readyState === "loading") {
      document.addEventListener("DOMContentLoaded", fn, { once: true });
      return;
    }
    fn();
  }

  function number(value, fallback) {
    var parsed = Number(value);
    return Number.isFinite(parsed) ? parsed : fallback;
  }

  function clamp(value, min, max) {
    return Math.min(max, Math.max(min, value));
  }

  function emit(canvas, name, detail) {
    canvas.dispatchEvent(new CustomEvent(name, { bubbles: true, detail: detail || {} }));
  }

  function initCanvas(canvas) {
    if (!canvas || canvas.dataset.gosxStudioSiteCanvasBound === "true") return;
    canvas.dataset.gosxStudioSiteCanvasBound = "true";
    var viewport = canvas.querySelector("[data-gosx-studio-canvas-viewport]");
    var surface = canvas.querySelector("[data-gosx-studio-canvas-surface]");
    if (!viewport || !surface) return;

    var state = {
      zoom: number(canvas.getAttribute("data-gosx-studio-canvas-zoom"), 1),
      panX: number(canvas.getAttribute("data-gosx-studio-canvas-pan-x"), 0),
      panY: number(canvas.getAttribute("data-gosx-studio-canvas-pan-y"), 0),
      selected: canvas.getAttribute("data-gosx-studio-canvas-selected") || ""
    };

    function apply(reason) {
      state.zoom = clamp(state.zoom, 0.25, 2.8);
      surface.style.setProperty("--gosx-studio-canvas-pan-x", state.panX + "px");
      surface.style.setProperty("--gosx-studio-canvas-pan-y", state.panY + "px");
      surface.style.setProperty("--gosx-studio-canvas-zoom", state.zoom);
      canvas.setAttribute("data-gosx-studio-canvas-zoom", String(state.zoom));
      canvas.setAttribute("data-gosx-studio-canvas-pan-x", String(Math.round(state.panX)));
      canvas.setAttribute("data-gosx-studio-canvas-pan-y", String(Math.round(state.panY)));
      emit(canvas, "gosxstudio:canvas-viewport", {
        zoom: state.zoom,
        panX: state.panX,
        panY: state.panY,
        reason: reason || ""
      });
    }

    function nodeDetail(node) {
      return {
        key: node.getAttribute("data-gosx-studio-canvas-node") || "",
        kind: node.getAttribute("data-gosx-studio-canvas-node-kind") || "",
        label: node.getAttribute("data-gosx-studio-canvas-node-label") || node.textContent.trim(),
        href: node.getAttribute("data-gosx-studio-canvas-node-href") || ""
      };
    }

    function select(node, reason) {
      canvas.querySelectorAll("[data-gosx-studio-canvas-node]").forEach(function (candidate) {
        var active = candidate === node;
        candidate.classList.toggle("is-selected", active);
        candidate.setAttribute("aria-pressed", active ? "true" : "false");
      });
      var detail = nodeDetail(node);
      state.selected = detail.key;
      canvas.setAttribute("data-gosx-studio-canvas-selected", state.selected);
      detail.reason = reason || "";
      emit(canvas, "gosxstudio:canvas-select", detail);
    }

    function reset() {
      state.zoom = 1;
      state.panX = 0;
      state.panY = 0;
      apply("reset");
    }

    function zoomBy(delta, reason) {
      state.zoom = clamp(state.zoom + delta, 0.25, 2.8);
      apply(reason || "zoom");
    }

    var drag = null;
    viewport.addEventListener("pointerdown", function (event) {
      if (event.target && event.target.closest && event.target.closest("[data-gosx-studio-canvas-node], button, a, input, textarea, select")) return;
      drag = {
        x: event.clientX,
        y: event.clientY,
        panX: state.panX,
        panY: state.panY
      };
      viewport.setPointerCapture(event.pointerId);
      viewport.setAttribute("data-gosx-studio-canvas-dragging", "true");
    });
    viewport.addEventListener("pointermove", function (event) {
      var rect = viewport.getBoundingClientRect();
      emit(canvas, "gosxstudio:canvas-cursor", {
        clientX: event.clientX,
        clientY: event.clientY,
        x: (event.clientX - rect.left - state.panX) / state.zoom,
        y: (event.clientY - rect.top - state.panY) / state.zoom,
        selected: state.selected
      });
      if (!drag) return;
      state.panX = drag.panX + event.clientX - drag.x;
      state.panY = drag.panY + event.clientY - drag.y;
      apply("pan");
    });
    viewport.addEventListener("pointerup", function (event) {
      if (!drag) return;
      drag = null;
      viewport.removeAttribute("data-gosx-studio-canvas-dragging");
      try {
        viewport.releasePointerCapture(event.pointerId);
      } catch (error) {
        return;
      }
    });
    viewport.addEventListener("wheel", function (event) {
      event.preventDefault();
      if (event.ctrlKey || event.metaKey) {
        zoomBy(event.deltaY > 0 ? -0.08 : 0.08, "wheel");
        return;
      }
      state.panX -= event.deltaX;
      state.panY -= event.deltaY;
      apply("wheel-pan");
    }, { passive: false });

    canvas.addEventListener("click", function (event) {
      var node = event.target && event.target.closest && event.target.closest("[data-gosx-studio-canvas-node]");
      if (node && canvas.contains(node)) select(node, "click");
      if (event.target && event.target.closest && event.target.closest("[data-gosx-studio-canvas-reset]")) reset();
      if (event.target && event.target.closest && event.target.closest("[data-gosx-studio-canvas-zoom-in]")) zoomBy(0.12, "button");
      if (event.target && event.target.closest && event.target.closest("[data-gosx-studio-canvas-zoom-out]")) zoomBy(-0.12, "button");
    });
    canvas.addEventListener("dblclick", function (event) {
      var node = event.target && event.target.closest && event.target.closest("[data-gosx-studio-canvas-node]");
      if (!node || !canvas.contains(node)) return;
      var detail = nodeDetail(node);
      emit(canvas, "gosxstudio:canvas-open", detail);
      if (canvas.getAttribute("data-gosx-studio-canvas-open-links") === "true" && detail.href) {
        window.location.href = detail.href;
      }
    });
    canvas.addEventListener("keydown", function (event) {
      if (event.key === "+" || event.key === "=") {
        event.preventDefault();
        zoomBy(0.12, "keyboard");
      } else if (event.key === "-") {
        event.preventDefault();
        zoomBy(-0.12, "keyboard");
      } else if (event.key === "0") {
        event.preventDefault();
        reset();
      } else if (event.key === "ArrowLeft" || event.key === "ArrowRight" || event.key === "ArrowUp" || event.key === "ArrowDown") {
        event.preventDefault();
        state.panX += event.key === "ArrowLeft" ? 42 : event.key === "ArrowRight" ? -42 : 0;
        state.panY += event.key === "ArrowUp" ? 42 : event.key === "ArrowDown" ? -42 : 0;
        apply("keyboard-pan");
      } else if (event.key === "Escape") {
        canvas.querySelectorAll("[data-gosx-studio-canvas-node]").forEach(function (node) {
          node.classList.remove("is-selected");
          node.setAttribute("aria-pressed", "false");
        });
        state.selected = "";
        canvas.removeAttribute("data-gosx-studio-canvas-selected");
        emit(canvas, "gosxstudio:canvas-select", { key: "", reason: "clear" });
      }
    });
    apply("init");
  }

  function initAll(root) {
    var scope = root && root.querySelectorAll ? root : document;
    scope.querySelectorAll("[data-gosx-studio-site-canvas]").forEach(initCanvas);
  }

  ready(function () { initAll(document); });
  document.addEventListener("gosx:navigate", function () { initAll(document); });
  document.addEventListener("gosx:render", function () { initAll(document); });
})();
