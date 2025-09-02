mod style;
mod ui;

use std::sync::{Arc, Mutex};

use log::{info, warn};
use tokio::sync::mpsc;
use ui::general_ui;

use crate::{
    actor::{ActorMessage, State},
    app::style::apply_style,
    app::ui::{experience_ui, project_ui, side_panel_ui},
};

#[derive(PartialEq)]
pub enum AppTab {
    General,
    Experiences,
    Projects,
}

pub struct App {
    selected_tab: AppTab,
    shared_state: Arc<Mutex<State>>,
    local_state: State,
    actor_sender: mpsc::Sender<ActorMessage>,
}

impl App {
    /// Called once before the first frame.
    pub fn new(
        cc: &eframe::CreationContext<'_>,
        actor_sender: mpsc::Sender<ActorMessage>,
        shared_state: Arc<Mutex<State>>,
    ) -> Self {
        // This is also where you can customize the look and feel of egui using
        // `cc.egui_ctx.set_visuals` and `cc.egui_ctx.set_fonts`.
        // Default::default()

        let style = apply_style();
        cc.egui_ctx.set_style(style);

        Self {
            selected_tab: AppTab::General,
            local_state: Arc::clone(&shared_state).lock().unwrap().clone(),
            shared_state,
            actor_sender,
        }
    }
}

impl eframe::App for App {
    /// Called by the framework to save state before shutdown.
    fn save(&mut self, _storage: &mut dyn eframe::Storage) {
        // eframe::set_value(storage, eframe::APP_KEY, self);
    }

    /// Called each time the UI needs repainting, which may be many times per second.
    fn update(&mut self, ctx: &egui::Context, _frame: &mut eframe::Frame) {
        // For inspiration and more examples, go to https://emilk.github.io/egui

        let top_panel_frame = style::top_panel_frame();

        // Top Panel
        egui::TopBottomPanel::top("top_panel")
            .frame(top_panel_frame)
            .show(ctx, |ui| {
                ui.horizontal(|ui| {
                    // --- Left side of the panel ---
                    ui.with_layout(egui::Layout::left_to_right(egui::Align::Center), |ui| {
                        ui.label(egui::RichText::new("AutoCV Config").strong().size(16.0));
                    });

                    // --- Right side of the panel ---
                    ui.with_layout(egui::Layout::right_to_left(egui::Align::Center), |ui| {
                        let compile_button = egui::Button::new(
                            egui::RichText::new("Compile PDF")
                                .color(crate::app::style::ayu_dark::BG_DARK)
                                .strong(),
                        )
                        .fill(crate::app::style::ayu_dark::ACCENT_YELLOW);

                        if ui
                            .add(compile_button)
                            .on_hover_text("Compile the CV configuration to a PDF file")
                            .clicked()
                        {
                            // TODO: Send compile message to actor
                            info!("Compile button clicked!");
                        }
                    });
                });
            });

        side_panel_ui(ctx, self);

        egui::CentralPanel::default().show(ctx, |ui| {
            ui.heading("Configuration");
            ui.separator();

            egui::ScrollArea::vertical().show(ui, |ui| {
                match self.selected_tab {
                    AppTab::General => general_ui(ui, &mut self.local_state.general_cv),
                    AppTab::Experiences => {
                        for experience in self.local_state.experiences_cv.iter_mut() {
                            experience_ui(ui, experience);
                            ui.separator();
                        }
                    }
                    AppTab::Projects => {
                        for project in self.local_state.projects_cv.iter_mut() {
                            project_ui(ui, project);
                            ui.separator();
                        }
                    }
                }
                ui.allocate_space(ui.available_size());
            });

            // ui.separator();
            // ui.heading("Content from Actor:");
            //
            // if let Ok(state) = self.shared_state.try_lock() {
            //     info!("Successfully locked the shared state.",);
            //     ui.label(&state.general_cv.personal_info.name);
            // } else {
            //     warn!("Failed to lock the shared state.");
            //     ui.label("State is being updated...");
            // }
            //
            // Lock the mutex to safely read the data
            // let state = self.shared_state.lock().unwrap();
            // ui.label(&state);

            // ctx.request_repaint(); // Ensures the UI updates when the actor changes the state
        });
    }
}
