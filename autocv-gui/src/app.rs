mod ui;
use ui::general_ui;

use crate::app::ui::{experience_ui, project_ui, side_panel_ui};

#[derive(PartialEq)] // We need this to compare enum variants
pub enum AppTab {
    General,
    Experiences,
    Projects,
}

pub struct App {
    label: String,
    value: f32,
    selected_tab: AppTab,
}

impl Default for App {
    fn default() -> Self {
        Self {
            label: "Hello World!".to_owned(),
            value: 2.7,
            selected_tab: AppTab::General,
        }
    }
}

impl App {
    /// Called once before the first frame.
    pub fn new(cc: &eframe::CreationContext<'_>) -> Self {
        // This is also where you can customize the look and feel of egui using
        // `cc.egui_ctx.set_visuals` and `cc.egui_ctx.set_fonts`.

        Default::default()
    }
}

impl eframe::App for App {
    /// Called by the framework to save state before shutdown.
    fn save(&mut self, storage: &mut dyn eframe::Storage) {
        // eframe::set_value(storage, eframe::APP_KEY, self);
    }

    /// Called each time the UI needs repainting, which may be many times per second.
    fn update(&mut self, ctx: &egui::Context, _frame: &mut eframe::Frame) {
        // For inspiration and more examples, go to https://emilk.github.io/egui

        // egui::TopBottomPanel::top("top_panel").show(ctx, |ui| {
        //     // The top panel is often a good place for a menu bar:
        //
        //     egui::MenuBar::new().ui(ui, |ui| {
        //         ui.menu_button("File", |ui| {
        //             if ui.button("Quit").clicked() {
        //                 ctx.send_viewport_cmd(egui::ViewportCommand::Close);
        //             }
        //         });
        //         ui.add_space(16.0);
        //
        //         egui::widgets::global_theme_preference_buttons(ui);
        //     });
        // });

        side_panel_ui(ctx, self);

        egui::CentralPanel::default().show(ctx, |ui| {
            ui.heading("Configuration");
            ui.separator();

            egui::ScrollArea::vertical().show(ui, |ui| {
                match self.selected_tab {
                    AppTab::General => general_ui(ui),
                    AppTab::Experiences => experience_ui(ui),
                    AppTab::Projects => project_ui(ui),
                }
                ui.allocate_space(ui.available_size());
            });

            // ui.with_layout(egui::Layout::bottom_up(egui::Align::LEFT), |ui| {
            //     powered_by_egui_and_eframe(ui);
            //     egui::warn_if_debug_build(ui);
            // });
        });
    }
}

fn _powered_by_egui_and_eframe(ui: &mut egui::Ui) {
    ui.horizontal(|ui| {
        ui.spacing_mut().item_spacing.x = 0.0;
        ui.label("Powered by ");
        ui.hyperlink_to("egui", "https://github.com/emilk/egui");
        ui.label(" and ");
        ui.hyperlink_to(
            "eframe",
            "https://github.com/emilk/egui/tree/master/crates/eframe",
        );
        ui.label(".");
    });
}
