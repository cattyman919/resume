// src/app/style.rs

use egui::{
    Color32, Stroke, Style,
    style::{Interaction, Visuals, WidgetVisuals},
};

// Define the color palette for our Ayu Dark theme
mod ayu_dark {
    use egui::Color32;
    pub const BG_DARK: Color32 = Color32::from_rgb(13, 17, 23); // Background
    pub const BG_MEDIUM: Color32 = Color32::from_rgb(22, 27, 34); // Panels, frames
    pub const BG_LIGHT: Color32 = Color32::from_rgb(33, 38, 45); // Widget backgrounds
    pub const FG_DARK: Color32 = Color32::from_rgb(122, 129, 142); // Less important text
    pub const FG_LIGHT: Color32 = Color32::from_rgb(182, 189, 202); // Main text
    pub const ACCENT_YELLOW: Color32 = Color32::from_rgb(255, 198, 0); // Primary accent
    pub const ACCENT_ORANGE: Color32 = Color32::from_rgb(255, 167, 26); // Other accents
}

pub fn apply_style() -> Style {
    let mut style = Style::default();
    let mut visuals = Visuals::dark();

    // --- Main Colors ---
    visuals.override_text_color = Some(ayu_dark::FG_LIGHT);
    visuals.window_fill = ayu_dark::BG_MEDIUM; // Window background
    visuals.panel_fill = ayu_dark::BG_DARK; // Side/central panel background

    // --- Widget Visuals ---
    visuals.widgets.noninteractive.bg_fill = ayu_dark::BG_LIGHT; // Buttons, text boxes
    visuals.widgets.noninteractive.fg_stroke = Stroke::new(1.0, ayu_dark::FG_DARK);
    visuals.widgets.noninteractive.bg_stroke = Stroke::new(1.0, ayu_dark::BG_LIGHT);

    // --- Interaction States ---
    let mut interaction = WidgetVisuals {
        bg_fill: ayu_dark::ACCENT_YELLOW.linear_multiply(0.4),
        bg_stroke: Stroke::new(2.0, ayu_dark::ACCENT_ORANGE),
        fg_stroke: Stroke::new(1.5, ayu_dark::ACCENT_ORANGE),
        weak_bg_fill: ayu_dark::BG_LIGHT,
        corner_radius: egui::CornerRadius::same(4),
        expansion: 0.0,
    };
    interaction.bg_fill = ayu_dark::ACCENT_YELLOW.linear_multiply(0.2); // Hover color
    interaction.fg_stroke = Stroke::new(2.0, ayu_dark::ACCENT_YELLOW); // Hover text/outline

    visuals.widgets.hovered = interaction;
    visuals.widgets.active = interaction; // Color when clicking

    // --- Specific Widget Styling ---
    visuals.selection.bg_fill = ayu_dark::ACCENT_YELLOW.linear_multiply(0.4);
    visuals.selection.stroke = Stroke::new(2.0, ayu_dark::ACCENT_ORANGE);
    visuals.hyperlink_color = ayu_dark::ACCENT_ORANGE;
    visuals.collapsing_header_frame = true;

    style.visuals = visuals;

    // --- Spacing and Sizing ---
    style.spacing.item_spacing = egui::vec2(8.0, 8.0);
    style.spacing.button_padding = egui::vec2(10.0, 6.0);
    style.spacing.window_margin = egui::Margin::same(12);

    style
}
