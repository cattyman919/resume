(function () {
  function initKeyboardShortcuts() {
    document.addEventListener('keydown', function (e) {
      const isMac = navigator.platform.toUpperCase().indexOf('MAC') >= 0;
      const mod = isMac ? e.metaKey : e.ctrlKey;

      if (mod && e.key === 's') {
        e.preventDefault();
        return;
      }

      if (mod && e.key === 'p') {
        e.preventDefault();
        const typeName = document.body.dataset.currentType;
        if (typeName && typeof generatePDF === 'function') {
          generatePDF(typeName);
        }
        return;
      }

      if (mod && e.key === 'd') {
        e.preventDefault();
        if (window.AutoCV && window.AutoCV.toggleTheme) {
          window.AutoCV.toggleTheme();
        }
        return;
      }

      if (e.key === '?' && !e.ctrlKey && !e.metaKey && !e.altKey) {
        const active = document.activeElement;
        if (active && (active.tagName === 'INPUT' || active.tagName === 'TEXTAREA' || active.tagName === 'SELECT')) {
          return;
        }
        toggleShortcutHelp();
        return;
      }

      if (e.key === 'Escape') {
        const help = document.getElementById('shortcut-help');
        if (help && !help.classList.contains('hidden')) {
          help.classList.add('hidden');
        }
      }
    });
  }

  function toggleShortcutHelp() {
    let help = document.getElementById('shortcut-help');
    if (!help) {
      help = document.createElement('div');
      help.id = 'shortcut-help';
      help.className = 'shortcut-help';
      help.innerHTML =
        '<div class="modal-backdrop" onclick="document.getElementById(\'shortcut-help\').classList.add(\'hidden\')">' +
        '<div class="modal-content" onclick="event.stopPropagation()">' +
        '<h3 class="text-lg font-bold mb-4" style="color: var(--color-foreground);">Keyboard Shortcuts</h3>' +
        '<div class="space-y-2" style="color: var(--color-foreground);">' +
        '<div class="flex justify-between text-sm"><span style="color: var(--color-muted-foreground);">Generate PDF</span><kbd class="kbd">Ctrl+P</kbd></div>' +
        '<div class="flex justify-between text-sm"><span style="color: var(--color-muted-foreground);">Toggle Dark Mode</span><kbd class="kbd">Ctrl+D</kbd></div>' +
        '<div class="flex justify-between text-sm"><span style="color: var(--color-muted-foreground);">Show Shortcuts</span><kbd class="kbd">?</kbd></div>' +
        '<div class="flex justify-between text-sm"><span style="color: var(--color-muted-foreground);">Close Dialog</span><kbd class="kbd">Esc</kbd></div>' +
        '<div class="flex justify-between text-sm"><span style="color: var(--color-muted-foreground);">Resize Panel</span><kbd class="kbd">\u2190 \u2192</kbd></div>' +
        '</div>' +
        '<div class="flex justify-end mt-4"><button class="btn btn-ghost" onclick="document.getElementById(\'shortcut-help\').classList.add(\'hidden\')">Close</button></div>' +
        '</div></div>';
      document.body.appendChild(help);
    }
    help.classList.toggle('hidden');
  }

  document.addEventListener('DOMContentLoaded', initKeyboardShortcuts);

  window.AutoCV = window.AutoCV || {};
  window.AutoCV.toggleShortcutHelp = toggleShortcutHelp;
})();
