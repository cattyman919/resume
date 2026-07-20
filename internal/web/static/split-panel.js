(function () {
  const PANEL_KEY = 'autocv-panel-sizes';
  const MIN_EDITOR = 320;
  const MIN_PREVIEW = 400;

  function getStoredSizes() {
    try {
      const raw = localStorage.getItem(PANEL_KEY);
      if (raw) return JSON.parse(raw);
    } catch (e) {}
    return null;
  }

  function storeSizes(editorW, previewW) {
    try {
      localStorage.setItem(PANEL_KEY, JSON.stringify({ editorW, previewW }));
    } catch (e) {}
  }

  function createDivider() {
    const main = document.getElementById('main-layout');
    if (!main || main.dataset.splitInit === 'true') return;
    main.dataset.splitInit = 'true';

    const editor = document.getElementById('editor-panel');
    const preview = document.getElementById('pdf-panel');
    if (!editor || !preview) return;

    const divider = document.createElement('div');
    divider.id = 'panel-divider';
    divider.className = 'panel-divider';
    divider.setAttribute('role', 'separator');
    divider.setAttribute('aria-orientation', 'vertical');
    divider.setAttribute('tabindex', '0');

    editor.parentNode.insertBefore(divider, preview);

    const stored = getStoredSizes();
    if (stored && stored.editorW && stored.previewW) {
      editor.style.width = stored.editorW + 'px';
      preview.style.width = stored.previewW + 'px';
    } else {
      editor.style.width = '50%';
      preview.style.width = '50%';
    }

    let isDragging = false;
    let startX = 0;
    let startEditorW = 0;

    function onPointerDown(e) {
      isDragging = true;
      startX = e.clientX || (e.touches && e.touches[0].clientX);
      startEditorW = editor.offsetWidth;
      document.body.style.cursor = 'col-resize';
      document.body.style.userSelect = 'none';
      divider.classList.add('active');
      e.preventDefault();
    }

    function onPointerMove(e) {
      if (!isDragging) return;
      const clientX = e.clientX || (e.touches && e.touches[0].clientX);
      const dx = clientX - startX;
      const newEditorW = Math.max(MIN_EDITOR, startEditorW + dx);
      const containerW = main.offsetWidth - divider.offsetWidth;
      const maxEditorW = containerW - MIN_PREVIEW;
      const clampedEditorW = Math.min(newEditorW, maxEditorW);
      const clampedPreviewW = containerW - clampedEditorW;

      if (clampedPreviewW >= MIN_PREVIEW) {
        editor.style.width = clampedEditorW + 'px';
        preview.style.width = clampedPreviewW + 'px';
      }
    }

    function onPointerUp() {
      if (!isDragging) return;
      isDragging = false;
      document.body.style.cursor = '';
      document.body.style.userSelect = '';
      divider.classList.remove('active');
      storeSizes(editor.offsetWidth, preview.offsetWidth);
      window.dispatchEvent(new Event('resize'));
    }

    divider.addEventListener('mousedown', onPointerDown);
    divider.addEventListener('touchstart', onPointerDown, { passive: false });
    document.addEventListener('mousemove', onPointerMove);
    document.addEventListener('touchmove', onPointerMove, { passive: false });
    document.addEventListener('mouseup', onPointerUp);
    document.addEventListener('touchend', onPointerUp);

    divider.addEventListener('keydown', function (e) {
      const step = 20;
      const containerW = main.offsetWidth - divider.offsetWidth;
      if (e.key === 'ArrowLeft' || e.key === 'ArrowRight') {
        e.preventDefault();
        const currentW = editor.offsetWidth;
        const delta = e.key === 'ArrowLeft' ? -step : step;
        const newW = Math.max(MIN_EDITOR, Math.min(currentW + delta, containerW - MIN_PREVIEW));
        editor.style.width = newW + 'px';
        preview.style.width = (containerW - newW) + 'px';
        storeSizes(newW, containerW - newW);
        window.dispatchEvent(new Event('resize'));
      }
    });
  }

  function initMobileTabs() {
    const tabs = document.querySelectorAll('.mobile-tab');
    tabs.forEach(function (tab) {
      if (tab.dataset.bound) return;
      tab.dataset.bound = 'true';
      tab.addEventListener('click', function () {
        const target = tab.dataset.tab;
        document.querySelectorAll('.mobile-tab').forEach(function (t) {
          t.classList.remove('active');
        });
        tab.classList.add('active');

        const editor = document.getElementById('editor-panel');
        const preview = document.getElementById('pdf-panel');
        const divider = document.getElementById('panel-divider');

        if (target === 'editor') {
          editor.style.display = '';
          preview.style.display = 'none';
          if (divider) divider.style.display = 'none';
        } else {
          editor.style.display = 'none';
          preview.style.display = '';
          if (divider) divider.style.display = 'none';
          window.dispatchEvent(new Event('resize'));
        }
      });
    });
  }

  function handleResize() {
    const main = document.getElementById('main-layout');
    const divider = document.getElementById('panel-divider');
    const editor = document.getElementById('editor-panel');
    const preview = document.getElementById('pdf-panel');
    const mobileTabs = document.getElementById('mobile-tabs');

    if (!main || !editor || !preview) return;

    const isMobile = window.innerWidth < 768;

    if (isMobile) {
      editor.style.width = '100%';
      preview.style.width = '100%';
      if (divider) divider.style.display = 'none';
      if (mobileTabs) mobileTabs.style.display = 'flex';

      const activeTab = document.querySelector('.mobile-tab.active');
      if (activeTab && activeTab.dataset.tab === 'editor') {
        editor.style.display = '';
        preview.style.display = 'none';
      } else {
        editor.style.display = 'none';
        preview.style.display = '';
      }
    } else {
      if (mobileTabs) mobileTabs.style.display = 'none';
      editor.style.display = '';
      preview.style.display = '';
      if (divider) divider.style.display = '';
      createDivider();
    }
  }

  document.addEventListener('DOMContentLoaded', function () {
    handleResize();
    initMobileTabs();
    let resizeTimer;
    window.addEventListener('resize', function () {
      clearTimeout(resizeTimer);
      resizeTimer = setTimeout(handleResize, 150);
    });
  });

  window.AutoCV = window.AutoCV || {};
  window.AutoCV.handleResize = handleResize;
})();
