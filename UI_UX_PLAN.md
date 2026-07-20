# AutoCV UI/UX Improvement Plan

## Current State

AutoCV's web UI is a functional but basic split-panel editor built with:
- **Go + Templ** (server-side rendered HTML)
- **HTMX** (dynamic partial updates)
- **Tailwind CSS** (via CDN, no config)
- **SortableJS** (drag-and-drop layout reorder)
- **PDF.js** (PDF preview rendering)

### Key UX Problems

| # | Problem | Severity |
|---|---------|----------|
| 1 | No dark mode — light-only theme strains eyes during extended use | Critical |
| 2 | Unicode symbols used as icons (✎ ☰ ✕ ▶ ▼) — inconsistent, platform-dependent | Critical |
| 3 | No visual hierarchy — all accordion sections look identical | High |
| 4 | Debounce delay exposed in header — developer-facing, confuses users | High |
| 5 | No responsive layout — 50/50 split breaks on mobile/tablet | High |
| 6 | No keyboard shortcuts — power users can't navigate efficiently | High |
| 7 | No save feedback — silent saves, no toast/confirmation | High |
| 8 | No inline form validation — errors surface only at PDF generation | High |
| 9 | PDF preview lacks zoom, fit, and page navigation controls | Medium |
| 10 | No search/filter in long editor scroll | Medium |
| 11 | Dialog modals have no animation or backdrop blur | Medium |
| 12 | No onboarding or empty state guidance | Low |
| 13 | No undo/redo or change tracking | Low |
| 14 | Env variable UI is confusing for non-technical users | Low |

---

## Design System

### Style: Minimalism & Swiss Style

Clean, functional, professional — ideal for productivity tools and editor interfaces.

- **Light mode**: Full support ✓
- **Dark mode**: Full support ✓
- **Performance**: ⚡ Excellent
- **Accessibility**: ✓ WCAG AAA
- **Complexity**: Low

**Effects**: Subtle hover (200-250ms), smooth transitions, clear type hierarchy, fast loading. No unnecessary shadows or gradients.

**Avoid**: Flat design without depth, text-heavy pages without visual hierarchy.

### Color Palette: B2B Service (Professional Navy + Blue CTA)

| Token | Light | Dark | Purpose |
|-------|-------|------|---------|
| `--color-primary` | `#0F172A` | `#E2E8F0` | Headers, primary text |
| `--color-on-primary` | `#FFFFFF` | `#0F172A` | Text on primary bg |
| `--color-secondary` | `#334155` | `#CBD5E1` | Secondary text, borders |
| `--color-accent` | `#0369A1` | `#38BDF8` | CTA buttons, links, active states |
| `--color-on-accent` | `#FFFFFF` | `#0C4A6E` | Text on accent bg |
| `--color-background` | `#F8FAFC` | `#0F172A` | Page background |
| `--color-surface` | `#FFFFFF` | `#1E293B` | Cards, panels, modals |
| `--color-foreground` | `#020617` | `#F8FAFC` | Primary text |
| `--color-muted` | `#E8ECF1` | `#334155` | Muted backgrounds |
| `--color-muted-foreground` | `#64748B` | `#94A3B8` | Secondary/helper text |
| `--color-border` | `#E2E8F0` | `#334155` | Borders, dividers |
| `--color-destructive` | `#DC2626` | `#EF4444` | Errors, delete actions |
| `--color-success` | `#16A34A` | `#22C55E` | Success states |
| `--color-warning` | `#EA580C` | `#FB923C` | Warning states |
| `--color-ring` | `#0F172A` | `#38BDF8` | Focus ring |

### Typography: Inter

| Role | Size | Weight | Line Height |
|------|------|--------|-------------|
| Page title | 1.5rem (24px) | 700 | 1.2 |
| Section heading | 0.875rem (14px) | 600 | 1.4 |
| Body text | 0.8125rem (13px) | 400 | 1.5 |
| Label | 0.75rem (12px) | 500 | 1.4 |
| Caption/helper | 0.6875rem (11px) | 400 | 1.4 |

Google Fonts: `https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&display=swap`

### Spacing Scale (8px rhythm)

| Token | Value | Usage |
|-------|-------|-------|
| `--space-1` | 4px | Tight gaps |
| `--space-2` | 8px | Component internal |
| `--space-3` | 12px | Small gaps |
| `--space-4` | 16px | Standard padding |
| `--space-5` | 24px | Section gaps |
| `--space-6` | 32px | Large gaps |
| `--space-8` | 48px | Section dividers |

### Icon System: Lucide Icons

Replace all Unicode symbols with SVG icons from [Lucide](https://lucide.dev/):
- ☰ → `<svg>` grip-vertical / menu
- ✎ → `<svg>` pencil / edit
- ✕ → `<svg>` x
- ▶ / ▼ → `<svg>` chevron-right / chevron-down
- ● → `<svg>` circle-dot (status)

---

## Implementation Phases

### Phase 1: Foundation (Critical)

**Goal**: Establish design tokens, dark mode, proper icon system.

#### 1.1 — Design Token Infrastructure

**Files**: `internal/web/components.templ`, new `internal/web/static/styles.css`

- Remove Tailwind CDN `<script>`
- Add proper `tailwind.config.js` with:
  ```js
  module.exports = {
    darkMode: 'class',
    theme: {
      extend: {
        colors: {
          primary: 'var(--color-primary)',
          accent: 'var(--color-accent)',
          surface: 'var(--color-surface)',
          // ... all semantic tokens
        },
        fontFamily: {
          sans: ['Inter', 'system-ui', 'sans-serif'],
        },
      },
    },
  }
  ```
- Add CSS custom properties for both light and dark themes
- Add Google Fonts `<link>` for Inter

#### 1.2 — Dark Mode Toggle

**Files**: `internal/web/components.templ`

- Add dark mode toggle button in header (sun/moon icon)
- JS: Toggle `dark` class on `<html>`, persist to `localStorage`
- Respect `prefers-color-scheme` on first visit
- All existing Tailwind classes get `dark:` variants

#### 1.3 — SVG Icon System

**Files**: `internal/web/components.templ`

- Create reusable Templ components for common icons (chevron, pencil, x, grip, sun, moon, plus, search, etc.)
- Replace all Unicode symbols throughout the template
- Use consistent 16px/20px sizing with `currentColor` for theming

#### 1.4 — Visual Hierarchy

**Files**: `internal/web/components.templ`

- Differentiate "Shared" sections (blue badge) from "Type-specific" sections (purple badge)
- Add subtle left-border accent colors per section type
- Use `--color-surface` for panels, `--color-muted` for section headers
- Improve spacing rhythm with the 8px scale

---

### Phase 2: Layout & Navigation (High)

**Goal**: Responsive, efficient editor layout with better navigation.

#### 2.1 — Responsive Split-Panel

**Files**: `internal/web/components.templ`, new `internal/web/static/app.js`

- Add resizable divider between editor and preview (drag to resize)
- Minimum widths: editor 320px, preview 400px
- On mobile (<768px): stack vertically with tab switcher (Editor / Preview)
- Persist panel sizes to `localStorage`

#### 2.2 — Sidebar Navigation

**Files**: `internal/web/components.templ`

Replace the long scroll of accordion sections with a sidebar navigation pattern:

```
┌──────────────────────────────────────────────────┐
│ AutoCV   [Type ▾] ✎ ✕    🌙  [+ New]           │
├──────┬───────────────────┬───────────────────────┤
│ Nav  │  Section Content  │   PDF Preview         │
│      │                   │                       │
│ ▸ Lay│  (active section  │   ┌─────────────┐     │
│ ▸ Per│   form fields)    │   │             │     │
│ ▸ Set│                   │   │  PDF pages   │     │
│ ▸ Des│                   │   │             │     │
│ ▸ Exp│                   │   └─────────────┘     │
│ ▸ Pro│                   │                       │
│ ▸ Ski│                   │   Zoom: - 100% +      │
│ ▸ Edu│                   │   Page 1 of 1         │
│ ▸ Cer│                   │                       │
│ ▸ Awa│                   │                       │
├──────┴───────────────────┴───────────────────────┤
```

- Click nav item → HTMX loads that section's content into the main area
- Active nav item highlighted with accent border + background
- "Shared" sections get a subtle indicator (dot or badge)
- Reduce vertical scrolling from "scroll through 10 sections" to "click one item"

#### 2.3 — PDF Preview Improvements

**Files**: `internal/web/components.templ`

- Add zoom controls: fit-to-width, zoom in/out, percentage display
- Add page indicator: "Page 1 of 2"
- Add download button (direct PDF download)
- Add "Open in new tab" button
- Show page thumbnails on hover

#### 2.4 — Remove Debounce Control from Header

**Files**: `internal/web/components.templ`

- Move debounce to a settings/preferences dropdown
- Default to 500ms (current is 50ms which is aggressive)
- Advanced users can still adjust via gear icon dropdown

---

### Phase 3: Interaction Quality (High)

**Goal**: Polish micro-interactions and feedback.

#### 3.1 — Toast Notifications

**Files**: New `internal/web/static/toast.js`, `internal/web/components.templ`

- Add toast container (bottom-right)
- Toast types: success (green), error (red), warning (orange), info (blue)
- Auto-dismiss after 3s, manual dismiss on click
- Show toast on: save success, save error, PDF generated, PDF error
- HTMX response headers: `HX-Trigger: showToast` with toast data

#### 3.2 — Inline Form Validation

**Files**: `internal/web/components.templ`

- Validate email format on blur
- Validate required fields (name) with red border + error message
- Validate URL format for website/LinkedIn/GitHub fields
- Show character count for description textarea
- Validate accent-color hex format in settings

#### 3.3 — Loading States

**Files**: `internal/web/components.templ`

- PDF generation: Show skeleton placeholder in preview while generating
- HTMX requests: Show subtle loading bar at top of section
- Use `htmx-indicator` class for loading spinners on buttons
- Progress bar animation during PDF compilation

#### 3.4 — Keyboard Shortcuts

**Files**: `internal/web/components.templ`

- `Ctrl/Cmd + S` — Save current section (prevent default browser save)
- `Ctrl/Cmd + P` — Generate PDF preview
- `Ctrl/Cmd + 1-9` — Switch to section N
- `Ctrl/Cmd + D` — Toggle dark mode
- `?` — Show keyboard shortcut help overlay

#### 3.5 — Better Modal Dialogs

**Files**: `internal/web/components.templ`

- Backdrop blur (`backdrop-filter: blur(4px)`)
- Smooth enter/exit transitions (scale + opacity)
- Click outside to close
- `Escape` key to close
- Focus trap within modal

#### 3.6 — Smooth HTMX Transitions

**Files**: `internal/web/components.templ`

- Use `htmx-settling` class for CSS transitions on content swap
- Fade-in animation for swapped content
- Smooth height transitions when sections expand/collapse

---

### Phase 4: Productivity Features (Medium)

**Goal**: Power-user features for efficient CV editing.

#### 4.1 — Search/Filter in Editor

- Search bar at top of editor sidebar
- Filters sections and fields matching query
- Highlight matching text

#### 4.2 — Auto-Save Indicator

- Show "Saving..." indicator when HTMX request is in flight
- Show "Saved" with timestamp when complete
- Show "Unsaved changes" dot when content has been modified but not saved

#### 4.3 — Experience/Project Inline Editing

- Click-to-edit instead of form fields for simple values (company name, role, dates)
- Expandable detail panels for points/achievements
- Drag-to-reorder points within an experience

#### 4.4 — Section Visibility Toggles

- Quick toggle to show/hide sections not in the current CV type's layout
- Grayed-out sections with "Not in layout" badge
- One-click add to layout

#### 4.5 — Bulk Actions

- Duplicate experience/project entry
- Move experience between CV types
- Bulk reorder via drag-and-drop (already have SortableJS)

---

### Phase 5: Advanced (Nice-to-have)

#### 5.1 — Onboarding Walkthrough
- First-time user guide with step-by-step tooltips
- "Create your first CV type" flow

#### 5.2 — Version History
- Track changes per CV type with timestamps
- Diff view between versions
- Rollback to previous version

#### 5.3 — Template Gallery
- Browse and apply CV templates (pre-configured layout + style combos)
- Preview before applying

#### 5.4 — Export Options
- Export to DOCX, JSON, Markdown
- Share via link (read-only)

#### 5.5 — Real-time Preview
- WebSocket-based live preview as you type
- Eliminate debounce delay entirely

---

## Tailwind Dark Mode Strategy

Since we're using Tailwind with the `dark` class strategy:

```html
<html lang="en" class="">
  <!-- Toggle adds/removes 'dark' class -->
</html>
```

```html
<!-- Example: Light bg white, dark bg slate-800 -->
<div class="bg-white dark:bg-slate-800 text-slate-900 dark:text-slate-100">
```

**Migration approach**: Incrementally add `dark:` variants to existing classes. Start with the most impactful: backgrounds, text, borders.

---

## File Structure (Proposed)

```
internal/web/
├── app.go                    # Routes (unchanged)
├── components.templ          # Main templates (refactored)
├── components_templ.go       # Generated (auto)
├── static/
│   ├── styles.css            # CSS custom properties, transitions, animations
│   ├── app.js                # Dark mode, keyboard shortcuts, panel resize
│   └── toast.js              # Toast notification system
├── icons.templ               # Reusable SVG icon components
└── partials/
    ├── personal_info.templ   # Section partials (extracted from components.templ)
    ├── settings.templ
    ├── experiences.templ
    ├── projects.templ
    └── ...
```

---

## Pre-Delivery Checklist

### Visual Quality
- [ ] No emojis/Unicode used as icons (all SVG)
- [ ] Consistent icon family (Lucide) with uniform stroke width
- [ ] Semantic theme tokens used (no hardcoded hex values)
- [ ] Dark mode fully functional with proper contrast ratios

### Interaction
- [ ] All interactive elements have hover/focus states
- [ ] Toast notifications for all save/generate actions
- [ ] Form validation with clear error messages
- [ ] Keyboard shortcuts working with help overlay
- [ ] Modal dialogs have backdrop blur and focus trap

### Accessibility
- [ ] Text contrast ≥4.5:1 in both light and dark mode
- [ ] All form inputs have proper `<label>` elements
- [ ] Focus-visible rings on all interactive elements
- [ ] `prefers-reduced-motion` respected
- [ ] Screen reader labels on icon-only buttons
- [ ] ARIA attributes on accordion sections

### Responsive
- [ ] Tested on 375px (mobile), 768px (tablet), 1440px (desktop)
- [ ] Mobile: stacked layout with editor/preview tab switcher
- [ ] Touch targets ≥44px on mobile
- [ ] No horizontal scroll on any breakpoint

### Layout
- [ ] 8px spacing rhythm maintained
- [ ] Safe areas respected on mobile
- [ ] Resizable split-panel with persistent sizes
- [ ] PDF preview with zoom controls
