let debounceTimers = {};
let defaultDebounceMs = 50;
const DEBOUNCE_KEY = 'autocv-debounce-ms';
let currentZoom = 1.0;
let currentPdf = null;

(function() {
	const stored = localStorage.getItem(DEBOUNCE_KEY);
	if (stored) {
		const val = parseInt(stored, 10);
		if (!isNaN(val) && val >= 0 && val <= 5000) defaultDebounceMs = val;
	}
})();

function getDebounceMs() {
	return defaultDebounceMs;
}

function setDebounceMs(ms) {
	defaultDebounceMs = Math.max(0, Math.min(5000, ms));
	localStorage.setItem(DEBOUNCE_KEY, String(defaultDebounceMs));
}

function debounceGenerate(typeName) {
	if (debounceTimers[typeName]) clearTimeout(debounceTimers[typeName]);
	debounceTimers[typeName] = setTimeout(() => generatePDF(typeName), getDebounceMs());
}

function generatePDF(typeName) {
	setPdfStatus('generating');
	if (window.AutoCV && window.AutoCV.showToast) window.AutoCV.showToast('Generating PDF...', 'info', 2000);
	fetch('/api/generate/' + typeName, { method: 'POST' })
		.then(r => r.json())
		.then(data => {
			if (data.status === 'done') {
				loadPdf(typeName);
				setPdfStatus('success');
				if (window.AutoCV && window.AutoCV.showToast) window.AutoCV.showToast('PDF generated successfully', 'success');
			} else if (data.status === 'error') {
				setPdfStatus('error', data.message || 'Generation failed');
				if (window.AutoCV && window.AutoCV.showToast) window.AutoCV.showToast('PDF generation failed: ' + (data.message || 'Unknown error'), 'error');
			} else if (data.status === 'already_generating') {
				setPdfStatus('generating');
				setTimeout(() => debounceGenerate(typeName), 1000);
			}
		})
		.catch(err => {
			setPdfStatus('error', err.message);
			if (window.AutoCV && window.AutoCV.showToast) window.AutoCV.showToast('Network error: ' + err.message, 'error');
		});
}

function loadPdf(typeName) {
	const url = '/api/pdf/' + typeName + '?t=' + Date.now();
	if (typeof pdfjsLib !== 'undefined') {
		pdfjsLib.getDocument(url).promise.then(pdf => {
			currentPdf = pdf;
			window.pdfDoc = pdf;
			renderAllPages(pdf);
			updatePageInfo(pdf);
		}).catch(err => setPdfStatus('error', err.message));
	}
}

function updatePageInfo(pdf) {
	const el = document.getElementById('page-info');
	if (el && pdf) {
		el.textContent = pdf.numPages + ' page' + (pdf.numPages > 1 ? 's' : '');
	}
}

function renderAllPages(pdf) {
	const container = document.getElementById('pdf-pages-container');
	if (!container) return;
	container.innerHTML = '';
	const viewerWidth = document.getElementById('pdf-viewer').clientWidth - 40;
	const dpr = window.devicePixelRatio || 1;
	let rendered = 0;
	for (let i = 1; i <= pdf.numPages; i++) {
		const wrapper = document.createElement('div');
		wrapper.className = 'pdf-page-wrapper';
		const canvas = document.createElement('canvas');
		wrapper.appendChild(canvas);
		container.appendChild(wrapper);
		pdf.getPage(i).then(page => {
			const baseViewport = page.getViewport({scale: 1});
			const scale = (viewerWidth / baseViewport.width) * currentZoom;
			const viewport = page.getViewport({scale: scale});
			canvas.width = Math.floor(viewport.width * dpr);
			canvas.height = Math.floor(viewport.height * dpr);
			canvas.style.width = Math.floor(viewport.width) + 'px';
			canvas.style.height = Math.floor(viewport.height) + 'px';
			const ctx = canvas.getContext('2d');
			ctx.scale(dpr, dpr);
			page.render({canvasContext: ctx, viewport: viewport}).promise.then(() => {
				rendered++;
				if (rendered === pdf.numPages) {
					document.getElementById('pdf-status-text').textContent = 'Generated (' + pdf.numPages + ' page' + (pdf.numPages > 1 ? 's' : '') + ')';
				}
			});
		});
	}
}

function setZoom(zoom) {
	currentZoom = Math.max(0.5, Math.min(3.0, zoom));
	const display = document.getElementById('zoom-display');
	if (display) display.textContent = Math.round(currentZoom * 100) + '%';
	if (currentPdf) renderAllPages(currentPdf);
}

function zoomIn() { setZoom(currentZoom + 0.1); }
function zoomOut() { setZoom(currentZoom - 0.1); }
function zoomFit() {
	currentZoom = 1.0;
	if (currentPdf) renderAllPages(currentPdf);
	const display = document.getElementById('zoom-display');
	if (display) display.textContent = '100%';
}

function setPdfStatus(status, msg) {
	const icon = document.getElementById('pdf-status-icon');
	const text = document.getElementById('pdf-status-text');
	const errorEl = document.getElementById('pdf-error');
	const colors = { idle: '#94A3B8', generating: '#F59E0B', success: '#22C55E', error: '#EF4444' };
	const labels = { idle: 'Idle', generating: 'Generating...', success: 'Generated', error: 'Failed' };
	const dotEl = icon ? icon.querySelector('.status-dot') : null;
	if (dotEl) dotEl.setAttribute('fill', colors[status] || colors.idle);
	if (text) text.textContent = labels[status] || 'Idle';
	if (status === 'error') {
		if (errorEl) { errorEl.textContent = msg || 'Unknown error'; errorEl.classList.remove('hidden'); }
	} else {
		if (errorEl) errorEl.classList.add('hidden');
	}
}

function makeSvgStr(tag, attrs, inner) {
	let s = '<' + tag;
	for (const [k, v] of Object.entries(attrs)) s += ' ' + k + '="' + v + '"';
	if (inner) return s + '>' + inner + '</' + tag + '>';
	return s + '/>';
}

function toggleAccordion(id) {
	const content = document.getElementById(id);
	const icon = document.getElementById(id + '-icon');
	if (!content) return;
	if (content.classList.contains('open')) {
		content.classList.remove('open');
		if (icon) icon.dataset.state = 'closed';
	} else {
		content.classList.add('open');
		if (icon) icon.dataset.state = 'open';
	}
}

function toggleEnvValue(inputId, envVar, resolved) {
	const input = document.getElementById(inputId);
	const toggle = document.getElementById(inputId + '-toggle');
	if (!input) return;
	if (input.dataset.showingEnv === 'true') {
		input.value = resolved;
		input.dataset.showingEnv = 'false';
		input.readOnly = false;
		if (toggle) toggle.textContent = 'Show env template';
	} else {
		const tmpl = String.fromCharCode(123,123) + ' env ' + String.fromCharCode(34) + envVar + String.fromCharCode(34) + ' ' + String.fromCharCode(125,125);
		input.value = tmpl;
		input.dataset.showingEnv = 'true';
		input.readOnly = true;
		if (toggle) toggle.textContent = 'Show resolved value';
	}
}

function initSortable() {
	document.querySelectorAll('.sortable-list:not(.sortable-initialized)').forEach(el => {
		el.classList.add('sortable-initialized');
		new Sortable(el, {
			animation: 150,
			ghostClass: 'sortable-ghost',
			chosenClass: 'sortable-chosen',
			onEnd: function(evt) {
				const typeName = el.dataset.typeName;
				const items = el.querySelectorAll('[data-layout]');
				const layouts = Array.from(items).map(i => i.dataset.layout);
				fetch('/api/layout/' + typeName, {
					method: 'POST',
					headers: {'Content-Type': 'application/x-www-form-urlencoded'},
					body: 'data=' + encodeURIComponent(JSON.stringify(layouts))
				}).then(() => debounceGenerate(typeName));
			}
		});
	});
}

function initAccordions() {
	document.querySelectorAll('[data-accordion]').forEach(btn => {
		if (btn.dataset.accordionBound) return;
		btn.dataset.accordionBound = 'true';
		btn.addEventListener('click', function(e) {
			e.preventDefault();
			e.stopPropagation();
			toggleAccordion(btn.dataset.accordion);
		});
	});
}

function initEnvToggles() {
	document.querySelectorAll('[data-env-toggle]').forEach(toggle => {
		if (toggle.dataset.envBound) return;
		toggle.dataset.envBound = 'true';
		toggle.addEventListener('click', function(e) {
			e.preventDefault();
			toggleEnvValue(toggle.dataset.envToggle, toggle.dataset.envVar, toggle.dataset.envResolved);
		});
	});
}

function openDialog(id) {
	const el = document.getElementById(id);
	if (!el) return;
	el.classList.remove('hidden');
	const firstInput = el.querySelector('input:not([type=hidden]), select');
	if (firstInput) setTimeout(() => firstInput.focus(), 100);
}

function closeDialog(id) {
	document.getElementById(id)?.classList.add('hidden');
}

function handleDialogKeydown(e) {
	if (e.key === 'Escape') {
		['create-cv-dialog', 'rename-cv-dialog', 'delete-cv-dialog'].forEach(id => closeDialog(id));
	}
}

document.addEventListener('DOMContentLoaded', function() {
	const currentType = document.body.dataset.currentType;
	initSortable();
	initAccordions();
	initEnvToggles();

	if (document.getElementById('pdf-pages-container')) {
		import('https://cdnjs.cloudflare.com/ajax/libs/pdf.js/4.4.168/pdf.min.mjs').then(module => {
			window.pdfjsLib = module;
			module.GlobalWorkerOptions.workerSrc = 'https://cdnjs.cloudflare.com/ajax/libs/pdf.js/4.4.168/pdf.worker.min.mjs';
			loadPdf(currentType);
		}).catch(err => console.error('Failed to load PDF.js:', err));
	}

	document.getElementById('btn-create-cv')?.addEventListener('click', () => openDialog('create-cv-dialog'));
	document.getElementById('btn-create-cancel')?.addEventListener('click', () => closeDialog('create-cv-dialog'));
	document.getElementById('btn-create-confirm')?.addEventListener('click', () => closeDialog('create-cv-dialog'));

	document.getElementById('btn-rename-cv')?.addEventListener('click', () => openDialog('rename-cv-dialog'));
	document.getElementById('btn-rename-cancel')?.addEventListener('click', () => closeDialog('rename-cv-dialog'));
	document.getElementById('btn-rename-confirm')?.addEventListener('click', () => closeDialog('rename-cv-dialog'));

	document.getElementById('btn-delete-cv')?.addEventListener('click', () => {
		const name = document.body.dataset.currentType;
		document.getElementById('delete-cv-name').textContent = name;
		document.getElementById('delete-cv-confirm').setAttribute('hx-vals', '{"type_name": "' + name + '"}');
		openDialog('delete-cv-dialog');
	});
	document.getElementById('btn-delete-cancel')?.addEventListener('click', () => closeDialog('delete-cv-dialog'));

	document.addEventListener('keydown', handleDialogKeydown);

	document.getElementById('btn-generate-pdf')?.addEventListener('click', function() {
		generatePDF(document.body.dataset.currentType);
	});

	document.getElementById('btn-zoom-in')?.addEventListener('click', zoomIn);
	document.getElementById('btn-zoom-out')?.addEventListener('click', zoomOut);
	document.getElementById('btn-zoom-fit')?.addEventListener('click', zoomFit);

	const debounceRange = document.getElementById('debounce-range');
	const debounceInput = document.getElementById('debounce-input');
	if (debounceRange && debounceInput) {
		debounceRange.value = defaultDebounceMs;
		debounceInput.value = defaultDebounceMs;
		debounceRange.addEventListener('input', function () {
			const val = parseInt(this.value, 10);
			setDebounceMs(val);
			debounceInput.value = defaultDebounceMs;
		});
		debounceInput.addEventListener('change', function () {
			const val = parseInt(this.value, 10);
			if (isNaN(val)) return;
			setDebounceMs(val);
			this.value = defaultDebounceMs;
			debounceRange.value = defaultDebounceMs;
		});
	}
});

document.body.addEventListener('htmx:afterSettle', function() {
	initSortable();
	initAccordions();
	initEnvToggles();
});

document.body.addEventListener('htmx:afterRequest', function(event) {
	try {
		const path = event.detail?.pathInfo?.requestPath || event.detail?.xhr?.responseURL || '';
		if (path.includes('/api/generate')) return;
		const xhr = event.detail?.xhr;
		const triggerHeader = xhr?.getResponseHeader('HX-Trigger');
		if (triggerHeader === 'pdf-regenerate') {
			debounceGenerate(document.body.dataset.currentType);
		}
		if (xhr && xhr.status >= 200 && xhr.status < 300 && !path.includes('/api/cv-type/')) {
			if (window.AutoCV && window.AutoCV.showToast) {
				const section = path.split('/').pop() || path.split('/').slice(-2).join('/');
				window.AutoCV.showToast('Saved: ' + section, 'success', 1500);
			}
			updateSaveStatus('saved');
		}
		if (xhr && xhr.status >= 400) {
			if (window.AutoCV && window.AutoCV.showToast) {
				window.AutoCV.showToast('Save failed', 'error');
			}
			updateSaveStatus('error');
		}
	} catch(e) {
		console.error('htmx:afterRequest error:', e);
	}
});

document.body.addEventListener('htmx:beforeRequest', function(event) {
	const path = event.detail?.pathInfo?.requestPath || '';
	if (!path.includes('/api/generate') && !path.includes('/api/cv-type/')) {
		updateSaveStatus('saving');
	}
});

function updateSaveStatus(status) {
	const el = document.getElementById('save-status');
	if (!el) return;
	if (status === 'saving') {
		el.className = 'save-status saving';
		el.innerHTML = '<div class="spinner"></div> Saving...';
	} else if (status === 'saved') {
		el.className = 'save-status saved';
		el.textContent = 'Saved';
		setTimeout(function() {
			if (el.classList.contains('saved')) {
				el.className = 'save-status';
				el.textContent = '';
			}
		}, 2000);
	} else if (status === 'error') {
		el.className = 'save-status error';
		el.textContent = 'Error';
		setTimeout(function() {
			if (el.classList.contains('error')) {
				el.className = 'save-status';
				el.textContent = '';
			}
		}, 3000);
	}
}
