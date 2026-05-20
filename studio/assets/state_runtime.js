(function () {
  "use strict";

  var states = {
    saved: "Saved",
    dirty: "Unsaved",
    saving: "Saving",
    autosaving: "Autosaving",
    error: "Save failed"
  };

  function ready(fn) {
    if (document.readyState === "loading") {
      document.addEventListener("DOMContentLoaded", fn, { once: true });
      return;
    }
    fn();
  }

  function frameTask(fn) {
    var pending = false;
    return function () {
      if (pending) return;
      pending = true;
      window.requestAnimationFrame(function () {
        pending = false;
        fn();
      });
    };
  }

  function skipField(field) {
    if (!field || !field.name || field.disabled) return true;
    if (field.name === "csrf_token") return true;
    if (field.matches && field.matches("[data-gosx-studio-state-ignore]")) return true;
    var type = String(field.type || "").toLowerCase();
    return type === "button" || type === "submit" || type === "reset" || type === "file";
  }

  function fieldValues(field) {
    var type = String(field.type || "").toLowerCase();
    if (type === "checkbox" || type === "radio") {
      return [field.checked ? (field.value || "on") : ""];
    }
    if (field.tagName === "SELECT" && field.multiple) {
      return Array.prototype.slice.call(field.options)
        .filter(function (option) { return option.selected; })
        .map(function (option) { return option.value; });
    }
    return [field.value || ""];
  }

  function formSignature(form) {
    var parts = [];
    Array.prototype.forEach.call(form.elements, function (field) {
      if (skipField(field)) return;
      fieldValues(field).forEach(function (value) {
        parts.push(encodeURIComponent(field.name) + "=" + encodeURIComponent(value));
      });
    });
    return parts.join("&");
  }

  function dirtyCount(savedSignature, currentSignature) {
    if (savedSignature === currentSignature) return 0;
    var saved = savedSignature ? savedSignature.split("&") : [];
    var current = currentSignature ? currentSignature.split("&") : [];
    var total = Math.max(saved.length, current.length);
    var count = 0;
    for (var index = 0; index < total; index += 1) {
      if (saved[index] !== current[index]) count += 1;
    }
    return Math.max(count, 1);
  }

  function statusNodes(form) {
    return Array.prototype.slice.call(form.querySelectorAll("[data-gosx-studio-save-state], [data-editor-save-state]"));
  }

  function detailNodes(form) {
    return Array.prototype.slice.call(form.querySelectorAll("[data-gosx-studio-save-detail]"));
  }

  function lastSavedNodes(form) {
    return Array.prototype.slice.call(form.querySelectorAll("[data-gosx-studio-last-saved]"));
  }

  function dirtyCountNodes(form) {
    return Array.prototype.slice.call(form.querySelectorAll("[data-gosx-studio-dirty-count]"));
  }

  function saveButtons(form) {
    return Array.prototype.slice.call(form.querySelectorAll("[data-gosx-studio-save-button], [data-editor-save-button]"));
  }

  function baseClass(node) {
    if (!node.dataset.gosxStudioSaveStateBaseClass) {
      var keep = Array.prototype.slice.call(node.classList).filter(function (className) {
        return className.indexOf("gosx-studio-save-state--") !== 0 && className.indexOf("editor-save-state--") !== 0;
      });
      node.dataset.gosxStudioSaveStateBaseClass = keep.join(" ");
    }
    return node.dataset.gosxStudioSaveStateBaseClass;
  }

  function setStatusClass(node, state) {
    var classes = baseClass(node);
    if (classes) classes += " ";
    classes += "gosx-studio-save-state--" + state;
    if (node.classList.contains("editor-save-state") || classes.indexOf("editor-save-state") >= 0) {
      classes += " editor-save-state--" + state;
    }
    node.className = classes;
  }

  function dispatchState(form, state, dirty, reason, options) {
    options = options || {};
    form.dispatchEvent(new CustomEvent("gosxstudio:save-state", {
      bubbles: true,
      detail: {
        state: state,
        dirty: dirty,
        reason: reason || "",
        dirtyCount: options.dirtyCount || 0,
        savedAt: options.savedAt || ""
      }
    }));
  }

  function timeLabel(date) {
    try {
      return date.toLocaleTimeString([], { hour: "numeric", minute: "2-digit" });
    } catch (error) {
      return "";
    }
  }

  function saveDetail(state, options) {
    options = options || {};
    if (state === "dirty") {
      return options.dirtyCount === 1 ? "1 change waiting" : String(options.dirtyCount || 0) + " changes waiting";
    }
    if (state === "autosaving") return "Syncing changes";
    if (state === "saving") return "Sending checkpoint";
    if (state === "error") return "Autosave could not reach the server";
    if (options.savedAtLabel) return "Saved at " + options.savedAtLabel;
    return "Ready";
  }

  function updateLastSaved(form, savedAt) {
    lastSavedNodes(form).forEach(function (node) {
      if (!savedAt) {
        node.textContent = node.getAttribute("data-gosx-studio-last-saved-empty") || "Not saved this session";
        node.removeAttribute("datetime");
        return;
      }
      var date = new Date(savedAt);
      var label = timeLabel(date);
      node.textContent = label ? "Last saved " + label : "Last saved";
      if (node.tagName === "TIME") node.setAttribute("datetime", savedAt);
    });
  }

  function updateDirtyCount(form, count) {
    form.setAttribute("data-gosx-studio-dirty-count", String(count || 0));
    dirtyCountNodes(form).forEach(function (node) {
      node.hidden = !count;
      node.textContent = count === 1 ? "1 change" : String(count || 0) + " changes";
    });
  }

  function setState(form, state, reason, options) {
    options = options || {};
    var dirty = state === "dirty" || state === "autosaving" || state === "error";
    var count = options.dirtyCount || 0;
    form.setAttribute("data-gosx-studio-save-state", state);
    form.setAttribute("data-studio-dirty-state", dirty ? "dirty" : "clean");
    statusNodes(form).forEach(function (node) {
      node.textContent = states[state] || states.saved;
      node.setAttribute("data-gosx-studio-save-state-value", state);
      setStatusClass(node, state);
    });
    detailNodes(form).forEach(function (node) {
      node.textContent = saveDetail(state, options);
    });
    updateDirtyCount(form, count);
    updateLastSaved(form, options.savedAt || form.getAttribute("data-gosx-studio-last-saved-at") || "");
    saveButtons(form).forEach(function (button) {
      button.classList.toggle("button--attention", state === "dirty" || state === "error");
      button.setAttribute("data-gosx-studio-save-button-state", state);
      button.setAttribute("aria-busy", state === "saving" || state === "autosaving" ? "true" : "false");
    });
    dispatchState(form, state, dirty, reason, options);
  }

  function autosaveEnabled(form) {
    return form.getAttribute("data-gosx-studio-autosave") === "true";
  }

  function autosaveDelay(form) {
    var delay = Number(form.getAttribute("data-gosx-studio-autosave-delay") || 1400);
    if (!Number.isFinite(delay) || delay < 250) return 1400;
    return delay;
  }

  function autosaveURL(form) {
    return form.getAttribute("data-gosx-studio-autosave-url") || form.getAttribute("action") || window.location.href;
  }

  function autosaveMethod(form) {
    return String(form.getAttribute("method") || "post").toUpperCase();
  }

  function initForm(form) {
    if (!form || form.dataset.gosxStudioStateBound === "true") return;
    form.dataset.gosxStudioStateBound = "true";
    var saved = formSignature(form);
    var submitting = false;
    var autosaving = false;
    var autosaveTimer = 0;
    var lastSavedAt = "";

    function isDirty() {
      return formSignature(form) !== saved;
    }

    function stateOptions(extra) {
      var current = formSignature(form);
      var options = {
        dirtyCount: dirtyCount(saved, current),
        savedAt: lastSavedAt,
        savedAtLabel: lastSavedAt ? timeLabel(new Date(lastSavedAt)) : ""
      };
      if (extra) {
        Object.keys(extra).forEach(function (key) {
          options[key] = extra[key];
        });
      }
      return options;
    }

    function update(reason) {
      if (submitting || autosaving) return;
      setState(form, isDirty() ? "dirty" : "saved", reason, stateOptions());
    }

    var updateFrame = frameTask(function () {
      update("field");
    });

    function scheduleAutosave() {
      if (!autosaveEnabled(form)) return;
      window.clearTimeout(autosaveTimer);
      autosaveTimer = window.setTimeout(runAutosave, autosaveDelay(form));
    }

    function runAutosave() {
      if (!autosaveEnabled(form) || submitting || autosaving || !isDirty()) return;
      var signature = formSignature(form);
      var target = autosaveURL(form);
      if (!target) return;
      autosaving = true;
      setState(form, "autosaving", "autosave", stateOptions());
      window.fetch(target, {
        method: autosaveMethod(form),
        body: new FormData(form),
        credentials: "same-origin",
        headers: {
          "X-GoSX-Studio-Autosave": "true"
        }
      }).then(function (response) {
        if (!response.ok) throw new Error("Autosave failed");
        saved = signature;
        lastSavedAt = new Date().toISOString();
        form.setAttribute("data-gosx-studio-last-saved-at", lastSavedAt);
        autosaving = false;
        if (isDirty()) {
          setState(form, "dirty", "autosave-stale", stateOptions());
          scheduleAutosave();
          return;
        }
        setState(form, "saved", "autosave", stateOptions({ dirtyCount: 0 }));
      }, function () {
        autosaving = false;
        setState(form, "error", "autosave", stateOptions());
      });
    }

    form.addEventListener("input", function () {
      updateFrame();
      scheduleAutosave();
    });
    form.addEventListener("change", function () {
      updateFrame();
      scheduleAutosave();
    });
    form.addEventListener("click", function (event) {
      if (event.target && event.target.closest && event.target.closest("[data-block-studio-move], [data-gosx-studio-state-dirty-action]")) {
        updateFrame();
        scheduleAutosave();
      }
    });
    form.addEventListener("drop", function () {
      updateFrame();
      scheduleAutosave();
    });
    form.addEventListener("dragend", function () {
      updateFrame();
      scheduleAutosave();
    });
    form.addEventListener("submit", function () {
      submitting = true;
      window.clearTimeout(autosaveTimer);
      setState(form, "saving", "submit", stateOptions());
    });

    document.addEventListener("keydown", function (event) {
      if (event.defaultPrevented || !document.contains(form)) return;
      if (!(event.metaKey || event.ctrlKey) || event.key.toLowerCase() !== "s") return;
      event.preventDefault();
      if (form.requestSubmit) form.requestSubmit();
      else form.submit();
    });

    document.addEventListener("click", function (event) {
      var link = event.target && event.target.closest && event.target.closest("a[href]");
      if (!link || !document.contains(form) || !isDirty() || submitting) return;
      if (link.target && link.target !== "_self") return;
      var next = new URL(link.href, window.location.href);
      if (next.origin === window.location.origin && next.pathname === window.location.pathname && next.search === window.location.search) return;
      if (!window.confirm("Discard unsaved editor changes?")) event.preventDefault();
    }, true);

    window.addEventListener("beforeunload", function (event) {
      if (!document.contains(form) || !isDirty() || submitting) return;
      event.preventDefault();
      event.returnValue = "";
    });

    update("init");
  }

  function initAll(root) {
    var scope = root && root.querySelectorAll ? root : document;
    Array.prototype.forEach.call(scope.querySelectorAll("[data-gosx-studio-state]"), initForm);
  }

  ready(function () { initAll(document); });
  document.addEventListener("gosx:navigate", function () { initAll(document); });
  document.addEventListener("gosx:render", function () { initAll(document); });
})();
