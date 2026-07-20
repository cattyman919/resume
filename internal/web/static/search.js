(function () {
  function initSearch() {
    const input = document.getElementById('section-search');
    if (!input || input.dataset.bound) return;
    input.dataset.bound = 'true';

    input.addEventListener('input', function () {
      const query = input.value.toLowerCase().trim();
      const sections = document.querySelectorAll('.section-card');
      sections.forEach(function (section) {
        if (!query) {
          section.style.display = '';
          return;
        }
        const text = section.textContent.toLowerCase();
        section.style.display = text.includes(query) ? '' : 'none';
      });
    });
  }

  document.addEventListener('DOMContentLoaded', initSearch);
  document.body.addEventListener('htmx:afterSettle', initSearch);

  window.AutoCV = window.AutoCV || {};
  window.AutoCV.initSearch = initSearch;
})();
